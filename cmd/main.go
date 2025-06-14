package main

import "fmt"

type user_name string
type user_password string

func main() {
	fmt.Println("Go ---SECURE IT APP--- Running")

	user_name := "chowta"
	user_password := "12345678"

	savePassword(user_name, user_password) //encrpyt

	password := getPassword(user_name) //decrypt

	fmt.Println("User:", user_name, " is ", password)
}

func savePassword(user_name, user_password) {
	fmt.Print("pass1")
}

func getPassword(user_name) {
	return "1"
}
