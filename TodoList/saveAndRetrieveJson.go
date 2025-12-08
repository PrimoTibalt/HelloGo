package main

import (
	"encoding/json"
	"io"
	"os"
	"path"
	"time"
)

type JsonNotes struct {
	Notes []FormattedNote `json:"notes"`
}

type FormattedNote struct {
	Order            int    `json:"order"`
	Text             string `json:"text"`
	RegistrationDate string `json:"registrationDate"`
	DueDate          string `json:"dueDate,omitempty"`
}

func registerNoteJson(note Note) {
	homeDir, homeDirGetError := os.UserHomeDir()
	if homeDirGetError != nil {
		panic(homeDirGetError)
	}
	fullPathToProgramConfigurationJSONFile := path.Join(
		homeDir,
		pathToProgramConfigurationJsonFile)
	_, error := os.Stat(fullPathToProgramConfigurationJSONFile)

	if os.IsNotExist(error) {
		createFileAndDirectoryForApp(homeDir)
	}

	file, openFileError := os.OpenFile(
		fullPathToProgramConfigurationJSONFile, os.O_RDWR, 0644)
	if openFileError != nil {
		panic(openFileError)
	}
	defer file.Close()

	jsonData, readingError := io.ReadAll(file)
	if readingError != nil {
		panic(readingError)
	}

	var jsonModel JsonNotes
	jsonParsingError := json.Unmarshal(jsonData, &jsonModel)
	if jsonParsingError != nil {
		panic(jsonParsingError)
	}

	topNumber := 0
	for _, note := range jsonModel.Notes {
		if topNumber < note.Order {
			topNumber = note.Order
		}
	}

	formattedNote := FormattedNote{
		note.Order,
		note.Text,
		note.RegistrationDate.Format(TimeFormat),
		note.DueDate.Format(TimeFormat),
	}
	jsonModel.Notes = append(jsonModel.Notes, formattedNote)
	jsonResult, marshallingError := json.Marshal(&jsonModel)
	if marshallingError != nil {
		panic(marshallingError)
	}
	_, writeError := file.WriteAt(jsonResult, 0)
	if writeError != nil {
		panic(writeError)
	}
}

func getAllEntriesJson() []Note {
	homeDir, homeDirGetError := os.UserHomeDir()
	if homeDirGetError != nil {
		panic(homeDirGetError)
	}
	fullPathToProgramConfigurationJSONFile := path.Join(
		homeDir,
		pathToProgramConfigurationJsonFile)
	_, error := os.Stat(fullPathToProgramConfigurationJSONFile)

	if os.IsNotExist(error) {
		createFileAndDirectoryForApp(homeDir)
		return []Note{}
	}

	file, openFileError := os.OpenFile(fullPathToProgramConfigurationJSONFile, os.O_RDWR, 0644)
	if openFileError != nil {
		panic(openFileError)
	}

	var jsonModel JsonNotes
	readBytes, readingError := io.ReadAll(file)
	if readingError != nil {
		panic(readingError)
	}
	parsingError := json.Unmarshal(readBytes, &jsonModel)
	if parsingError != nil {
		panic(parsingError)
	}

	notesResult := []Note{}
	for _, fNote := range jsonModel.Notes {
		regDate, regDateParsingError := time.Parse(TimeFormat, fNote.RegistrationDate)

		dueDate, dueDateParsingError := time.Parse(TimeFormat, fNote.DueDate)
		if dueDateParsingError != nil || regDateParsingError != nil {
			continue
		}

		notesResult = append(notesResult, Note{
			fNote.Order,
			fNote.Text,
			regDate,
			dueDate,
		})
	}

	return notesResult
}
