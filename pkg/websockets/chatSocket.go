package websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"real-time-forum/pkg/database"
	newTime "time"

	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
)

type ChatMessage struct {
	Type          messageType              `json:"type,omitempty"`
	Timestamp     string                   `json:"timestamp,omitempty"`
	Conversations []*database.Conversation `json:"conversations"`
}

func GetChats(w http.ResponseWriter, r *http.Request) {
	var participants database.ConvoParticipants
	err := json.NewDecoder(r.Body).Decode(&participants)
	if err != nil {
		log.Println(err)
	}
	log.Println("data asked")
	log.Println(participants)
	var usersMatch, jsonChats = usersPartofConvo(participants.Participant1, participants.Participant2, false)

	if usersMatch {
		w.Write(jsonChats)
	}
}

func Notification(w http.ResponseWriter, r *http.Request) {
	// adding or removing notification?
	var chat database.Chat
	err := json.NewDecoder(r.Body).Decode(&chat)
	if err != nil {
		log.Println(err)
	}
	fmt.Println("Notification change")
	if chat.Notification {
		fmt.Println("adding Notification", chat.Sender.ID, chat.Reciever.ID)
		CreateNotification(chat.Sender.ID, chat.Reciever.ID)
	} else {
		fmt.Println("removing Notification")
		RemoveNotification(chat.Reciever.ID, chat.Sender.ID)
	}
	data, _ := json.Marshal(chat)
	w.Write(data)
}

func usersPartofConvo(user1, user2 string, newMessage bool) (bool, []byte) {

	// var participant1 = convo.Participants[0]
	// var p1 = fmt.Sprintf("%+v", participant1.ID)
	// var participant2 = convo.Participants[1]
	var convoCheck, _ = database.GetUserFromConversations(user1, user2)
	fmt.Println("convo check", convoCheck)
	fmt.Println(convoCheck.Participant1, convoCheck.Participant2, user1, user2)
	if convoCheck.Participant1 == user1 && convoCheck.Participant2 == user2 {
		var chats, _ = database.GetChat(convoCheck.ConvoID)
		var jsondata []byte
		if newMessage {
			jsondata, _ = json.Marshal(chats[len(chats)-1])
			fmt.Println("found new chat", chats[len(chats)-1])
		} else {
			jsondata, _ = json.Marshal(chats)
		}

		return true, jsondata
	}
	return false, nil
}
func CreateNotification(user string, from string) {
	stmt, err := database.DB.Prepare("INSERT INTO notificationss (user, user2) VALUES (?, ?);")
	defer stmt.Close()
	if err != nil {
		fmt.Printf("CreateNotification DB Prepare error: %+v\n", err)
	}

	_, err = stmt.Exec(user, from)
	if err != nil {
		fmt.Printf("CreateNotifictaion Exec error: %+v\n", err)

	}

}
func RemoveNotification(user string, from string) {
	_, err := database.DB.Exec("DELETE FROM notificationss WHERE user = ? AND user2 = ?", user, from)
	if err != nil {
		log.Println("Remove Notification Database Error:,", err)
	}

}

func (m *ChatMessage) Handle(s *socket) error {

	fmt.Println("chat message func", m.Conversations, "type", m.Type)
	fmt.Println("time after this", m.Timestamp, "chat")
	time := newTime.Now().String()
	fmt.Println("time in golang", time)
	var convoCheckID string
	var reciever string
	//find last sent message before new chat or convo is added to DB
	lastLatestSender := LatestChatConvo(m.Conversations[0].Participants[1].ID)
	fmt.Println("Function worked: found the latest user:", lastLatestSender)

	//Create new chat and convo if needed
	for i, convo := range m.Conversations {
		var participant1 = convo.Participants[0]
		// var p1 = fmt.Sprintf("%+v", participant1.ID)
		var participant2 = convo.Participants[1]
		reciever = convo.Participants[1].ID
		var partOfConvo, _ = usersPartofConvo(participant1.ID, participant2.ID, true)
		var convoCheck, _ = database.GetUserFromConversations(participant1.ID, participant2.ID)
		if partOfConvo {
			convo.ConvoID = convoCheck.ConvoID
			convoCheckID = convoCheck.ConvoID
		}
		// creates a new conversation if the convoID is missing
		if convo.ConvoID == "" {
			newConvoID, err := CreateConversation(convo)
			if err != nil {
				return fmt.Errorf("ChatSocket Handle (CreateConversation) error: %w", err)
			}
			convo.ConvoID = newConvoID
			convoCheckID = newConvoID
		}
		for j, chat := range convo.Chats {
			chat.Date = time
			// for new chats, the chat.ConvoID is given the conversation's convoID if it is missing
			if chat.ConvoID == "" {
				chat.ConvoID = convo.ConvoID
			}
			if chat.ChatID == "" {
				newChatID, err := CreateChat(chat)
				if err != nil {
					return fmt.Errorf("ChatSocket Handle (CreateChat) error: %w", err)
				}
				chat.ChatID = newChatID
			}
			convo.Chats[j] = chat
		}
		m.Conversations[i] = convo
	}
	//get latest chat from database
	var chats, _ = database.GetChat(convoCheckID)
	var newChat = &chats[len(chats)-1]
	//add timestamp of latest message to convo DB
	AddDatetoConvoDB(m.Timestamp, m.Conversations[0].Participants[0].ID, m.Conversations[0].Participants[1].ID)

	//if the last person to send this user a message is not the same as the current chat sender then the presence list needs to be updated
	if m.Conversations[0].Participants[0].ID != lastLatestSender {
		//call presence handler
		m := &PresenceMessage{
			// reciever
			Username: m.Conversations[0].Participants[1].ID,
			Login:    "chat",
		}
		err := m.Handle()
		if err != nil {
			fmt.Print("Chat Handler: Presence Handler Error: ")
			return err
		}
	}
	// send to sender
	err2 := Broadcast(s.con, newChat)
	if err2 != nil {
		fmt.Print("Chat Handler: unable to broadcast new message to sender")
		return err2
	}
	fmt.Println("Chat Handler: broadcast new message to sender")

	//send to reciever
	if BrowserSockets[reciever] == nil {
		fmt.Println("Broser Socket", reciever, BrowserSockets[reciever])
		// notification when receiver is logged out
		CreateNotification(m.Conversations[0].Participants[0].ID, m.Conversations[0].Participants[1].ID)
	} else {
		err3 := Broadcast(BrowserSockets[reciever][0], newChat)
		if err3 != nil {
			fmt.Print("Chat Handler: unable to broadcast new message to reciever", reciever)
			return err3
		}
		fmt.Println("Chat Handler: broadcast new message to reciever", reciever)
	}
	//change convo DB: update latest message time

	return nil
}

func Broadcast(s *websocket.Conn, m *database.Chat) error {
	// if s.t == m.Type {
	if err := s.WriteJSON(m); err != nil {
		return fmt.Errorf("unable to send (chat) message: %w", err)
	}
	// } else {
	// 	return fmt.Errorf("cannot send chat message down ws of type %s", s.t.String())
	// }
	return nil
}
func CreateChat(chat database.Chat) (string, error) {
	stmt, err := database.DB.Prepare("INSERT INTO chats (convoID, chatID, sender, date, body) VALUES (?, ?, ?, ?, ?);")
	if err != nil {
		return "", fmt.Errorf("CreateChat DB Prepare error: %+v\n", err)
	}
	defer stmt.Close()
	if chat.ChatID == "" {
		chat.ChatID = uuid.NewV4().String()
	}

	_, err = stmt.Exec(chat.ConvoID, chat.ChatID, chat.Sender.ID, chat.Date, chat.Body)
	if err != nil {
		return "", fmt.Errorf("CreateChat Exec error: %+v\n", err)
	}
	return chat.ChatID, err
}

func CreateConversation(conversations *database.Conversation) (string, error) {
	log.Println("inserting into convo db")
	stmt, err := database.DB.Prepare("INSERT INTO conversations (convoID, participants, participants2, lastMessageTime) VALUES (?, ?,?,?);")
	defer stmt.Close()
	if err != nil {
		return "", fmt.Errorf("CreateConversations DB Prepare error: %+v\n", err)
	}
	if conversations.ConvoID == "" {
		conversations.ConvoID = uuid.NewV4().String()
	}

	_, err = stmt.Exec(conversations.ConvoID, conversations.Participants[0].ID, conversations.Participants[1].ID, "")
	if err != nil {
		return "", fmt.Errorf("CreateConversations Exec error: %+v\n", err)

	}
	_, err2 := stmt.Exec(conversations.ConvoID, conversations.Participants[1].ID, conversations.Participants[0].ID, "")
	if err2 != nil {
		return "", fmt.Errorf("CreateConversations Exec error: %+v\n", err2)

	}
	return conversations.ConvoID, err
}

// add the time of the last chat sent in this conversation to the convo database
func AddDatetoConvoDB(date string, sender string, reciever string) {
	stmt, err := database.DB.Prepare(`UPDATE "conversations" SET "lastMessageTime" = ? WHERE "participants" = ? AND "participants2" = ? OR "participants" = ? AND "participants2" = ?`)
	defer stmt.Close()
	if err != nil {
		fmt.Printf("AddDatetoConvoDB: DB Prepare Error:%+v\n", err)
	}
	_, err2 := stmt.Exec(date, sender, reciever, reciever, sender)
	if err2 != nil {
		fmt.Printf("AddDatetoConvoDB: DB Exec Error:%+v\n", err)
	}
}

// who was the last person to send this user a message
func LatestChatConvo(user string) string {
	rows, err := database.DB.Query("SELECT * FROM conversations WHERE participants = ?", user)
	if err != nil {
		fmt.Printf("LatestChatConvo: DB Query Error:%+v\n", err)
	}
	var convoId string
	var participant1 string
	var participant2 string
	var time string
	var latestTime newTime.Time
	var latestUser string
	for rows.Next() {
		fmt.Println("rows latest user:", rows)
		scanErr := rows.Scan(&convoId, &participant1, &participant2, &time)
		if scanErr != nil {
			fmt.Printf("LatestChatConvo: Scan Error:%+v\n", err)
		}
		// time = strings.Replace(time, ",", "", 1)
		// fullTimeArr := strings.Split(time, " ")
		// //reformat date
		// dateTimeArr := strings.Split(fullTimeArr[0], "/")
		// dateTimeArr[2], dateTimeArr[0] = dateTimeArr[0], dateTimeArr[2]
		// fullTimeArr[0] = strings.Join(dateTimeArr, "-")
		// // reformat time
		// timeArr := strings.Split(fullTimeArr[1], ":")
		// timeArr = timeArr[:2]
		// fullTimeArr[1] = strings.Join(timeArr, ":")
		// //rejoin date and time
		// time = strings.Join(fullTimeArr, " ")
		// fmt.Println("time", time)
		times, err := newTime.Parse("2006-01-02 15:04", time)
		if err != nil {
			fmt.Printf("Error LatestChatConvo: time.Parse %v\n", err)
		}
		fmt.Println("latest message time", latestTime)
		fmt.Println("conversation time", times)
		if times.After(latestTime) {

			latestTime = times
			latestUser = participant2
		}
	}
	rowErr := rows.Err()
	if rowErr != nil {
		fmt.Printf("LatestChatConvo: Rows Loop Error:%+v\n", err)
	}
	return latestUser
}
