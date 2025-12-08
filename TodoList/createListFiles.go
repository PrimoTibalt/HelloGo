package main

import (
	"os"
	"path"
)

const (
	pathToUserSpecificDirectory         string = "/.local/share/primotibalt"
	pathToProgramConfigurationDirectory string = pathToUserSpecificDirectory + "/todolist"
	pathToProgramConfigurationCsvFile   string = pathToProgramConfigurationDirectory + "/config.csv"
	pathToProgramConfigurationJsonFile  string = pathToProgramConfigurationDirectory + "/config.json"
	// pathToProgramConfigurationSqliteDb  string = pathToProgramConfigurationDirectory + "/config.db"
)

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
	fileInfo.Close()

	fileInfo, configurationFileCreationError = os.Create(
		path.Join(homeDir, pathToProgramConfigurationJsonFile))
	if configurationFileCreationError != nil && !os.IsExist(configurationFileCreationError) {
		panic("Couldn't create a configuration file for you. There is something deeply wrong about your system")
	}

	fileInfo.WriteString("{}")
	fileInfo.Close()
}
