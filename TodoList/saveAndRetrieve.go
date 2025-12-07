package main

import (
	"encoding/csv"
	"errors"
	"os"
	"os/user"
	"path"
	"strconv"
	"time"
)

const TimeFormat string = "02.01.2006"

func (n *Note) convertToCsvArray(index int) []string {
	result := []string{strconv.Itoa(index), n.Text}

	if !n.RegistrationDate.IsZero() {
		result = append(result, n.RegistrationDate.Format(TimeFormat))
	}
	if !n.DueDate.IsZero() {
		return append(result, n.DueDate.Format(TimeFormat))
	}

	return result
}

func NewFromCsv(line []string) Note {
	if len(line) < 1 {
		panic("Empty notes are not allowed")
	}

	result := Note{Text: line[1]}
	if len(line) > 1 {
		parsedTime, parsingError := time.Parse(TimeFormat, line[2])
		if parsingError != nil {
			panic(parsingError)
		}
		result.RegistrationDate = parsedTime
	}

	if len(line) > 2 {
		parsedTime, parsingError := time.Parse(TimeFormat, line[3])
		if parsingError != nil {
			panic(parsingError)
		}
		result.DueDate = parsedTime
	}

	return result
}

const (
	pathToUserSpecificDirectory         string = "/.local/share/primotibalt"
	pathToProgramConfigurationDirectory string = pathToUserSpecificDirectory + "/todolist"
	pathToProgramConfigurationCsvFile   string = pathToProgramConfigurationDirectory + "/config.csv"
	// pathToProgramConfigurationJsonFile  string = pathToProgramConfigurationDirectory + "config.json"
	// pathToProgramConfigurationSqliteDb  string = pathToProgramConfigurationDirectory + "config.db"
)

func createAnEntry(note Note) {
	currentUser, currentUserErrorOnGet := user.Current()
	if currentUserErrorOnGet != nil {
		panic(currentUserErrorOnGet)
	}

	fullPathToProgramConfigurationCsvFile := path.Join(currentUser.HomeDir, pathToProgramConfigurationCsvFile)
	_, error := os.Stat(fullPathToProgramConfigurationCsvFile)

	if os.IsNotExist(error) {
		createFileAndDirectoryForApp(currentUser.HomeDir)
	}

	file, readingCsvError := os.OpenFile(fullPathToProgramConfigurationCsvFile, os.O_RDWR, 0644)
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

	fileWriter := csv.NewWriter(file)
	writingNewRecordError := fileWriter.Write(
		note.convertToCsvArray(topNumber + 1),
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

func getAllEntries() []Note {
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
		result = append(result, NewFromCsv(line))
	}

	return result
}

func createFileAndDirectoryForApp(homeDir string) {
	myUserDirectoryCreationError := os.Mkdir(
		path.Join(homeDir, pathToUserSpecificDirectory), 0o700)
	if myUserDirectoryCreationError != nil && !os.IsExist(myUserDirectoryCreationError) {
		panic("Wasn't able to create my user directory, give me some permissions, boi \n" +
			myUserDirectoryCreationError.Error())
	}

	applicationDirectoryCreationError := os.Mkdir(
		path.Join(homeDir, pathToProgramConfigurationDirectory), 0o700)
	if applicationDirectoryCreationError != nil && !os.IsExist(applicationDirectoryCreationError) {
		panic("Wasn't able to create my application directory, that's strange")
	}

	fileInfo, configurationFileCreationError := os.Create(
		path.Join(homeDir, pathToProgramConfigurationCsvFile))
	if configurationFileCreationError != nil && !os.IsExist(configurationFileCreationError) {
		panic("Couldn't create a configuration file for you. There is something deeply wrong about your system")
	}

	defer fileInfo.Close()
}
