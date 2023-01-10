package websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"real-time-forum/pkg/database"

	uuid "github.com/satori/go.uuid"
)

type ChatMessage struct {
	Type          messageType              `json:"type,omitempty"`
	Timestamp     string                   `json:"timestamp,omitempty"`
	Conversations []*database.Conversation `json:"conversations"`
}

func (m *ChatMessage) Broadcast(s *socket) error {
	if s.t == m.Type {
		if err := s.con.WriteJSON(m); err != nil {
			return fmt.Errorf("unable to send (chat) message: %w", err)
		}
	} else {
		return fmt.Errorf("cannot send chat message down ws of type %s", s.t.String())
	}
	return nil
}
func (m *ChatMessage) Handle(s *socket) error {
	fmt.Println("chat message func", m.Conversations, "type", m.Type)
	fmt.Println("time after this", m.Timestamp, "chat")
	var time = m.Timestamp
	if len(m.Conversations) == 0 {
		fmt.Println("no converation to be handled..")
		conversations, err := database.GetPopulatedConversations(nil)
		if err != nil {
			return err
		}
		c := &ChatMessage{
			Type:          chat,
			Conversations: conversations,
		}
		return c.Broadcast(s)
	}
	for i, convo := range m.Conversations {
		var participant1 = convo.Participants[0]
		// var p1 = fmt.Sprintf("%+v", participant1.ID)
		var participant2 = convo.Participants[1]
		var partOfConvo, _ = usersPartofConvo(participant1.ID, participant2.ID)
		var convoCheck, _ = database.GetUserFromConversations(participant1.ID, participant2.ID)

		if partOfConvo {
			convo.ConvoID = convoCheck.ConvoID
		}
		// creates a new conversation if the convoID is missing
		if convo.ConvoID == "" {
			newConvoID, err := CreateConversation(convo)
			if err != nil {
				return fmt.Errorf("ChatSocket Handle (CreateConversation) error: %w", err)
			}
			convo.ConvoID = newConvoID
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
	fmt.Println("chat handler: new convo created", m.Conversations)
	c, err := database.GetPopulatedConversations(m.Conversations)
	if err != nil {
		return fmt.Errorf("ChatSocket Handle (GetPopulatedConversations) error: %w", err)
	}
	m.Conversations = c
	// b, _ := json.Marshal(m.Conversations)
	// fmt.Println(string(b))
	return m.Broadcast(s)
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
	// TODO: remove placeholder nickname once login/sessions are working
	if chat.Sender.ID == "" {
		//this is foo's userID in the database
		chat.Sender.ID = "6d01e668-2642-4e55-af73-46f057b731f9"
	}
	_, err = stmt.Exec(chat.ConvoID, chat.ChatID, chat.Sender.ID, chat.Date, chat.Body)
	if err != nil {
		return "", fmt.Errorf("CreateChat Exec error: %+v\n", err)
	}
	return chat.ChatID, err
}
func CreateConversation(conversations *database.Conversation) (string, error) {
	log.Println("inserting into convo db")
	stmt, err := database.DB.Prepare("INSERT INTO conversations (convoID, participants, participants2) VALUES (?, ?,?);")
	defer stmt.Close()
	if err != nil {
		return "", fmt.Errorf("CreateConversations DB Prepare error: %+v\n", err)
	}
	if conversations.ConvoID == "" {
		conversations.ConvoID = uuid.NewV4().String()
	}

	_, err = stmt.Exec(conversations.ConvoID, conversations.Participants[0].ID, conversations.Participants[1].ID)
	if err != nil {
		return "", fmt.Errorf("CreateConversations Exec error: %+v\n", err)

	}
	return conversations.ConvoID, err
}

func GetChats(w http.ResponseWriter, r *http.Request) {
	var participants database.ConvoParticipants
	err := json.NewDecoder(r.Body).Decode(&participants)
	if err != nil {
		log.Println(err)
	}
	log.Println("data asked")
	log.Println(participants)
	var usersMatch, jsonChats = usersPartofConvo(participants.Participant1, participants.Participant2)
	//	 if usersPartofConvo(participants.Participant1, participants.Participant2){
	//		fmt.Println("users match")
	//	 }
	if usersMatch {
		w.Write(jsonChats)
	}
}

func usersPartofConvo(user1, user2 string) (bool, []byte) {

	// var participant1 = convo.Participants[0]
	// var p1 = fmt.Sprintf("%+v", participant1.ID)
	// var participant2 = convo.Participants[1]
	var convoCheck, _ = database.GetUserFromConversations(user1, user2)
	fmt.Println("convo check", convoCheck)
	fmt.Println(convoCheck.Participant1, convoCheck.Participant2, user1, user2)
	if convoCheck.Participant1 == user1 && convoCheck.Participant2 == user2 {
		var chats, _ = database.GetChat(convoCheck.ConvoID)
		// fmt.Println(chats)
		var jsondata, _ = json.Marshal(chats)
		return true, jsondata
		// fmt.Println(convoCheck.Participant1 == participant1.ID, convoCheck.Participant2 ==participant2.ID)
		// convo.ConvoID = convoCheck.ConvoID
	}
	return false, nil
}
