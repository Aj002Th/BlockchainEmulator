package network

import (
	"bufio"
	"bytes"
	"errors"
	"io"
	"log"
	"net"
	"sync"

	"github.com/Aj002Th/BlockchainEmulator/signal"
)

// TcpCustomProtocolNetwork 基于 tcp 自定义的应用层协议
// 以 '\n' 作为消息的边界
type TcpCustomProtocolNetwork struct {
	connMapLock    sync.Mutex
	connectionPool map[string]net.Conn
	OnUpload       signal.Signal[int]
	OnDownload     signal.Signal[int]
	logger         log.Logger
	tcpLock        sync.Mutex
	recvChan       chan []byte
	doneChan       chan struct{}
}

func NewTcpCustomProtocolNetwork() *TcpCustomProtocolNetwork {
	return &TcpCustomProtocolNetwork{
		connectionPool: make(map[string]net.Conn),
		OnUpload:       signal.NewAsyncSignalImpl[int]("TcpOnUpload"),
		OnDownload:     signal.NewAsyncSignalImpl[int]("TcpOnDownload"),
		logger:         *log.Default(),
		recvChan:       make(chan []byte),
		doneChan:       make(chan struct{}),
	}
}

// Send 发送消息
func (t *TcpCustomProtocolNetwork) Send(content []byte, addr string) {
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
				delete(t.connectionPool, addr)
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
		go t.readFromConn(addr)
	}

	_, err = conn.Write(append(content, '\n'))
	t.UpdateMetric(len(content)+1, 0) // 带上反斜杠n的长度计算。
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

// Serve 传入endpoint应当满足 IPv4地址:端口号
func (t *TcpCustomProtocolNetwork) Serve(endpoint string) chan []byte { // 不停听并且起goroutine
	go func() {
		tcpLn, err := net.Listen("tcp", endpoint)
		if err != nil {
			log.Panic(err)
		}

		go func() {
			<-t.doneChan
			tcpLn.Close() // SO说，如何让net.Listener.Accept()优雅Cancellationn呢？可以直接Ln.Close()，诱导它返回错误。这样就能顺利解决问题。
		}()

		for {
			conn, err := tcpLn.Accept()
			if err != nil {
				close(t.recvChan)
				t.logger.Println("recvChan close")
				return
			}
			t.logger.Printf("Accepted the: %v. Now Start a session.\n", conn.RemoteAddr())
			go t.startSession(conn)
		}
	}()
	return t.recvChan
}

func (d *TcpCustomProtocolNetwork) startSession(con net.Conn) {
	defer con.Close()
	clientReader := bufio.NewReader(con)
	for {
		clientRequest, err := clientReader.ReadBytes('\n')
		d.OnDownload.Emit(len(clientRequest))
		switch err {
		case nil:
			d.tcpLock.Lock()
			d.recvChan <- clientRequest
			d.tcpLock.Unlock()
		case io.EOF:
			log.Println("client closed the connection by terminating the process")
			return
		default:
			log.Printf("error: %v\n", err)
			return
		}
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

	t.doneChan <- struct{}{}
}

// 读取完整消息并打印到日志
func (t *TcpCustomProtocolNetwork) readFromConn(addr string) {
	conn := t.connectionPool[addr]

	buffer := make([]byte, 1024)
	var messageBuffer bytes.Buffer

	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if errors.Is(err, net.ErrClosed) {
				break
			}
			if err != io.EOF {
				log.Println("Unexpected conn.Read error for address", addr, ":", err, ". Now break the reading loop")
			}
			break
		}

		// 将消息写入 buffer
		messageBuffer.Write(buffer[:n])
		t.UpdateMetric(0, n)

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

func (t *TcpCustomProtocolNetwork) GetOnUpload() signal.Signal[int] {
	return t.OnUpload
}

func (t *TcpCustomProtocolNetwork) GetOnDownload() signal.Signal[int] {
	return t.OnDownload
}

func (t *TcpCustomProtocolNetwork) UpdateMetric(up int, down int) { // 单位是字节数
	t.OnUpload.Emit(up)
	t.OnDownload.Emit(down)
}

// 从 buffer 中读取一行
func readMessage(buffer *bytes.Buffer) (string, error) {
	message, err := buffer.ReadBytes('\n')
	if err != nil && err != io.EOF {
		return "", err
	}
	return string(message), nil
}
