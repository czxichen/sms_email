package sns

import "sync"

// Sender 发送器
type Sender interface {
	Send(contact, subject, content string) error
}

var senderMap = struct {
	mux  *sync.RWMutex
	maps map[string]Sender
}{mux: new(sync.RWMutex), maps: make(map[string]Sender)}

// SenderRegister 注册sender
func SenderRegister(provider, typ string, sender Sender) {
	senderMap.mux.Lock()
	senderMap.maps[provider+typ] = sender
	senderMap.mux.Unlock()
}

// SenderGet 获取sender
func SenderGet(provider, typ string) (Sender, bool) {
	senderMap.mux.RLock()
	defer senderMap.mux.RUnlock()
	sender, exist := senderMap.maps[provider+typ]
	return sender, exist
}
