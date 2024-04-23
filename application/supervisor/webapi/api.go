package webapi

import (
	"encoding/json"
	"errors"
	"flag"
	"log"
	"net/http"
	"sync"

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
	Total int `json:"total`
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

func NewGoodApiProxy() GoodApiProxy {
	return GoodApiProxy{
		queue:         make([]Msg, 0),
		Append_Signal: signal.NewAsyncSignalImpl[Msg]("appendSignal"),
	}
}

// 给PBFT模块用。
func (ap GoodApiProxy) Enqueue(m Msg) {
	ap.mtx.Lock() //写锁成对
	ap.queue = append(ap.queue, m)
	ap.mtx.Unlock()
	ap.Append_Signal.Emit(m)
}

func (ap GoodApiProxy) dequeue(m *Msg) error {
	if len(ap.queue) > 1 {
		h := ap.queue[0]
		ap.queue = ap.queue[1:]
		*m = h
		return nil
	} else {
		return errors.New("len of q is zero")
	}
}

// TODO：尝试用CSP来实现。。算了，太抽象了，用signal吧。

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
		log.Println("write msg to websocket error: ", err)
		return errors.New("the Conn.WriteMessage got some err")
	}
	return nil
}

// 不是private，是包内共享方法。给echo用的
func (ap GoodApiProxy) writeToConnNoConsume(c *websocket.Conn) error { // without dequeue any elems
	// FIXME: 这里要修复, 需要生产者消费者来同步。
	var chased chan Nothing
	var append_cb func(m Msg)
	append_cb = func(m Msg) {
		<-chased //为了同步。
		err := writeToCMsg(c, m)
		if err != nil {
			log.Print("append_Cb err")
			ap.Append_Signal.Disconnect(append_cb)
		}
	}
	ap.Append_Signal.Connect(append_cb)
	ap.mtx.RLock() // 读锁成对。
	var q_copy []Msg
	copy(q_copy, ap.queue)
	ap.mtx.RUnlock()

	for _, m := range q_copy {
		// writeToCMsg(c*websocket.Conn,c)
		err := writeToCMsg(c, m)
		if err != nil {
			log.Print("writeToConn Chase Failed")
			return errors.New("writeToConn Chase Failed")
		}
	}
	chased <- ValNothing
	return nil
}

var G_Proxy ApiProxy = NewGoodApiProxy()

type DumbProxy struct {
}

func (dp DumbProxy) Enqueue(m Msg) {
	// do nothing，intentially left blank
}

func (dp DumbProxy) writeToConnNoConsume(c *websocket.Conn) error {
	// do nothing, intentially left blank
	return nil
}

var addr = flag.String("addr", "0.0.0.0:7697", "http service address")

var upgrader = websocket.Upgrader{} // use default options

// 为了让echo处理器附加上一个chan状态，我们用闭包实现穷人的对象。
func echo(w http.ResponseWriter, r *http.Request) {
	// 这个websocket很奇葩，是依赖标准库的http服务器，并且用http提升成websocket协议。
	// 原因是websocket其实是HTTP/2.0兼容的。
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()
	// 事件循环。现在c变成了buffered的双向队列。
	G_Proxy.
		writeToConnNoConsume(c)
}

// 构建一个富有状态的echo函数和相应状态chan。并且返回chan供写入。
func RunApiServer() {
	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/", echo)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

type PbftItem struct {
	TxpoolSize int `json:"txpool_size"`
	Tx         int `json:"tx"`
	Ctx        int `json:"ctx"`
}

type MeasureItem struct {
	Name string    `json:"name"`
	Vals []float64 `json:"vals"`
}

type Re struct {
	PbftShardCsv   []PbftItem    `json:"pbftShardCsv"`
	MeasureOutputs []MeasureItem `json:"measureOutputs"`
}

// 示例代码。不应该放在这里。
func rrTestCase() {
	go RunApiServer()

}

var Hello = Msg{Type: "hello"}
var Started = Msg{Type: "started"}
var Computing = func(Total int, Count int) Msg {
	return Msg{Type: "computing", Content: Progress{Total: Total, Count: Count}}
}
var Completed = func(PbftShardCsv []PbftItem, MeasureOutputs []MeasureItem) Msg {
	return Msg{Type: "completed", Content: Re{PbftShardCsv: PbftShardCsv, MeasureOutputs: MeasureOutputs}}
}
var Bye = Msg{Type: "bye"}
