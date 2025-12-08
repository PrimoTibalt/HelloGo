package main

import (
	"fmt"
	"os"
	"strings"
	"time"
)

type Note struct {
	Order            int
	Text             string
	RegistrationDate time.Time
	DueDate          time.Time
}

type DateTime time.Time

var (
	TextArgumentIndex    = 2
	DueDateArgumentIndex = 3
)

func main() {
	switch strings.ToLower(os.Args[1]) {
	case "create":
		createNote()
	case "getall":
		notes := getAllEntries()
		displayNotes(notes)
	}
}

func getAllEntries() []Note {
	source := "csv"
	if len(os.Args) > 2 {
		source = strings.ToLower(os.Args[2])
	}

	switch source {
	case "csv":
		return getAllEntriesCsv()
	case "json":
		return getAllEntriesJson()
	default:
		panic("unknown source type '" + source + "'")
	}
}

func createNote() {
	if len(os.Args) < 3 {
		fmt.Println("go create [title for a note should have been there]")
	}

	var registerFn func(note Note)
	switch strings.ToLower(os.Args[TextArgumentIndex]) {
	case "json":
		registerFn = registerNoteJson
		TextArgumentIndex++
		DueDateArgumentIndex++
	default:
		registerFn = registerNoteCsv
	}

	note := Note{Order: 0, Text: os.Args[TextArgumentIndex], RegistrationDate: time.Now()}
	if len(os.Args) >= DueDateArgumentIndex-1 {
		dueDate, parsigError := time.Parse(
			TimeFormat,
			os.Args[DueDateArgumentIndex])
		if parsigError != nil {
			panic(parsigError)
		}
		note.DueDate = dueDate
	}

	registerFn(note)
}
