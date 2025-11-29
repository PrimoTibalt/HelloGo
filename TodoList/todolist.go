package main

import (
	"fmt"
	"os"
	"os/user"
	"path"
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

	_, error := os.Stat(
		path.Join(currentUser.HomeDir, pathToProgramConfigurationCsvFile))

	if !os.IsExist(error) {
		createFileAndDirectoryForApp(currentUser.HomeDir)
	}
}

func getAllEntries() {
	fmt.Println("nothing")
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
