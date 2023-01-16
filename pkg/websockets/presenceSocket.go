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
	Login     string              `json:"Login,omitempty"`
	Presences []database.Presence `json:"presences,omitempty"`
}

func (m *PresenceMessage) Broadcast(conn *websocket.Conn) error {
	fmt.Println("presence message", m)
	if err := conn.WriteJSON(m); err != nil {
		return fmt.Errorf("unable to send (presence) message: %w", err)
	}

	return nil
}
func changePresList(pres []database.Presence, username string, conn *websocket.Conn) *PresenceMessage {
	c2 := &PresenceMessage{Type: presence,
		Timestamp: ""}

	for j, p := range pres {
		//browser you are logging in on
		if p.Nickname == username {
			if len(pres) == 1 {
				users := make([]database.Presence, 0)
				c2.Presences = users
				return c2
			}
			users := make([]database.Presence, 0)
			users = append(users, pres[:j]...)
			users = append(users, pres[j+1:]...)
			c2.Presences = users
			return c2

		}
	}
	//browser you are not logging in too
	c2.Presences = pres
	fmt.Println("other browser presence", c2.Presences)
	return c2
}

func (m *PresenceMessage) Handle() error {
	//get username of user logged in on this browser
	var username string
	var login bool

	if m != nil {
		fmt.Println("Presence Username", m.Username)
		username = m.Username
	}
	// get all presences
	presences, err := GetPresences()
	if err != nil {
		return fmt.Errorf("OnPresenceConnect (GetPresences) error: %+v", err)
	}
	fmt.Println("OnPresenceConnect: get presences successful", presences)

	var broadErr error
	fmt.Println("username", username, "sockets", BrowserSockets[username])
	if m.Login == "false" {
		fmt.Println("Browser Logout")
		//loggin out
		login = false
	} else {
		//login or browser refresh
		if BrowserSockets[username] != nil {
			fmt.Println("Browser refresh")
		} else {
			fmt.Println("Browser Login")

		}
		var newUser = make([]*websocket.Conn, 0)
		newUser = append(newUser, SavedSockets[len(SavedSockets)-2:]...)
		BrowserSockets[username] = newUser
		fmt.Println("Browser Sockets", BrowserSockets)
		login = true
	}

	for name, brow := range BrowserSockets {
		fmt.Println("Browser", name, brow)

		for i, conn := range brow {
			if i == 1 {
				if (name == username && login) || name != username {
					fmt.Println("updating presences...")
					//for browser that users have logged in to, remove logged in user from user list before broadcasting
					finalM := changePresList(presences, name, conn)
					if name == username {
						finalM.Presences = checkNotifications(username, finalM.Presences)
					}
					broadErr = finalM.Broadcast(conn)
					if broadErr != nil {
						return broadErr
					}
					fmt.Println("presences in current browser updated")
				}

			}

		}

	}

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

		presences = append(presences, database.Presence{
			ID:       user.ID,
			Nickname: user.Nickname,
			Online:   user.LoggedIn,
		})

	}
	return presences, nil
}

func OnPresenceConnect(s *socket) error {
	time.Sleep(1 * time.Second)
	presences := make([]database.Presence, 0)

	c := &PresenceMessage{
		Type:      presence,
		Timestamp: "",
		Presences: presences,
	}
	return c.Broadcast(s.con)
}

func checkNotifications(username string, presences []database.Presence) []database.Presence {
	var Users []string
	rows, err := database.DB.Query(`SELECT * FROM notificationss`)
	if err != nil {
		fmt.Printf("checkNotifictaions db.query error: %+v\n", err)
	}
	var user string
	var user2 string
	for rows.Next() {
		err := rows.Scan(&user, &user2)
		if err != nil {
			fmt.Printf("checkNotifictaions rows.Scan error: %+v\n", err)
		}
		if user2 == username {
			Users = append(Users, user)
		}
	}
	fmt.Println("presence with notifictaion", Users)

	for i, pres := range presences {
		for _, user := range Users {
			if pres.Nickname == user {
				presences[i].Notification = true
			}
		}
	}
	fmt.Println("New presences", presences)

	return presences

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
