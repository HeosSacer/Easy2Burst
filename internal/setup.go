package internal

import (
	"path/filepath"
	"os"
	"log"
)

type Status struct {
	Name string
	Message string
	Progress  string
	Size	string
}

var (
	toolPath = filepath.ToSlash(os.Getenv("APPDATA") + "/Easy2Burst/")
	downloadCachePath = toolPath + "downloadCache/"
	relevantFileNames = []string{"AppInfo.xml", "BurstWallet", "MariaDB"}
	downloadUrl = "https://download.cryptoguru.org/burst/qbundle/Easy2Burst/"
	stat = Status{
		Name: "starting",
		Message: "Starting Setup...",
		Progress: "0%",
		Size: "",
	}
)

func CheckTools(statusCh chan Status) {
	//set logs
	os.Remove(toolPath + "startup.log")
	logFile, _ := os.OpenFile(toolPath + "startup.log", os.O_WRONLY|os.O_CREATE, 0644)
	defer logFile.Close()
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Llongfile)
	//test case
	log.Print("Started Easy2Burst Setup")
	statusCh <- stat
	listOfMissingFiles := checkFileExistences()
	if len(listOfMissingFiles) > 0{
		stat.Name = "startSetup"
		statusCh <- stat
		log.Printf("Files %s missing.", listOfMissingFiles)
		processFiles(listOfMissingFiles, statusCh)
	}
	stat.Name = "setupFinished"
	statusCh <- stat
	CheckForUpdates(statusCh)
}

func checkFileExistences() []string {
	missingFiles := make([]string, 0)
	for _, fileName := range relevantFileNames {
		if _, err := os.Stat(toolPath + fileName); os.IsNotExist(err) {
			missingFiles = append(missingFiles, fileName)
		}
	}
	return missingFiles
}

func processFiles(missingFiles []string, statusCh chan Status) {
	for _, fileName := range missingFiles {
		if fileName != "AppInfo.xml" {
			DownloadFile(downloadCachePath + fileName + ".zip", downloadUrl + fileName + ".zip", statusCh)
			unzip(downloadCachePath + fileName + ".zip", toolPath + fileName, statusCh)
		} else{
			DownloadFile(downloadCachePath + fileName, downloadUrl + fileName, statusCh)
			CopyFile(downloadCachePath + fileName, toolPath + fileName)
		}
	}
	os.RemoveAll(downloadCachePath)
}


