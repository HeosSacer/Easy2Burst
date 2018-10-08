package internal

import (
	"fmt"
	"path/filepath"
	"os"
	"os/exec"
	"strings"
	"strconv"
	"log"
)

var (
	toolPath = filepath.ToSlash(os.Getenv("APPDATA") + "/Easy2Burst/")
	downloadCachePath = toolPath + "downloadCache/"
	relevantFileNames = []string{"AppInfo.xml", "BurstWallet", "MariaDB"}
	downloadUrl = "https://download.cryptoguru.org/burst/qbundle/Easy2Burst/"
	logFile, errLogger = os.OpenFile("startup.log", os.O_WRONLY|os.O_CREATE, 0644)
)

func CheckTools() {
	//set logs
	log.SetOutput(logFile)
	log.SetFlags(log.LstdFlags | log.Llongfile)
	//test case
	log.Print("Started Easy2Go Setup")

	if NeedsJava("1.8.0"){
		fmt.Print("TODO")
	}
	listOfMissingFiles := checkFileExistences()
	fmt.Print(listOfMissingFiles)
	processMissingFiles(listOfMissingFiles)
	if errLogger != nil {
		log.Fatal(errLogger)
	}
	//defer to close when you're done with it, not because you think it's idiomatic!
	defer logFile.Close()
}

func NeedsJava(javaVersion string) bool{
	cmd := exec.Command("java", "-version")
	cmd.Env = append(os.Environ())
	out, err := cmd.CombinedOutput()
	currentJavaVersion := strings.Split(string(out[:])," ")[2][1:9]
	if err != nil {
		log.Fatal(err)
	}
	javaVersionNum, _ := strconv.ParseFloat(javaVersion[0:3], 32)
	currentJavaVersionNum, _ := strconv.ParseFloat(currentJavaVersion[0:3],32)
	if javaVersionNum > currentJavaVersionNum{
		return true
	}else {
		return false
	}
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

func processMissingFiles(missingFiles []string) {
	for _, fileName := range missingFiles {
		if fileName != "AppInfo.xml" {
			DownloadFile(downloadCachePath + fileName + ".zip", downloadUrl + fileName + ".zip")
			unzip(downloadCachePath + fileName + ".zip", toolPath + fileName)
		} else{
			DownloadFile(downloadCachePath + fileName, downloadUrl + fileName)
			CopyFile(downloadCachePath + fileName, toolPath + fileName)
		}
	}
	os.RemoveAll(downloadCachePath)
}


