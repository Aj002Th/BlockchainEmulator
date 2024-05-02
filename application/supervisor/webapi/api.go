package webapi

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"sync"

	"github.com/Aj002Th/BlockchainEmulator/application/supervisor/metrics"
	"github.com/Aj002Th/BlockchainEmulator/signal"
	"github.com/gorilla/websocket"
)

// 现在的代码还得改改。
// 为了实现：
// 当检测到-f frontend，就在创建sup的时候先弄到那个upgrade的c，并且放进一个Proxy，塞给sup用。
// 没检测到就给个哑的Proxy给sup用。

type Msg struct { // json消息反序列化。
	Type    string      `json:"type"`
	Content interface{} `json:"content"`
}

type Progress struct {
	Count int `json:"count"`
	Total int `json:"total"`
}

type ApiProxy interface {
	Enqueue(m Msg)
	writeToConnNoConsume(c *websocket.Conn) error
}

type GoodApiProxy struct {
	queue         []Msg
	Append_Signal signal.Signal[Msg]
	mtx           sync.RWMutex
}

func NewGoodApiProxy() *GoodApiProxy {
	return &GoodApiProxy{
		queue:         make([]Msg, 0),
		Append_Signal: signal.NewAsyncSignalImpl[Msg]("appendSignal"),
	}
}

// 给PBFT模块用。
func (ap *GoodApiProxy) Enqueue(m Msg) {
	ap.mtx.Lock() //写锁成对
	ap.queue = append(ap.queue, m)
	ap.mtx.Unlock()
	ap.Append_Signal.Emit(m)
}

func (ap *GoodApiProxy) dequeue(m *Msg) error {
	if len(ap.queue) > 1 {
		h := ap.queue[0]
		ap.queue = ap.queue[1:]
		*m = h
		return nil
	} else {
		return errors.New("len of q is zero")
	}
}

type Nothing struct{} // 表示不携带东西的空元组事件data

var ValNothing = Nothing{} // 一个方便的单例。

func writeToCMsg(c *websocket.Conn, m Msg) error {

	bts, err := json.Marshal(m)
	if err != nil {
		log.Print("upgrade:", err)
		return errors.New("the json error in proxy writeToConn")
	}
	err = c.WriteMessage(websocket.TextMessage, bts)
	if err != nil {
		log.Printf("write msg %v to websocket meet error: %v. \n", m, err)
		return errors.New("the Conn.WriteMessage got some err")
	}
	return nil
}

// 不是private，是包内共享方法。给echo用的
func (ap *GoodApiProxy) writeToConnNoConsume(c *websocket.Conn) error { // without dequeue any elems
	others := make(chan Msg)

	ap.mtx.RLock() // 读锁成对。在入队前注册监听器确保不错过事件。
	var appendCB func(m Msg)
	appendCB = func(m Msg) { // 那个asyncImpl能保证顺序
		others <- m
		if m.Type == "bye" {
			close(others)
			log.Print("Bye Sent")
			ap.Append_Signal.Disconnect(appendCB)
		}
	}
	ap.Append_Signal.Connect(appendCB)

	// 把之前的历史消息发一遍。
	qCopy := make([]Msg, 0)
	qCopy = append(qCopy, ap.queue...)
	ap.mtx.RUnlock()

	for _, m := range qCopy {
		// writeToCMsg(c*websocket.Conn,c)
		err := writeToCMsg(c, m)
		if err != nil {
			log.Print("writeToConn Chase Failed")
			return errors.New("writeToConn Chase Failed")
		}
	}

	// 现在发送新到达的
	for {
		m, isOpen := <-others
		if !isOpen {
			break
		}
		err := writeToCMsg(c, m)
		if err != nil {
			log.Print("writeToConn Append Failed: ", err)
			return errors.New("writeToConn Append Failed")
		}
	}
	return nil
}

var GlobalProxy ApiProxy

type DumbProxy struct {
}

func (dp DumbProxy) Enqueue(m Msg) {
	// do nothing，intentially left blank
}

func (dp DumbProxy) writeToConnNoConsume(c *websocket.Conn) error {
	// do nothing, intentially left blank
	return nil
}

var upgrader = websocket.Upgrader{} // use default options

// 为了让echo处理器附加上一个chan状态，我们用闭包实现穷人的对象。
func echo(w http.ResponseWriter, r *http.Request) { // golang http server是多线程的。
	// 这个websocket很奇葩，是依赖标准库的http服务器，并且用http提升成websocket协议。
	// 原因是websocket其实是HTTP/2.0兼容的。
	upgrader.CheckOrigin = func(r *http.Request) bool { return true } // 首先解决一下跨域问题。
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	// defer c.Close()
	// 事件循环。现在c变成了buffered的双向队列。
	GlobalProxy.writeToConnNoConsume(c)
}

// 构建一个富有状态的echo函数和相应状态chan。并且返回chan供写入。
func RunApiServer() {
	log.SetFlags(0)
	http.HandleFunc("/api/status", echo)
	http.ListenAndServe("0.0.0.0:7697", nil)
}

type PbftItem struct {
	TxpoolSize int `json:"txpool_size"`
	Tx         int `json:"tx"`
}

type MeasureItem struct {
	Name string    `json:"name"`
	Desc string    `json:"desc"`
	Vals []float64 `json:"vals"`
}

type Re struct {
	PbftShardCsv   []PbftItem    `json:"pbftShardCsv"`
	MeasureOutputs []MeasureItem `json:"measureOutputs"`
}

type Re1 struct {
	PbftShardCsv   []PbftItem     `json:"pbftShardCsv"`
	MeasureOutputs []metrics.Desc `json:"measureOutputs"`
}

// 通信消息

var Hello = Msg{Type: "hello"}

var Started = Msg{Type: "started"}

var Computing = func(Total int, Count int) Msg {
	return Msg{
		Type: "computing",
		Content: Progress{
			Total: Total,
			Count: Count,
		},
	}
}

var Completed = func(PbftShardCsv []PbftItem, desc []metrics.Desc) Msg {
	return Msg{
		Type:    "completed",
		Content: Re1{PbftShardCsv: PbftShardCsv, MeasureOutputs: desc},
	}
}

var Bye = Msg{Type: "bye"}
