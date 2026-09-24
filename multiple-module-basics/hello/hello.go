package main

import (
	"fmt"
	"log"

	"example.com/greetings"
)

func main() {
	log.SetPrefix("Greetings Module Error: ")
	log.SetFlags(0)

	// message, err := greetings.Hello("Nica")

	names := []string{
		"Steven",
		"Nica",
		"Kemp",
	}
	messages, err := greetings.Hellos(names)

	if err != nil {
		log.Fatal(err)
	}

	// fmt.Println(message)
	fmt.Println(messages["Kemp"])
}