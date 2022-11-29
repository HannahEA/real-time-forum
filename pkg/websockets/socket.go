package websockets
import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	auth "real-time-forum/pkg/authentication"
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
	savedSockets = make([]*socket, 0)
)
func SocketCreate(w http.ResponseWriter, r *http.Request) {
	c1, err1 := r.Cookie("1st-cookie")
	if err1 == nil && !auth.Person.Accesslevel {
		// first home page access
		c1.MaxAge = -1
		http.SetCookie(w, c1)
	}
	_, err := r.Cookie("1st-cookie")
	if err != nil && auth.Person.Accesslevel {
		// logged in and on 2nd browser
		auth.Person.CookieChecker = false
	} else if err == nil && auth.Person.Accesslevel {
		// Original browser and logged in
		auth.Person.CookieChecker = true
	} else {
		// not logged in yet
		auth.Person.CookieChecker = false
	}
	fmt.Println("Socket Request on " + r.RequestURI)
	con, _ := upgrader.Upgrade(w, r, nil)
	ptrSocket := &socket{
		con:  con,
		uuid: uuid.NewV4(),
	}
	// add new case here when added to main.go for handlers
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
	savedSockets = append(savedSockets, ptrSocket)
	ptrSocket.pollSocket()
	for i, so := range savedSockets {
		if so.uuid == ptrSocket.uuid {
			ret := make([]*socket, 0)
			ret = append(ret, savedSockets[:i]...)
			savedSockets = append(ret, savedSockets[i+1:]...)
		}
	}
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
				if err := m.Handle(s); err != nil {
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

