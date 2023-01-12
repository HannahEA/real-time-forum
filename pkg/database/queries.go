package database

import (
	"fmt"
)

type User struct {
	ID        string `json:"id,omitempty"`
	Nickname  string `json:"nickname,omitempty"`
	Nickname2 string `json:"nickname2,omitempty"`
	Age       string `json:"age,omitempty"`
	Gender    string `json:"gender,omitempty"`
	FirstName string `json:"firstname,omitempty"`
	LastName  string `json:"lastname,omitempty"`
	Email     string `json:"email,omitempty"`
	Password  string `json:"password,omitempty"`
	LoggedIn  string `json:"loggedin,omitempty"`
}
type Post struct {
	PostID     string    `json:"post_id,omitempty"`
	Nickname   string    `json:"nickname"`
	Title      string    `json:"title,omitempty"`
	Categories string    `json:"categories,omitempty"`
	Body       string    `json:"body,omitempty"`
	Comments   []Comment `json:"comments,omitempty"`
}
type Comment struct {
	CommentID string `json:"comment_id,omitempty"`
	PostID    string `json:"post_id,omitempty"`
	Nickname  string `json:"nickname"`
	Body      string `json:"body,omitempty"`
}
type ConvoParticipants struct {
	ConvoID      string `json:"convo_id"`
	Participant1 string `json:"user1"`
	Participant2 string `json:"user2"`
}

type Conversation struct {
	ConvoID      string `json:"convo_id"`
	Participants []User `json:"participants"`
	Chats        []Chat `json:"chats"`
}
type Chat struct {
	ConvoID string `json:"convo_id"`
	ChatID  string `json:"chat_id"`
	Sender  User   `json:"sender"`
	Date    string `json:"date,omitempty"`
	Body    string `json:"body,omitempty"`
}
type Presence struct {
	ID                string `json:"id"`
	Nickname          string `json:"nickname"`
	Online            string `json:"online"`
	LastContactedTime string `json:"last_contacted_time"`
}

type Login struct {
	Nickname string `json:"nickname,omitempty"`
	Password string `json:"password,omitempty"`
}


func GetUsers() ([]User, error) {
	users := []User{}
	rows, err := DB.Query(`SELECT * FROM users`)
	if err != nil {
		return users, fmt.Errorf("GetUsers DB Query error: %+v\n", err)
	}
	var id string
	var nickname string
	var age string
	var gender string
	var firstname string
	var lastname string
	var email string
	var password string
	var loggedin string

	for rows.Next() {
		err := rows.Scan(&id, &nickname, &age, &gender, &firstname, &lastname, &email, &password, &loggedin)
		if err != nil {
			return users, fmt.Errorf("GetUsers rows.Scan error: %+v\n", err)
		}
		users = append(users, User{
			ID:        id,
			Nickname:  nickname,
			Age:       age,
			Gender:    gender,
			FirstName: firstname,
			LastName:  lastname,
			Email:     email,
			Password:  password,
			LoggedIn:  loggedin,
		})
	}
	err = rows.Err()
	if err != nil {
		return users, err
	}
	return users, nil
}
func GetPosts() ([]*Post, error) {
	posts := []*Post{}
	rows, err := DB.Query(`SELECT * FROM posts`)
	if err != nil {
		return posts, fmt.Errorf("GetPosts DB Query error: %+v\n", err)
	}
	var postid string
	var nickname string
	var title string
	var category string
	var postcontent string
	for rows.Next() {
		err := rows.Scan(&postid, &nickname, &title, &category, &postcontent)
		if err != nil {
			return posts, fmt.Errorf("GetPosts rows.Scan error: %+v\n", err)
		}
		posts = append(posts, &Post{
			PostID:     postid,
			Nickname:   nickname,
			Title:      title,
			Categories: category,
			Body:       postcontent,
		})
	}
	err = rows.Err()
	if err != nil {
		return posts, err
	}
	return posts, nil
}
func GetPostForComment(c Comment) (Post, error) {
	posts, err := GetPosts()
	if err != nil {
		return Post{}, err
	}
	for _, p := range posts {
		if p.PostID == c.PostID {
			return *p, nil
		}
	}
	return Post{}, fmt.Errorf("no matching post found for id: %s", c.PostID)
}
func GetPopulatedPosts() ([]*Post, error) {
	posts, err := GetPosts()
	if err != nil {
		return nil, fmt.Errorf("GetPopulatedPosts (GetPosts) error: %+v\n", err)
	}
	populatedPosts, err := populateCommentsForPosts(posts)
	if err != nil {
		return nil, fmt.Errorf("GetPopulatedPosts (populateCommentsForPosts) error: %+v\n", err)
	}
	return populatedPosts, nil
}
func populateCommentsForPosts(posts []*Post) ([]*Post, error) {
	comments, err := GetComments()
	if err != nil {
		return nil, fmt.Errorf("populateCommentsForPosts (GetComments) error: %+v\n", err)
	}
	outPost := []*Post{}
	for _, pts := range posts {
		newPost := pts
		for _, cmts := range comments {
			if pts.PostID == cmts.PostID {
				newPost.Comments = append(newPost.Comments, cmts)
			}
		}
		outPost = append(outPost, newPost)
	}
	return outPost, nil
}
func GetComments() ([]Comment, error) {
	comments := []Comment{}
	rows, err := DB.Query(`SELECT * FROM comments`)
	if err != nil {
		return comments, fmt.Errorf("GetComments DB Query error: %+v\n", err)
	}
	var postid string
	var commentid string
	var nickname string
	var commentcontent string
	for rows.Next() {
		err := rows.Scan(&commentid, &postid, &nickname, &commentcontent)
		if err != nil {
			return comments, fmt.Errorf("GetComments rows.Scan error: %+v\n", err)
		}
		comments = append(comments, Comment{
			CommentID: commentid,
			PostID:    postid,
			Nickname:  nickname,
			Body:      commentcontent,
		})
	}
	err = rows.Err()
	if err != nil {
		return comments, err
	}
	return comments, nil
}
func FilterCommentsForPost(postID string, comments []Comment) []Comment {
	out := []Comment{}
	for _, c := range comments {
		if postID == c.PostID {
			out = append(out, c)
		}
	}
	return out
}

func populateUsersForChats(chats []Chat) ([]Chat, error) {
	users, err := GetUsers()
	if err != nil {
		return nil, fmt.Errorf("populateUsersForChats (GetUsers) error: %+v\n", err)
	}
	outChats := []Chat{}
	for _, chat := range chats {
		newChat := chat
		for _, user := range users {
			if chat.Sender.ID == user.ID {
				newChat.Sender = user
			}
		}
		outChats = append(outChats, newChat)
	}
	return outChats, nil
}


func GetChat(convoID string) ([]Chat, error) {
	chats := []Chat{}
	rows, err := DB.Query(`SELECT * FROM chats WHERE convoID=?`, convoID)
	if err != nil {
		return chats, fmt.Errorf("GetChats DB Query error: %+v\n", err)
	}
	var convoid string
	var chatid string
	var sender string
	var date string
	var body string
	for rows.Next() {
		err := rows.Scan(&convoid, &chatid, &sender, &date, &body)
		if err != nil {
			return chats, fmt.Errorf("GetChats rows.Scan error: %+v\n", err)
		}
		chats = append(chats, Chat{
			ConvoID: convoid,
			ChatID:  chatid,
			Sender: User{
				ID: sender,
			},
			Date: date,
			Body: body,
		})
	}
	err = rows.Err()
	if err != nil {
		return chats, err
	}
	return populateUsersForChats(chats)
}

func GetUserFromConversations(nickname, nickname2 string) (ConvoParticipants, error) {
	user := ConvoParticipants{}

	rows, err := DB.Query(`SELECT * FROM conversations WHERE participants =? AND participants2 =?`, nickname, nickname2)
	if err != nil {
		return user, fmt.Errorf("GetUsers DB Query error: %+v\n", err)
	}
	var convoId string
	var participant1 string
	var participant2 string
	for rows.Next() {
		err := rows.Scan(&convoId, &participant1, &participant2)
		if err != nil {
			return user, fmt.Errorf("GetUsers rows.Scan error: %+v\n", err)
		}
		user = ConvoParticipants{
			ConvoID:      convoId,
			Participant1: participant1,
			Participant2: participant2,
		}
	}
	err = rows.Err()
	if err != nil {
		return user, err
	}
	fmt.Println(user, "users part of convo")
	return user, nil
}
