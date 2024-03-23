package network

import (
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"sync"
)

// TcpCustomProtocolNetwork 基于 tcp 自定义的应用层协议
// 以 '\n' 作为消息的边界
type TcpCustomProtocolNetwork struct {
	connMapLock    sync.Mutex
	connectionPool map[string]net.Conn
}

func NewTcpCustomProtocolNetwork() *TcpCustomProtocolNetwork {
	return &TcpCustomProtocolNetwork{
		connectionPool: make(map[string]net.Conn),
	}
}

// Send 发送消息
func (t *TcpCustomProtocolNetwork) Send(context []byte, addr string) {
	t.connMapLock.Lock()
	defer t.connMapLock.Unlock()

	var err error
	var conn net.Conn

	// 如果连接池中存在该连接，则复用该连接
	// 如果连接池中不存在该连接，则新建连接
	if c, ok := t.connectionPool[addr]; ok {
		if tcpConn, tcpOk := c.(*net.TCPConn); tcpOk {
			// 判断连接是否存活
			if err := tcpConn.SetKeepAlive(true); err != nil {
				delete(t.connectionPool, addr) // Remove if not alive
				conn, err = net.Dial("tcp", addr)
				if err != nil {
					log.Println("Reconnect error", err)
					return
				}
				t.connectionPool[addr] = conn
				go t.readFromConn(addr) // 用于打印日志
			} else {
				conn = c // 复用连接池中的连接
			}
		}
	} else {
		conn, err = net.Dial("tcp", addr)
		if err != nil {
			log.Println("Connect error", err)
			return
		}
		t.connectionPool[addr] = conn
		go t.readFromConn(addr) // Start reading from new connection
	}

	_, err = conn.Write(append(context, '\n'))
	if err != nil {
		log.Println("Write error", err)
		return
	}
}

// Broadcast 广播消息
func (t *TcpCustomProtocolNetwork) Broadcast(sender string, receivers []string, msg []byte) {
	for _, ip := range receivers {
		if ip == sender {
			continue
		}
		go t.Send(msg, ip)
	}
}

// Close 关闭所有连接并重置连接池
func (t *TcpCustomProtocolNetwork) Close() {
	t.connMapLock.Lock()
	defer t.connMapLock.Unlock()

	for _, conn := range t.connectionPool {
		_ = conn.Close()
	}
	t.connectionPool = make(map[string]net.Conn) // 重置连接池
}

// 读取完整消息并打印到日志
func (t *TcpCustomProtocolNetwork) readFromConn(addr string) {
	conn := t.connectionPool[addr]

	buffer := make([]byte, 1024)
	var messageBuffer bytes.Buffer

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				log.Println("Read error for address", addr, ":", err)
			}
			break
		}

		// 将消息写入 buffer
		messageBuffer.Write(buffer[:n])

		// 处理完整的消息
		for {
			message, err := readMessage(&messageBuffer)
			if errors.Is(err, io.ErrShortBuffer) {
				// 消息不完整, 需要继续从 buffer 中读取信息
				break
			} else if err == nil {
				// 将完整信息打印记录
				log.Println("Received from", addr, ":", message)
			} else {
				// 处理其他错误
				log.Println("Error processing message for address", addr, ":", err)
				break
			}
		}
	}
}

// 从 buffer 中读取一行
func readMessage(buffer *bytes.Buffer) (string, error) {
	message, err := buffer.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return string(message), nil
}
