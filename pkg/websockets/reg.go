package websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"real-time-forum/pkg/database"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var User1 database.User

// ***************************REGISTER**********************************************************8
// check if pasword meets criteria number length etc, if nickname is not taken
func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getting data")
	//registration data, unmarshall to struct
	var data database.User
	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
	}

	// User Exists message
	type Exist struct {
		Nickname string
		Email    string
	}
	var Exists Exist

	fmt.Println(data)
	fmt.Println(data.LoggedIn)

	// check database for nickname or email match
	rows, err := database.DB.Query(`SELECT nickname, email FROM users`)
	if err != nil {
		log.Println(err)
	}
	var nickname string
	var email string
	for rows.Next() {
		err := rows.Scan(&nickname, &email)
		if err != nil {
			log.Fatal(err)
		}
		//nickname match
		if data.Nickname == nickname {
			Exists.Nickname = "true"
		}
		//email match
		if data.Email == email {
			Exists.Email = "true"
		}
	}
	// create user if no matches found
	if Exists.Nickname == "" && Exists.Email == "" {
		Exists.Nickname = "false"
		Exists.Email = "false"
		data.Password = passwordHash(data.Password)
		CreateUser(data)
	}
	// send user exists message to javascript
	json.NewEncoder(w).Encode(Exists)
}

func passwordHash(str string) string {
	hashedPw, err := bcrypt.GenerateFromPassword([]byte(str), 8)
	if err != nil {
		log.Fatal(err)
	}
	return string(hashedPw)
}

func checkPwHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// *****************************LOGIN ***************************************
func Login(w http.ResponseWriter, r *http.Request) {
	var user database.Login

	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		log.Println(err)
	}

	var users []database.User

	// selects nickname and password from user database
	rows, err := database.DB.Query(`SELECT nickname, password,loggedin, email FROM users`)
	if err != nil {
		log.Println(err)
	}

	var nickname string
	var password string
	var loggedin string
	var email string

	for rows.Next() {
		err := rows.Scan(&nickname, &password, &loggedin, &email)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println(user.Nickname, nickname, email)

		// compares data with front end, if user nick match, checks pw if match stores value
		if user.Nickname == nickname || user.Nickname == email {
			if checkPwHash(user.Password, password) {
				users = append(users, database.User{
					Nickname: nickname,
					Password: password,
					LoggedIn: "true",
				})
			}

		}
	}

	// if len ==0, no matching user was found
	if len(users) == 0 {
		fmt.Println("pw mismatch")
	}
	fmt.Println(users)

	// checks len again to stop panic err && updates user logged in to true in DB and creates cookie
	if len(users) > 0 && users[0].LoggedIn == "true" {
		User1.Nickname = user.Nickname
		loggedin := "true"

		cookieValue := uuid.NewV4()
		UpdateUser(users[0].Nickname, loggedin, (cookieValue.String()))
		Cookie(w, r, users[0].Nickname, (cookieValue.String()))
		// sore username of logged in user so you can delete cookie on logout

	}
	// sends data to js front end
	json.NewEncoder(w).Encode(users)
}

// updates database table
func UpdateUser(nickname, loggedin string, id string) {
	fmt.Println("Updaate user:", nickname, " loggedin:", loggedin)
	stmt, err := database.DB.Prepare(`UPDATE "users" SET "loggedin" = ? WHERE "nickname" = ?`)
	if err != nil {
		log.Println("Update user:,", err)
	} else {
		stmt.Exec(loggedin, nickname)
	}
	stmt.Close()
	if loggedin == "true" {
		fmt.Println("Adding cookie to database...")
		stmt2, err := database.DB.Prepare("INSERT INTO cookies (name, sessionID) VALUES (?, ?)")
		if err != nil {
			log.Println("Update cookies table: Login:,", err)
		} else {
			stmt2.Exec(nickname, id)
		}
	} else {
		fmt.Println("Removing cookie from database...")
		_, err := database.DB.Exec("DELETE FROM cookies WHERE name = ?", nickname)
		if err != nil {
			log.Println("Update user: Logout:,", err)
		}
	}

}

// *****************************LOGOUT***************************************
type logoutDetails struct {
	Username string `json:"username"`
}

func Logout(w http.ResponseWriter, r *http.Request) {
	var details logoutDetails
	err := json.NewDecoder(r.Body).Decode(&details)
	// name := "test"
	loggedin := "false"
	fmt.Println("Logging out", details.Username)
	UpdateUser(details.Username, loggedin, "")

	var user = &PresenceMessage{
		Username: details.Username,
		Login:    "false",
	}
	if err := user.Handle(); err != nil {
		fmt.Println(err)
		return
	}

	//close websocket connection
	BrowserSockets[details.Username][0].Close()
	BrowserSockets[details.Username][1].Close()
	// delete logged out user from map
	delete(BrowserSockets, details.Username)
	fmt.Println("Browser sockets name updated..", BrowserSockets)

	c1, err := r.Cookie(details.Username)
	fmt.Println("Cookie---", c1)
	if err != nil {
		fmt.Println("On Logout: cookie for user cannot be found!")
	} else {
		fmt.Println("Logout: Cookie deleted, user successfully logged out")
		c1.MaxAge = -1
		http.SetCookie(w, c1)
	}
}

// creates cookie
func Cookie(w http.ResponseWriter, r *http.Request, Username string, id string) {
	// expiration := time.Now().Add(1 * time.Hour)
	cookie := http.Cookie{Name: Username, Value: id, MaxAge: 0}
	// cookie, _ := r.Cookie("username")
	http.SetCookie(w, &cookie)
	fmt.Println(cookie)
	// fmt.Fprintf((w, cookie))
}

func CheckCookies(w http.ResponseWriter, r *http.Request) {
}
