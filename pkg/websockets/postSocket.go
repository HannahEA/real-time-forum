package websockets
import (
	"fmt"
	"real-time-forum/pkg/database"
	"time"
	uuid "github.com/satori/go.uuid"
)
type PostMessage struct {
	Type      messageType      `json:"type,omitempty"`
	Timestamp string           `json:"timestamp,omitempty"`
	Posts     []*database.Post `json:"posts"`
}
func (m PostMessage) Handle(s *socket) error {
	if len(m.Posts) == 0 {
		p, err := database.GetPopulatedPosts()
		if err != nil {
			return err
		}
		c := &PostMessage{
			Type:  post,
			Posts: p,
		}
		return c.Broadcast(s)
	}
	for _, pt := range m.Posts {
		if pt.PostID == "" {
			newPost, err := CreatePost(pt)
			if err != nil {
				return fmt.Errorf("PostSocket Handle (CreatePost) error: %w", err)
			}
			pt = newPost
		}
		for _, comment := range pt.Comments {
			if comment.CommentID == "" {
				newComment, err := CreateComment(comment)
				if err != nil {
					return fmt.Errorf("PostSocket Handle (CreateComment) error: %w", err)
				}
				comment = newComment
			}
		}
		return m.Broadcast(s)
	}
	return nil
}
func (m *PostMessage) Broadcast(s *socket) error {
	if s.t == m.Type {
		if err := s.con.WriteJSON(m); err != nil {
			return fmt.Errorf("unable to send (post) message: %w", err)
		}
	} else {
		return fmt.Errorf("cannot send post message down ws of type %s", s.t.String())
	}
	return nil
}
// TODO: add timestamp in the PostMessage struct; c
func OnPostsConnect(s *socket) error {
	time.Sleep(1 * time.Second)
	p, err := database.GetPopulatedPosts()
	if err != nil {
		return err
	}
	c := &PostMessage{
		Type:      post,
		Timestamp: "",
		Posts:     p,
	}
	return c.Broadcast(s)
}
func CreatePost(post *database.Post) (*database.Post, error) {
	stmt, err := database.DB.Prepare("INSERT INTO posts (postID, nickname, title, categories, body) VALUES (?, ?, ?, ?, ?);")
	defer stmt.Close()
	if err != nil {
		return nil, fmt.Errorf("CreatePost DB Prepare error: %+v\n", err)
	}
	if post.PostID == "" {
		post.PostID = uuid.NewV4().String()
	}
	// TODO: remove placeholder nickname once login/sessions are working, and hook up the real user who is logged in
	// if post.Nickname == "" {
	// 	post.Nickname = "Cassidy"
	// }
	_, err = stmt.Exec(post.PostID, post.Nickname, post.Title, post.Categories, post.Body)
	if err != nil {
		return nil, fmt.Errorf("CreatePost Exec error: %+v\n", err)
	}
	return post, err
}
func CreateComment(comment database.Comment) (database.Comment, error) {
	stmt, err := database.DB.Prepare("INSERT INTO comments (commentID, postID, nickname, body) VALUES (?, ?, ?, ?);")
	defer stmt.Close()
	if err != nil {
		return comment, fmt.Errorf("CreateComment DB Prepare error: %+v\n", err)
	}
	if comment.CommentID == "" {
		comment.CommentID = uuid.NewV4().String()
	}
	// TODO: remove placeholder nickname once login/sessions are working
	// if comment.Nickname == "" {
	// 	comment.Nickname = "Cassidy"
	// }
	_, err = stmt.Exec(comment.CommentID, comment.PostID, comment.Nickname, comment.Body)
	if err != nil {
		return comment, fmt.Errorf("CreateComment Exec error: %+v\n", err)
	}
	return comment, err
}
func CreateUser (user database.User) (database.User, error){
	stmt, err := database.DB.Prepare("INSERT INTO users (ID, nickname, age, gender, firstname, lastname, email, password, loggedin) VALUES (?,?,?,?,?,?,?,?,?)")
	defer stmt.Close()
	if err !=nil {
		return user, fmt.Errorf("Create User DB Prepare error: %+v\n", err)
	}
	if user.ID == ""{
		user.ID = uuid.NewV4().String()
	}
	_,err = stmt.Exec(user.ID, user.Nickname, user.Age, user.Gender, user.FirstName, user.LastName, user.Email, user.Password, user.LoggedIn)
	if err != nil {
		return user, fmt.Errorf("Create User Exec error: %+v\n", err)
	}
	return user,err
}

