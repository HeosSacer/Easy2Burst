package internal

import (
	"path/filepath"
	"os"
	"log"
	"os/exec"
	"io/ioutil"
	"strings"
	"bufio"
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
	burstCmdPath = toolPath + "BurstWallet/"
	relevantFileNames = []string{"AppInfo.xml", "BurstWallet", "MariaDB"}
	downloadUrl = "https://download.cryptoguru.org/burst/qbundle/Easy2Burst/"
	stat = Status{
		Name: "starting",
		Message: "Starting Setup...",
		Progress: "0%",
		Size: "",
	}
)

func CheckTools(statusCh chan Status, commandCh chan string) {
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
	stat.Name = "checkBurstDB"
	statusCh <- stat
	CheckBurstDB()
	stat.Name = "setupFinished"
	statusCh <- stat
	CheckForUpdates(statusCh)
	stat.Name = "updaterFinished"
	statusCh <- stat
	StartWallet(statusCh, commandCh)
	stat.Name = "walletStarted"
	statusCh <- stat
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


func CheckBurstDB(){
	//Write the path of the .sql into the BAT file
	file, err := ioutil.ReadFile(toolPath + "MariaDB/bin/setupDb.BAT")
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	newContents := strings.Replace(string(file), "{PATH_TO_SQL_SCRIPT}", toolPath + "BurstWallet/init-mysql.sql", -1)
	err = ioutil.WriteFile(toolPath + "MariaDB/bin/setupDb.BAT", []byte(newContents), 0)
	if err != nil {
		log.Fatal(err)
		panic(err)
	}
	cmd := exec.Command(toolPath + "MariaDB/bin/setupDb.BAT")
	cmd.Env = append(os.Environ())
	out, err := cmd.CombinedOutput()
	log.Print(string(out))
	if err != nil{
		log.Fatal(err)
	}
}

func StartWallet(statusCh chan Status, commandCh chan string){
	cmd := exec.Command(burstCmdPath + "burst.cmd")
	stdout, err := cmd.StdoutPipe()
	cmd.Env = append(os.Environ())
	scanner := bufio.NewScanner(stdout)
	err = cmd.Start()
	if err != nil{
		log.Fatal(err)
	}
	go monitorWallet(statusCh, commandCh, cmd, scanner)
}

func monitorWallet(statusCh chan Status, commandCh chan string, cmd *exec.Cmd, scanner *bufio.Scanner){
	walletIO := make(chan string)
	go func(walletIO chan string){
		err := cmd.Wait()
		log.Print("Wallet stopped.")
		walletIO <- "walletStopped"
		if err != nil{
			log.Fatal(err)
		}
	}(walletIO)
	fullErrString := ""
	for{
		select {
		case msg:= <-commandCh:
			if msg == "stopWallet"{
				cmd.Process.Kill()
			}
		case msg:= <-walletIO:
			if msg == "walletStopped"{
				stat.Name = "walletStopped"
				statusCh <- stat
				return
			}
		default:
			scanner.Scan()
			if strings.Contains(scanner.Text(), "loadProperties"){
				stat.Name = "walletStarting"
				statusCh <- stat
			}
			if strings.Contains(scanner.Text(), "started successfully."){
				stat.Name = "walletStarted"
				statusCh <- stat
			}
			if strings.Contains(scanner.Text(),"Shutting down..."){
				stat.Name = "walletStopping"
				statusCh <- stat
			}
			if strings.Contains(scanner.Text(),"brs.statistics.StatisticsManagerImpl - handling"){
				stat.Name = "walletLoadingChain"
				statusCh <- stat
			}
			if strings.Contains(scanner.Text(),"[SEVERE]"){
				fullErrString = ""
				fullErrString = fullErrString + scanner.Text()
			}
			if strings.Contains(fullErrString,"[SEVERE]"){
				fullErrString = fullErrString + scanner.Text()
			}
			if strings.Contains(scanner.Text(),"[INFO]") && strings.Contains(fullErrString,"[SEVERE]"){
				stat.Name = "walletError"
				stat.Message = fullErrString
				statusCh <- stat
			}
		}
	}
}