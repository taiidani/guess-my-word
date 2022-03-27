package actions

import (
	"context"
	"guess_my_word/internal/model"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:    1024,
	WriteBufferSize:   1024,
	EnableCompression: true,
	HandshakeTimeout:  time.Second,
	CheckOrigin: func(r *http.Request) bool {

		allowedHosts := []string{
			os.Getenv("ORIGIN"),
			"localhost", "localhost:3000",
			"127.0.0.1", "127.0.0.1:3000",
		}
		for _, host := range allowedHosts {
			if r.Host == host {
				return true
			}
		}

		log.Println(r.Host, "host was denied websocket connection")
		return false
	},
}

var socketPeers = newPubSub()

type stats struct {
	Date     time.Time `form:"date" time_format:"unix"` // The unix date, no timestamp
	dateUser time.Time // The date under the user's timezone
	TZ       int       `form:"tz"`
	Mode     string    `form:"mode"`
}

type statsReply struct {
	Today model.Word `json:"today"`
	Word  model.Word `json:"word"`
	Error string     `json:"error,omitempty"`
}

// ErrRevealToday is emitted when the reveal request is for a current or future word
const ErrRevealToday = "It's too early to reveal this word. Please try again later!"

func wsHandler(c *gin.Context) {
	statsRequest := stats{}

	// Validate the request
	if err := c.ShouldBind(&statsRequest); err != nil {
		log.Println("Invalid request received: ", err)
		c.JSON(400, statsReply{Error: ErrInvalidRequest})
		return
	}

	wsStatsHandler(c, statsRequest)
}

func wsStatsHandler(c *gin.Context, request stats) {
	reply := statsReply{}

	if request.Date.Unix() == 0 {
		reply.Error = ErrInvalidStartTime
	} else {
		log.Println("TZ:", request.TZ)
		request.dateUser = convertUTCToUser(request.Date, request.TZ)
		log.Printf("Requested date is: %s", request.dateUser)

		y, m, d := time.Now().Date()
		cmp := time.Date(y, m, d, 0, 0, 0, 0, request.dateUser.Location())

		if request.dateUser.After(cmp) {
			log.Printf("Compared date was: %s", request.dateUser)
			reply.Error = ErrRevealToday
		}
	}

	if reply.Error != "" {
		c.JSON(400, reply)
		return
	}

	// Begin sending replies on a websocket connection
	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("Failed to set websocket upgrade: %+v", err)
		c.JSON(400, reply)
		return
	}
	defer conn.Close()
	log.Println("Opened websocket connection from", conn.RemoteAddr())

	// Spin up a no-op reader. This allows close connections to come in from the peer but
	// otherwise doesn't handle any messages
	peerCtx, peerCancel := context.WithCancel(c)
	go wsNoOpReaderLoop(conn, peerCancel)

	// Add to the list of subscribers that will be notified whenever a refresh is required
	refresher := socketPeers.Subscribe(conn)

	conn.WriteJSON(wsStatsHandlerLoop(c, request))
	for {
		select {
		case <-c.Request.Context().Done():
			log.Println("WebSocket context canceled: ", c.Request.Context().Err())
			return
		case <-peerCtx.Done():
			log.Printf("WebSocket client %s closed connection", conn.RemoteAddr())
			return
		case <-refresher:
			log.Println("Refreshing WebSocket client", conn.RemoteAddr())

			conn.SetWriteDeadline(time.Now().Add(time.Second * 5))
			err := conn.WriteJSON(wsStatsHandlerLoop(c, request))
			if err != nil {
				log.Printf("Websocket message send error: %s", err)
				return
			}
		}
	}
}

func wsNoOpReaderLoop(c *websocket.Conn, cancel func()) {
	defer cancel()
	for {
		if _, _, err := c.NextReader(); err != nil {
			c.Close()
			break
		}
	}
}

func wsStatsHandlerLoop(ctx context.Context, request stats) statsReply {
	reply := statsReply{}

	// Generate the word for the day
	word, err := wordStore.GetForDay(ctx, request.dateUser, request.Mode)
	if err != nil {
		reply.Error = err.Error()
		return reply
	}

	reply.Word = word

	// Now for today's stats. Similar, but without the word information!
	todayTm := request.dateUser.AddDate(0, 0, 1)
	word, err = wordStore.GetForDay(ctx, todayTm, request.Mode)
	if err != nil {
		reply.Error = err.Error()
		return reply
	}

	reply.Today = word
	reply.Today.Value = ""
	return reply
}

type pubSub struct {
	subscribers map[*websocket.Conn]chan interface{}
}

func newPubSub() *pubSub {
	return &pubSub{
		subscribers: map[*websocket.Conn]chan interface{}{},
	}
}

func (r *pubSub) Notify() {
	for _, ch := range r.subscribers {
		ch <- true
	}
}

func (r *pubSub) Subscribe(conn *websocket.Conn) chan interface{} {
	ret := make(chan interface{}, 10)
	r.subscribers[conn] = ret
	return ret
}

func (r *pubSub) Unsubscribe(conn *websocket.Conn) {
	delete(r.subscribers, conn)
}
