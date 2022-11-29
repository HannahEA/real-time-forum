package websockets

import (
	"fmt"
	"real-time-forum/pkg/database"
	"sort"
	"time"
)

type PresenceMessage struct {
	Type      messageType         `json:"type"`
	Timestamp string              `json:"timestamp,omitempty"`
	Presences []database.Presence `json:"presences,omitempty"`
}

func (m *PresenceMessage) Broadcast(s *socket) error {
	if s.t == m.Type {
		if err := s.con.WriteJSON(m); err != nil {
			return fmt.Errorf("unable to send (presence) message: %w", err)
		}
	} else {
		return fmt.Errorf("cannot send presence message down ws of type %s", s.t.String())
	}
	return nil
}

func (m *PresenceMessage) Handle(s *socket) error {
	presences, err := GetPresences()
	if err != nil {
		return fmt.Errorf("OnPresenceConnect (GetPresences) error: %+v", err)
	}
	fmt.Println("OnPresenceConnect: get presences successful")
	c := &PresenceMessage{
		Type:      presence,
		Timestamp: "",
		Presences: presences,
	}
	return c.Broadcast(s)
	// return m.Broadcast(s)
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
	return c.Broadcast(s)
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
