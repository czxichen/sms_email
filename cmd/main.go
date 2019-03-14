package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"net/rpc"
	"net/rpc/jsonrpc"
	"sns"

	"github.com/alecthomas/kingpin"
)

// EnvPrefix 环境变量前缀
const EnvPrefix = "Verification_CODE_"

var config = struct {
	Listen string `json:"listen"`
	User   string `json:"user"`
	Passwd string `json:"passwd"`
	DB     string `json:"db"`
	Crt    string `json:"crt"`
	Key    string `json:"key"`
	Debug  bool   `json:"debug"`
}{}

func init() {
	kingpin.Flag("listen", "listen address and port").Short('l').Default("0.0.0.0:1789").Envar(EnvPrefix + "LISTEN").StringVar(&config.Listen)
	kingpin.Flag("db", "mysql connection url:'user:passwd@tcp(ip:port)/dbName?charset=utf8&parseTime=True&loc=Local'").Envar(EnvPrefix + "DB").StringVar(&config.DB)
	kingpin.Flag("user", "basic auth user").Short('u').Envar(EnvPrefix + "USER").StringVar(&config.User)
	kingpin.Flag("passwd", "basic auth passwd").Short('p').Envar(EnvPrefix + "PASSWD").StringVar(&config.Passwd)
	kingpin.Flag("crt", "tls crt path").Short('c').Envar(EnvPrefix + "CRT").StringVar(&config.Crt)
	kingpin.Flag("key", "tls key path").Short('k').Envar(EnvPrefix + "KEY").StringVar(&config.Key)
	kingpin.Flag("debug", "open debug mode").Short('d').Envar(EnvPrefix + "DEBUG").BoolVar(&config.Debug)
	kingpin.Parse()
}

func main() {
	err := sns.InitMysql(config.DB, config.Debug)
	if err != nil {
		panic(err)
	}
	defer sns.CloseMysql()

	handler := new(sns.Handler)
	jsonRPCServer := &JSONRPCServer{Server: rpc.NewServer()}
	if err := jsonRPCServer.Server.Register(handler); err != nil {
		panic(err)
	}
	basicAuth := BasicAuth(config.User, config.Passwd)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		log.Printf("Remoter host:%s URI:%s\n", r.RemoteAddr, r.URL.Path)
		if !basicAuth(r) {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		switch r.URL.Path {
		case "/rpc":
			jsonRPCServer.ServeHTTP(w, r)
		case "/ping":
			fmt.Fprint(w, "ok")
		default:
			// TODO: http handler
		}
	})
	if config.Crt != "" || config.Key != "" {
		err = http.ListenAndServeTLS(config.Listen, config.Crt, config.Key, nil)
	} else {
		err = http.ListenAndServe(config.Listen, nil)
	}
	if err != nil {
		log.Printf("Listen error:%s\n", err.Error())
	}
}

// JSONRPCServer json rpc server
type JSONRPCServer struct {
	Server *rpc.Server
}

// ServeHTTP implements an http.Handler that answers RPC requests.
func (jsonRPCServer *JSONRPCServer) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.Method != "CONNECT" {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusMethodNotAllowed)
		io.WriteString(w, "405 must CONNECT\n")
		return
	}
	conn, _, err := w.(http.Hijacker).Hijack()
	if err == nil {
		io.WriteString(conn, "HTTP/1.0 200 Connected to Go RPC\n\n")
		jsonRPCServer.Server.ServeCodec(jsonrpc.NewServerCodec(conn))
	} else {
		log.Printf("rpc hijacking %s:%s\n", req.RemoteAddr, err.Error())
	}
}

// BasicAuth 基础认证
func BasicAuth(user, passwd string) func(r *http.Request) bool {
	if user != "" || passwd != "" {
		return func(r *http.Request) bool {
			u, p, ok := r.BasicAuth()
			if ok {
				return u == user && p == passwd
			}
			return false
		}
	}
	return func(r *http.Request) bool { return true }
}
