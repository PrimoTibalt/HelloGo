package main

import (
	"encoding/csv"
	"fmt"
	"os"
	"os/user"
	"path"
	"strconv"
	"strings"
)

var (
	pathToUserSpecificDirectory         string = "/.local/share/primotibalt"
	pathToProgramConfigurationDirectory string = pathToUserSpecificDirectory + "/todolist"
	pathToProgramConfigurationCsvFile   string = pathToProgramConfigurationDirectory + "/config.csv"
	// pathToProgramConfigurationJsonFile  string = pathToProgramConfigurationDirectory + "config.json"
	// pathToProgramConfigurationSqliteDb  string = pathToProgramConfigurationDirectory + "config.db"
)

func main() {
	switch strings.ToLower(os.Args[1]) {
	case "create":
		if len(os.Args) < 3 {
			fmt.Println("go create [title for a note should have been there]")
			panic(1)
		}

		createAnEntry(os.Args[2])
	case "getall":
		getAllEntries()
	}
}

func createAnEntry(text string) {
	currentUser, currentUserErrorOnGet := user.Current()
	if currentUserErrorOnGet != nil {
		fmt.Println("There is something wrong with getting current user, I'm out.")
		return
	}

	fullPathToProgramConfigurationCsvFile := path.Join(currentUser.HomeDir, pathToProgramConfigurationCsvFile)
	_, error := os.Stat(fullPathToProgramConfigurationCsvFile)

	if os.IsNotExist(error) {
		createFileAndDirectoryForApp(currentUser.HomeDir)
	}

	file, readingCsvError := os.OpenFile(fullPathToProgramConfigurationCsvFile, os.O_RDWR, 0644)
	if readingCsvError != nil {
		fmt.Println("There was an error reading our csv.")
		fmt.Println(readingCsvError.Error())
		closingError := file.Close()
		if closingError != nil {
			fmt.Println("Couldn't close connection to file because...")
			fmt.Println(closingError.Error())
		}
		return
	}

	fileReader := csv.NewReader(file)
	existingRecords, readingAllRecordsError := fileReader.ReadAll()
	if readingAllRecordsError != nil {
		fmt.Println(readingAllRecordsError.Error())
	}

	topNumber := 0
	for _, item := range existingRecords {
		number, conversionError := strconv.ParseInt(item[0], 10, 0)
		if conversionError != nil {
			fmt.Println("Check your " + fullPathToProgramConfigurationCsvFile + " . Data seems to be corrupted.")
		}

		if topNumber < int(number) {
			topNumber = int(number)
		}
	}

	fileWriter := csv.NewWriter(file)
	writingNewRecordError := fileWriter.Write([]string{strconv.Itoa(topNumber + 1), text})
	if writingNewRecordError != nil {
		fmt.Println("Error writing new line to csv.")
		fmt.Println(writingNewRecordError.Error())
		return
	}

	fileWriter.Flush()
	fileCloseError := file.Close()
	if fileCloseError != nil {
		fmt.Println("Couldn't close connection to file because...")
		fmt.Println(fileCloseError.Error())
	}

	fmt.Println("Your new todo item was added.")
}

func getAllEntries() {
	currentUser, currentUserErrorOnGet := user.Current()
	if currentUserErrorOnGet != nil {
		fmt.Println("There is something wrong with getting current user, I'm out.")
		return
	}

	fullPathToProgramConfigurationCsvFile := path.Join(currentUser.HomeDir, pathToProgramConfigurationCsvFile)
	_, error := os.Stat(fullPathToProgramConfigurationCsvFile)

	if os.IsNotExist(error) {
		fmt.Println("Couldn't file your todoitems file. Creating new one.")
		createFileAndDirectoryForApp(currentUser.HomeDir)
	}

	file, readingCsvError := os.OpenFile(fullPathToProgramConfigurationCsvFile, os.O_RDONLY, 0644)
	if readingCsvError != nil {
		fmt.Println("There was an error reading our csv.")
		fmt.Println(readingCsvError.Error())
		closingError := file.Close()
		if closingError != nil {
			fmt.Println("Couldn't close connection to file because...")
			fmt.Println(closingError.Error())
		}
		return
	}

	fileReader := csv.NewReader(file)
	content, readingLinesError := fileReader.ReadAll()
	if readingLinesError != nil {
		fmt.Println("Error reading lines of csv.")
		fmt.Println(readingLinesError.Error())
	}

	if len(content) < 1 {
		fmt.Println("You don't have anything planned.")
	}

	for _, line := range content {
		fmt.Println(line)
	}

	fileCloseError := file.Close()
	if fileCloseError != nil {
		fmt.Println("Couldn't close connection to file because...")
		fmt.Println(fileCloseError.Error())
	}
}

func createFileAndDirectoryForApp(homeDir string) {
	myUserDirectoryCreationError := os.Mkdir(
		path.Join(homeDir, pathToUserSpecificDirectory), 0o700)
	if myUserDirectoryCreationError != nil && !os.IsExist(myUserDirectoryCreationError) {
		fmt.Println("Wasn't able to create my user directory, give me some permissions, boi")
		fmt.Println(myUserDirectoryCreationError.Error())
		return
	}

	applicationDirectoryCreationError := os.Mkdir(
		path.Join(homeDir, pathToProgramConfigurationDirectory), 0o700)
	if applicationDirectoryCreationError != nil && !os.IsExist(applicationDirectoryCreationError) {
		fmt.Println("Wasn't able to create my application directory, that's strange")
		return
	}

	fileInfo, configurationFileCreationError := os.Create(
		path.Join(homeDir, pathToProgramConfigurationCsvFile))
	if configurationFileCreationError != nil && !os.IsExist(configurationFileCreationError) {
		fmt.Println("Couldn't create a configuration file for you. There is something deeply wrong about your system")
		return
	}

	defer fileInfo.Close()
}
