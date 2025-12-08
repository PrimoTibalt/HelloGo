package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"os/user"
	"path"
	"strconv"
	"time"
)

const TimeFormat string = "02.01.2006"

func (n *Note) convertToCsvArray() []string {
	result := []string{strconv.Itoa(n.Order), n.Text}

	if !n.RegistrationDate.IsZero() {
		result = append(result, n.RegistrationDate.Format(TimeFormat))
	}
	if !n.DueDate.IsZero() {
		return append(result, n.DueDate.Format(TimeFormat))
	}

	return result
}

func NewFromCsv(line []string) (Note, error) {
	if len(line) < 2 {
		panic("Empty notes are not allowed")
	}

	order, parsingError := strconv.Atoi(line[0])
	if parsingError != nil {
		return Note{}, errors.New("Error parsing order of a note")
	}

	result := Note{Order: order, Text: line[1]}
	if len(line) > 2 {
		parsedTime, parsingError := time.Parse(TimeFormat, line[2])
		if parsingError != nil {
			panic(parsingError)
		}
		result.RegistrationDate = parsedTime
	}

	if len(line) > 3 {
		parsedTime, parsingError := time.Parse(TimeFormat, line[3])
		if parsingError != nil {
			panic(parsingError)
		}
		result.DueDate = parsedTime
	}

	return result, nil
}

func registerNoteCsv(note Note) {
	currentUser, currentUserErrorOnGet := user.Current()
	if currentUserErrorOnGet != nil {
		panic(currentUserErrorOnGet)
	}

	fullPathToProgramConfigurationCsvFile := path.Join(currentUser.HomeDir, pathToProgramConfigurationCsvFile)
	_, error := os.Stat(fullPathToProgramConfigurationCsvFile)

	if os.IsNotExist(error) {
		createFileAndDirectoryForApp(currentUser.HomeDir)
	}

	file, readingCsvError := os.OpenFile(
		fullPathToProgramConfigurationCsvFile,
		os.O_RDWR,
		0644)
	if readingCsvError != nil {
		panic(readingCsvError)
	}

	fileReader := csv.NewReader(file)
	existingRecords, readingAllRecordsError := fileReader.ReadAll()
	if readingAllRecordsError != nil {
		panic(readingAllRecordsError)
	}

	topNumber := 0
	for _, item := range existingRecords {
		number, conversionError := strconv.ParseInt(item[0], 10, 0)
		if conversionError != nil {
			panic(conversionError)
		}

		if topNumber < int(number) {
			topNumber = int(number)
		}
	}
	note.Order = topNumber + 1

	fileWriter := csv.NewWriter(file)
	writingNewRecordError := fileWriter.Write(
		note.convertToCsvArray(),
	)
	if writingNewRecordError != nil {
		panic(writingNewRecordError)
	}

	fileWriter.Flush()
	fileCloseError := file.Close()
	if fileCloseError != nil {
		panic(fileCloseError)
	}
}

func getAllEntriesCsv() []Note {
	currentUser, currentUserErrorOnGet := user.Current()
	if currentUserErrorOnGet != nil {
		gettingCurrentUserError := errors.New("There is something wrong with getting current user, I'm out.")
		panic(gettingCurrentUserError)
	}

	fullPathToProgramConfigurationCsvFile := path.Join(currentUser.HomeDir, pathToProgramConfigurationCsvFile)
	_, error := os.Stat(fullPathToProgramConfigurationCsvFile)

	if os.IsNotExist(error) {
		createFileAndDirectoryForApp(currentUser.HomeDir)
	}

	file, readingCsvError := os.OpenFile(fullPathToProgramConfigurationCsvFile, os.O_RDONLY, 0644)
	if readingCsvError != nil {
		panic(readingCsvError)
	}

	fileReader := csv.NewReader(file)
	content, readingLinesError := fileReader.ReadAll()
	if readingLinesError != nil {
		panic(readingLinesError)
	}

	fileCloseError := file.Close()
	if fileCloseError != nil {
		panic(fileCloseError)
	}

	result := []Note{}
	for _, line := range content {
		note, parsingError := NewFromCsv(line)
		if parsingError != nil {
			fmt.Println(parsingError.Error())
		} else {
			result = append(result, note)
		}
	}

	return result
}
