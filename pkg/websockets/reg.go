package websockets

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	auth "real-time-forum/pkg/authentication"
	"real-time-forum/pkg/database"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

// ***************************REGISTER**********************************************************8
// check if pasword meets criteria number length etc, if nickname is not taken
func Register(w http.ResponseWriter, r *http.Request) {
	fmt.Println("getting data")

	var data database.User

	err := json.NewDecoder(r.Body).Decode(&data)
	if err != nil {
		log.Println(err)
	}

	fmt.Println(data)
	fmt.Println(data.LoggedIn)

	data.Password = passwordHash(data.Password)
	//  data.Password = checkPwHash(r.FormValue("password"), data.Password)

	CreateUser(data)
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
		loggedin := "true"
		UpdateUser(user.Nickname, loggedin)
		cookieValue := uuid.NewV4()
		Cookie(w, r, user.Nickname, (cookieValue.String()))
		// sore username of logged in user so you can delete cookie on logout
		auth.Person.Nickname = nickname

	}
	// sends data to js front end
	json.NewEncoder(w).Encode(users)
}

// updates user table
func UpdateUser(nickname, loggedin string) {
	stmt, err := database.DB.Prepare(`UPDATE users SET loggedin = ? WHERE nickname = ?`)
	if err != nil {
		log.Println("Update user:,", err)
	} else {
		stmt.Exec(loggedin, nickname)
	}
}

// logout user NOT WORKING YET
func Logout(w http.ResponseWriter, r *http.Request) {
	// name := "test"
	loggedin := "false"
	fmt.Println("Logging out", auth.Person.Nickname)
	UpdateUser(auth.Person.Nickname, loggedin)

	c1, err := r.Cookie(auth.Person.Nickname)
	fmt.Println("Cookie---", c1)
	if err != nil {
		fmt.Println("On Logout: cookie for user cannot be found!")
	} else {
		fmt.Println("Logout: Cookie deleted, user successfully logged out")
		c1.Name = "Deleted"
		c1.Value = ""
		c1.MaxAge = -1
		auth.Person.Nickname = ""
		fmt.Println(c1.Name)
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
