package websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime/debug"
	"time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type SocketMessage struct {
	Type messageType `json:"type,omitempty"`
}

type socket struct {
	con      *websocket.Conn
	nickname string
	t        messageType
	uuid     uuid.UUID
}

var (
	t        = time.Now()
	dateTime = t.Format("1/2/2006, 3:04:05 PM")
	upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	BrowserSockets = make(map[string][]*websocket.Conn)
	SavedSockets   = make([]*websocket.Conn, 0)
)

func SocketCreate(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Socket Request on " + r.RequestURI)
	con, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		fmt.Println(err)
	}
	var ptrSocket = &socket{
		con:  con,
		uuid: uuid.NewV4(),
	}
	SavedSockets = append(SavedSockets, ptrSocket.con)
	fmt.Println("SavedSocket", SavedSockets)
	switch r.RequestURI {
	case "/content":
		ptrSocket.t = content
		// loads the home page (which contains the posts form)
		if err := OnContentConnect(ptrSocket); err != nil {
			fmt.Println(err)
			return
		}
	case "/post":
		ptrSocket.t = post
		// loads the saved posts on window load
		if err := OnPostsConnect(ptrSocket); err != nil {
			fmt.Println(err)
			return
		}
	case "/chat":
		ptrSocket.t = chat
	case "/presence":
		ptrSocket.t = presence
		// loads the presence list on window load
		if err := OnPresenceConnect(ptrSocket); err != nil {
			fmt.Println(err)
			return
		}
	default:
		ptrSocket.t = unknown
	}

	// if len(SavedSockets) == 4 {
	// 	fmt.Println(SavedSockets)
	// 	name := strconv.Itoa(len(BrowserSockets))
	// 	BrowserSockets[name] = SavedSockets
	// 	fmt.Println("Browser Sockets", BrowserSockets)
	// 	var emptySockets []*websocket.Conn
	// 	SavedSockets = emptySockets
	// }
	ptrSocket.pollSocket()

}

func (s *socket) pollSocket() {
	go func() {
		defer func() {
			err := recover()
			if err != nil {
				fmt.Printf("recovered panic in %s socket: %+v\n%s\n", s.t.String(), err, string(debug.Stack()))
			}
		}()
		for {
			b, err := s.read()
			if err != nil {
				fmt.Println(s.t.String() + " socket closed")
				panic(err)
			} else if b == nil {
				fmt.Println(s.t.String() + " socket closed")
				return
			}
			sm := &SocketMessage{}
			if err := json.Unmarshal(b, sm); err != nil {
				panic(err)
			}
			switch sm.Type {
			case chat:
				m := &ChatMessage{}
				if err := json.Unmarshal(b, m); err != nil {
					panic(err)
				}
				if err := m.Handle(s); err != nil {
					panic(err)
				}
			case post:
				m := &PostMessage{}
				if err := json.Unmarshal(b, m); err != nil {
					panic(err)
				}
				if err := m.Handle(s); err != nil {
					panic(err)
				}
			case content:
				m := &ContentMessage{}
				if err := json.Unmarshal(b, m); err != nil {
					panic(err)
				}
				if err := m.Handle(s); err != nil {
					panic(err)
				}
			case presence:
				m := &PresenceMessage{}
				if err := json.Unmarshal(b, m); err != nil {
					panic(err)
				}
				if err := m.Handle(); err != nil {
					panic(err)
				}
			default:
				panic(fmt.Errorf("unable to determine message type for '%s'", string(b)))
			}
		}
	}()
}

func (s *socket) read() ([]byte, error) {
	_, b, err := s.con.ReadMessage()
	if err != nil {
		if websocket.IsCloseError(err, websocket.CloseGoingAway) {
			return nil, nil
		}
		log.Print(b)
		return nil, fmt.Errorf("unable to read message from socket, got: '%s', %w", string(b), err)
	}
	return b, nil
}
