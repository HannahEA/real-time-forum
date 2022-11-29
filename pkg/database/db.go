package database
import (
	"database/sql"
	"fmt"
	"log"
	"os"
	_ "github.com/mattn/go-sqlite3"
)
var DB *sql.DB
func InitialiseDB(path string, insertPlaceholders bool) {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		file, err := os.Create(path)
		if err != nil {
			log.Fatal(err.Error())
		}
		fmt.Println(insertPlaceholders)
		if insertPlaceholders {
			defer insertPlaceholdersInDB()
		}
		file.Close()
	}
	// Open the database SQLite file and create the database table
	sqliteDatabase, err1 := sql.Open("sqlite3", path)
	if err1 != nil {
		log.Fatal(err1.Error())
	}
	// Defer the close
	// defer sqliteDatabase.Close()
	// Create the database for each user
	_, errTbl := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "users" (
			"ID"	TEXT UNIQUE,
			"nickname" TEXT UNIQUE,
			"age" TEXT,
			"gender" TEXT,
			"firstname" TEXT,
			"lastname" TEXT,
			"email" 	TEXT UNIQUE,
			"password"	TEXT UNIQUE,
			"loggedin" 	TEXT 
		);
	`)
	if errTbl != nil {
		fmt.Println("USER TABLE ERROR")
		log.Fatal(errTbl.Error())
	}
	// Create the posts table
	_, errPosts := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "posts" (
			"postID"	TEXT,
			"nickname"		TEXT,
			"title" TEXT,
			"categories"	TEXT,
			"body" TEXT
		);
	`)
	if errPosts != nil {
		fmt.Println("POST ERROR")
		log.Fatal(errPosts.Error())
	}
	// Create the cookies table
	_, errCookie := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "cookies" (
			"name"	TEXT,
			"sessionID" 	TEXT UNIQUE
		);
	`)
	if errCookie != nil {
		fmt.Println("TABLE ERROR")
		log.Fatal(errCookie.Error())
	}
	// Create the table for each user
	_, errComments := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "comments" (
			"commentID" TEXT,
			"postID"	TEXT,
			"nickname"	TEXT,
			"body"	TEXT
		);
	`)
	if errComments != nil {
		fmt.Println("USER ERROR")
		log.Fatal(errComments.Error())
	}
	// Create the database for each user
	_, errCategories := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "categories" (
			"postID"	TEXT,
			"Javascript"	INTEGER,
			"Go"	INTEGER,
			"Rust"	INTEGER		);
	`)
	if errCategories != nil {
		fmt.Println("CATEGORY TABLE ERROR")
		log.Fatal(errCategories.Error())
	}
	// Create the database for each user
	_, errConversations := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "conversations" (
			"convoID" TEXT,
			"participants"	TEXT
			);
	`)
	if errConversations != nil {
		fmt.Println("CONVERSATIONS TABLE ERROR")
		log.Fatal(errConversations.Error())
	}
	_, errChats := sqliteDatabase.Exec(`
		CREATE TABLE IF NOT EXISTS "chats" (
			"convoID" TEXT,
			"chatID"	TEXT,
			"sender"	TEXT,
			"date"	TEXT,
			"body" TEXT
			);
	`)
	if errChats != nil {
		fmt.Println("CHATS TABLE ERROR")
		log.Fatal(errChats.Error())
	}
	DB = sqliteDatabase
}
func insertPlaceholdersInDB() {
	queries := map[string]string{
		"fake user 1": fmt.Sprintf(`INSERT INTO users values ("975496ca-9bfc-4d71-8736-da4b6383a575", "bar", "22", "female", "foodie", "barry", "foo@bar.com", "s0fj489fhjsof")`),
		"fake user 2": fmt.Sprintf(`INSERT INTO users values ("6d01e668-2642-4e55-af73-46f057b731f9", "foo", "30", "male", "barry", "fool", "bar@foo.com", "03444f89fsof")`),
		"fake post 1": fmt.Sprintf(`INSERT INTO posts values ("9b4bc963-ecb2-4767-a79b-b09cd102ce4a", "bar", "golang", "Best Coding Language ever", "Golang is really the best!")`),
		"fake post 2": fmt.Sprintf(`INSERT INTO posts values ("16f94e48-82bc-4884-96b3-c847d37f069c", "foo", "javascript", "I love Javascript!", "JS is really neat!")`),
		"fake comment 1": fmt.Sprintf(`INSERT INTO comments values ("49f89e2f-4d7d-4b03-beb6-8def55652d4a", "9b4bc963-ecb2-4767-a79b-b09cd102ce4a", "Cassidy", "I like it too!")`),
		"fake comment 2": fmt.Sprintf(`INSERT INTO comments values ("fbbd419a-e40f-49d5-867a-afa328127cbb", "16f94e48-82bc-4884-96b3-c847d37f069c", "Jeff", "Thanks for this post!")`),
		//bar
		"fake conversation 1 participant 1": fmt.Sprintf(`INSERT INTO conversations values ("0675de06-2d2c-444f-9d0a-ffd3303068d8", "975496ca-9bfc-4d71-8736-da4b6383a575")`),
		//foo
		"fake conversation 1 participant 2": fmt.Sprintf(`INSERT INTO conversations values ("0675de06-2d2c-444f-9d0a-ffd3303068d8", "6d01e668-2642-4e55-af73-46f057b731f9")`),
		//bar
		"fake conversation 2  participant 1": fmt.Sprintf(`INSERT INTO conversations values ("e1953c48-581c-4349-9de0-fb4a81d3745c", "6d01e668-2642-4e55-af73-46f057b731f9")`),
		//foo
		"fake conversation 2  participant 2": fmt.Sprintf(`INSERT INTO conversations values ("e1953c48-581c-4349-9de0-fb4a81d3745c", "975496ca-9bfc-4d71-8736-da4b6383a575")`),
		"fake chat 1": fmt.Sprintf(`INSERT INTO chats values ("e1953c48-581c-4349-9de0-fb4a81d3745c", "d5327a90-e76a-46ef-8b09-531875a534c8", "975496ca-9bfc-4d71-8736-da4b6383a575", "DATE", "Hey! How are you?")`),
		"fake chat 2": fmt.Sprintf(`INSERT INTO chats values ("0675de06-2d2c-444f-9d0a-ffd3303068d8", "d5327a90-e76a-46ef-8b09-531875a534c8", "6d01e668-2642-4e55-af73-46f057b731f9", "DATE", "Good! And you?")`),
	}
	for purpose, q := range queries {
		if _, err := DB.Exec(q); err != nil {
			panic(fmt.Errorf("placeholder query '%s' failed with err: %w", purpose, err))
		} else {
			fmt.Println("successfully executed " + purpose)
		}
	}
}

