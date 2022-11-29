package authentication

import (
	"database/sql"
	"fmt"

	uuid "github.com/satori/go.uuid"
	"golang.org/x/crypto/bcrypt"
)

var Person UserDetails

type UserDetails struct {
	ID                     string
	Email                  string
	Nickname               string
	Password               string
	Accesslevel            bool
	CookieChecker          bool
	Attempted              bool
	RegistrationAttempted  bool
	FailedRegister         bool
	SuccessfulRegistration bool
	PostAdded              bool
}

//register
func newUser(nickname string, email string, password string, db *sql.DB) {
	hash, err := HashPassword(password)
	if err != nil {
		fmt.Println("Error hasing the password", err)
	}
	u1 := uuid.NewV4()
	_, errNewUser := db.Exec("INSERT INTO users (ID, email, nickname, password) VALUES (?, ?, ?, ?)", u1, email, nickname, hash)
	if errNewUser != nil {
		fmt.Printf("The error is %v", errNewUser.Error())
	}
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// store user details in database
// login
// create a cookie and a sessionID on login
// check login info against database
func Register(nickname string, email string, password string, db *sql.DB) {
	rows, err := db.Query("SELECT email FROM users WHERE email = ?", email)
	if err != nil {
		fmt.Println("Registration Error - selecting email from database")
	}
	count := 0
	for rows.Next() {
		count++
	}
	rows2, err2 := db.Query("SELECT nickname FROM users WHERE nickname = ?", nickname)
	if err2 != nil {
		fmt.Println("Registration Error - selecting email from database")
	}
	count2 := 0
	for rows2.Next() {
		count2++
	}
	if count == 0 && count2 == 0 {
	}
}
