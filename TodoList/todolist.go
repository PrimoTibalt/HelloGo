package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Note struct {
	Text             string
	RegistrationDate time.Time
	DueDate          time.Time
}

func main() {
	switch strings.ToLower(os.Args[1]) {
	case "create":
		if len(os.Args) < 3 {
			fmt.Println("go create [title for a note should have been there]")
		}

		note := Note{Text: os.Args[2], RegistrationDate: time.Now()}
		if len(os.Args) > 3 {
			dueDate, parsigError := time.Parse(TimeFormat,
				os.Args[3])
			if parsigError != nil {
				panic(parsigError)
			}

			note.DueDate = dueDate
		}

		createAnEntry(note)
	case "getall":
		notes := getAllEntries()
		displayNotes(notes)
	}
}
