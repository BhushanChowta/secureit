package main

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
)

type SecuredUser struct {
	UserName          string
	EncryptedPassword string
}

var db *sql.DB

const key = "s3cr3t" // Simple key for XOR encryption

func init() {
	// Load .env file
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Get DB details from env
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	dbname := os.Getenv("DB_NAME")

	if user == "" || pass == "" || host == "" || port == "" || dbname == "" {
		log.Fatal("One or more required DB environment variables are missing")
	}

	// Build DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", user, pass, host, port, dbname)

	db, err = sql.Open("mysql", dsn)
	if err != nil {
		log.Fatal("Database connection failed:", err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("Database not reachable:", err)
	}

	createTableQuery := `
	CREATE TABLE IF NOT EXISTS secured_users (
		id INT AUTO_INCREMENT PRIMARY KEY,
		user_name VARCHAR(255) UNIQUE NOT NULL,
		encrypted_password VARCHAR(255) NOT NULL
	);`

	if _, err := db.Exec(createTableQuery); err != nil {
		log.Fatal("Error creating table:", err)
	}
}

func main() {
	fmt.Println("Go ---SECURE IT APP--- Running")

	userName := "chowta"
	userPassword := "12345678"

	err := savePassword(userName, userPassword)
	if err != nil {
		log.Println("Save error:", err)
		return
	}

	encryptedPassowrd, password, err := getPassword(userName)
	if err != nil {
		log.Println("Get error:", err)
		return
	}

	fmt.Println("User:", userName, "-Encrypted password is", encryptedPassowrd, "& Decrypted password is", password)
}

// XOR-based encryption (for demo)
func encryptDecrypt(input string) string {
	output := make([]rune, len(input))
	for i, c := range input {
		output[i] = c ^ rune(key[i%len(key)])
	}
	return string(output)
}

func savePassword(userName, userPassword string) error {
	encrypted := encryptDecrypt(userPassword)
	query := "INSERT INTO secured_users (user_name, encrypted_password) VALUES (?, ?) ON DUPLICATE KEY UPDATE encrypted_password = ?"
	_, err := db.Exec(query, userName, encrypted, encrypted)
	return err
}

func getPassword(userName string) (string, string, error) {
	var encrypted string
	query := "SELECT encrypted_password FROM secured_users WHERE user_name = ?"
	err := db.QueryRow(query, userName).Scan(&encrypted)
	if err != nil {
		return "", "", err
	}
	decrypted := encryptDecrypt(encrypted)
	return encrypted, decrypted, nil
}
