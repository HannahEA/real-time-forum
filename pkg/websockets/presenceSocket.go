package websockets

import (
	"fmt"
	"real-time-forum/pkg/database"
	"sort"
	"time"

	"github.com/gorilla/websocket"
)

type PresenceMessage struct {
	Type      messageType         `json:"type"`
	Timestamp string              `json:"timestamp,omitempty"`
	Username  string              `json:"username,omitempty"`
	Presences []database.Presence `json:"presences,omitempty"`
}

func (m *PresenceMessage) onLoadBroadcast(s *socket) error {

	if s.t == m.Type {
		if err := s.con.WriteJSON(m); err != nil {
			return fmt.Errorf("unable to send (presence) message: %w", err)
		}
	} else {
		return fmt.Errorf("cannot send presence message down ws of type %s", s.t.String())
	}
	return nil
}
func (m *PresenceMessage) Broadcast(conn *websocket.Conn) error {
	if err := conn.WriteJSON(m); err != nil {
		return fmt.Errorf("unable to send (presence) message: %w", err)
	}

	return nil
}

func (m *PresenceMessage) Handle(s map[string][]*websocket.Conn, user *websocket.Conn) error {
	//get username of user logged in on this browser
	var username string
	if m != nil {
		fmt.Println("Presence Username", m.Username)
		fmt.Println("Presence connection", &user)
		username = m.Username
	}
	// get all presences
	presences, err := GetPresences()
	if err != nil {
		return fmt.Errorf("OnPresenceConnect (GetPresences) error: %+v", err)
	}
	fmt.Println("OnPresenceConnect: get presences successful")
	//create message with al presnces to send back to javascript side
	c := &PresenceMessage{
		Type:      presence,
		Timestamp: "",
		Presences: presences,
	}

	// send presences update to all browsers
	var broadErr error
	var oldBrow string
	for k, brow := range s {
		fmt.Println("Browser", brow)
		for i, conn := range brow {
			fmt.Println("Connection", &conn)
			//if connection is not the one where the user is logging in send all presences
			if i == 3 && conn != user && ([]rune(k))[0] >= 48 && ([]rune(k))[0] <= 57 {
				broadErr = c.Broadcast(conn)
				if broadErr != nil {
					return broadErr
				}
				fmt.Println("presences in other browsers updated")
			} else if conn == user {
				fmt.Println("updating presences in current browser...")
				//strore key for browser of currenr user
				oldBrow = k
				//for brwser that users has logged in to, remove logged in user from user list before broadcasting
				for j, p := range c.Presences {
					if p.Nickname == username {
						c2 := &PresenceMessage{}
						users := make([]database.Presence, 0)
						users = append(users, c.Presences[:j]...)
						users = append(users, c.Presences[j+1:]...)
						c2 = c
						c2.Presences = users
						broadErr = c2.Broadcast(conn)
						if broadErr != nil {
							return broadErr
						}
						fmt.Println("presences in current browser updated")
					}

				}
			}

		}
	}
	//change name of browser websocket concections in the map to the username of the logged in user

	BrowserSockets[username] = BrowserSockets[oldBrow]
	delete(BrowserSockets, oldBrow)
	fmt.Println("Browser sockets name updated..", BrowserSockets)
	return broadErr
}

func GetPresences() ([]database.Presence, error) {
	presences := []database.Presence{}
	users, err := database.GetUsers()
	if err != nil {
		return nil, fmt.Errorf("GetUsers (getPresences) error: %+v", err)
	}
	sort.SliceStable(users[:], func(i, j int) bool {
		return users[i].Nickname < users[j].Nickname
	})
	for _, user := range users {
		if user.LoggedIn == "true" {
			presences = append(presences, database.Presence{
				ID:       user.ID,
				Nickname: user.Nickname,
				// Online:            bool,
				// LastContactedTime: created,
			})
		}
	}
	return presences, nil
}

func OnPresenceConnect(s *socket) error {
	time.Sleep(1 * time.Second)
	presences, err := GetPresences()
	if err != nil {
		return fmt.Errorf("OnPresenceConnect (GetPresences) error: %+v", err)
	}
	c := &PresenceMessage{
		Type:      presence,
		Timestamp: "",
		Presences: presences,
	}
	return c.onLoadBroadcast(s)
}

// func (data *Forum) GetSessions() ([]Session, error) {
// 	session := []Session{}
// 	rows, err := data.DB.Query(`SELECT * FROM session`)
// 	if err != nil {
// 		return session, fmt.Errorf("GetSession DB Query error: %+v\n", err)
// 	}
// 	var session_token string
// 	var uName string
// 	var exTime string
// 	for rows.Next() {
// 		err := rows.Scan(&session_token, &uName, &exTime)
// 		if err != nil {
// 			return nil, fmt.Errorf("GetSession rows.Scan error: %+v\n", err)
// 		}
// 		// time.Format("01-02-2006 15:04")
// 		convTime, err := time.Parse("2006-01-02 15:04:05.999999999Z07:00", exTime)
// 		if err != nil {
// 			return nil, fmt.Errorf("GetSession time.Parse(exTime) error: %+v\n", err)
// 		}
// 		session = append(session, Session{
// 			SessionID: session_token,
// 			Nickname:  uName,
// 			Expiry:    convTime,
// 		})
// 	}
// 	return session, nil
// }
