package sns

import (
	"bytes"
	"fmt"
	"sync"
	"text/template"
)

var templateMap = struct {
	mux   *sync.RWMutex
	temps map[string]*template.Template
}{
	mux:   new(sync.RWMutex),
	temps: make(map[string]*template.Template),
}

// TemplateRegister 模版注册
// 模版中至少应该有一个Code变量
func TemplateRegister(config *TemplateRegisterBase) error {
	temp := &template.Template{}
	for _, tempConfig := range config.Temps {
		t := temp.New(tempConfig.ContenName())
		if _, err := t.Parse(tempConfig.Content); err != nil {
			return err
		}
		if tempConfig.Extend != "" {
			t = t.New(tempConfig.ExtendName())
			if _, err := t.Parse(tempConfig.Extend); err != nil {
				return err
			}
		}
	}
	templateMap.mux.Lock()
	templateMap.temps[config.Application] = temp
	templateMap.mux.Unlock()
	return nil
}

// TemplateToMsg 模版生成
func TemplateToMsg(cfg *SendCodeConfig) (string, string, error) {
	// func TemplateToMsg(app, typ, tag string, data map[string]interface{}) (string, string, error) {
	templateMap.mux.RLock()
	temp, exist := templateMap.temps[cfg.Application]
	templateMap.mux.RLock()
	if !exist {
		return "", "", fmt.Errorf("application %s not register", cfg.Application)
	}
	var buf = bytes.NewBuffer(nil)
	var name = cfg.ContenName()
	t := temp.Lookup(name)
	if t == nil {
		return "", "", fmt.Errorf("template %s not register", name)
	}
	if err := t.Execute(buf, cfg.Data); err != nil {
		return "", "", fmt.Errorf("Execaute content template error:%s", err.Error())
	}
	content := buf.String()
	buf.Reset()
	var extend string
	name = cfg.ExtendName()
	if t = temp.Lookup(name); t != nil {
		if err := t.Execute(buf, cfg.Data); err != nil {
			return "", "", fmt.Errorf("Execaute extend template error:%s", err.Error())
		}
		extend = buf.String()
	}
	return extend, content, nil
}
