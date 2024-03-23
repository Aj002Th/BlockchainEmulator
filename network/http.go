package network

import (
	"bytes"
	"log"
	"net/http"
)

type HttpNetwork struct {
}

func (h *HttpNetwork) Send(context []byte, addr string) {
	buf := bytes.NewBuffer(context)
	resp, err := http.Post(addr, "application/json", buf)
	if err != nil {
		// 处理http错误
		log.Println("Error processing message for address", addr, ":", err)
		return
	}
	// 将完整信息打印记录
	log.Println("Received from", addr, ":", resp)
}

func (h *HttpNetwork) Broadcast(sender string, receivers []string, msg []byte) {
	for _, ip := range receivers {
		if ip == sender {
			continue
		}
		go h.Send(msg, ip)
	}
}

func (h *HttpNetwork) Close() {
	// nothing to do
}
