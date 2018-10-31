package internal

import (
	"path/filepath"
	"os"
	"log"
	"os/exec"
	"io/ioutil"
	"strings"
	"bufio"
	"fmt"
	"time"
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
	burstCmdPath = toolPath + "BurstWallet"
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
	if _, err := os.Stat(burstCmdPath + "/brs.log"); !os.IsNotExist(err) {
		err2 := os.Remove(burstCmdPath + "/brs.log")
		if err2 != nil{
			log.Print("Attempted to start wallet. Wallet still running!")
			stat.Name = "walletStillRunning"
			statusCh <- stat
		}
	}
	cmd := exec.Command("javaw", "-cp", "burst.jar;conf", "brs.Burst")
	cmd.Dir = burstCmdPath
	cmd.Env = append(os.Environ())
	stdout, err := cmd.StderrPipe()
	if err != nil{
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	//brsLog, err := os.OpenFile(burstCmdPath + "/brs.log", os.O_RDONLY, 0644)
	err = cmd.Start()
	if err != nil{
		log.Fatal(err)
	}
	//reader := bufio.NewScanner(stdout)
	reader := bufio.NewReader(stdout)
	go monitorWallet(statusCh, commandCh, cmd, reader)
}

func monitorWallet(statusCh chan Status, commandCh chan string, cmd *exec.Cmd, reader *bufio.Reader){
	walletIO := make(chan string)
	go func(walletIO chan string){
		defer cmd.Process.Signal(os.Interrupt)
		err := cmd.Wait()
		log.Print("Wallet stopped.")
		walletIO <- "walletStopped"
		if err != nil {
			log.Fatal(err)
		}
		return
	}(walletIO)
	fullErrString := ""

	for{
		select {
		case msg1:= <-commandCh:
			if msg1 == "stopWallet"{
				timer := time.NewTicker(5 * time.Second)
				stat.Name = "walletStopping"
				statusCh <- stat
				for{
					select {
					case <- timer.C:
						cmd.Process.Kill()  //really bad, don't do that!
						stat.Name = "walletStopped"
						statusCh <- stat
						return
					default:
						cmd.Process.Signal(os.Interrupt) //should not work on windows, but sometimes does?
					}
				}
			}
		case msg2:= <-walletIO:
			if msg2 == "walletStopped"{
				stat.Name = "walletStopped"
				statusCh <- stat
				return
			}
		default:
			out, _, err := reader.ReadLine()
			if err != nil{
				log.Fatal(err)
			}
			scannerText := string(out)
			if scannerText != ""{
				fmt.Printf(scannerText + "\n")
			}
			if strings.Contains(scannerText, "brs.db.sql.Db"){
				stat.Name = "walletStarting"
				statusCh <- stat
			}
			if strings.Contains(scannerText, "started successfully."){
				stat.Name = "walletStarted"
				statusCh <- stat
			}
			if strings.Contains(scannerText,"Shutting down..."){
				stat.Name = "walletStopping"
				statusCh <- stat
			}
			if strings.Contains(scannerText,"brs.statistics.StatisticsManagerImpl - handling"){
				stat.Name = "walletLoadingChain"
				statusCh <- stat
			}
			if strings.Contains(scannerText,"[SEVERE]"){
				fullErrString = ""
				fullErrString = fullErrString + scannerText
			}
			if strings.Contains(fullErrString,"[SEVERE]"){
				fullErrString = fullErrString + scannerText
			}
			if strings.Contains(scannerText,"[INFO]") && strings.Contains(fullErrString,"[SEVERE]"){
				stat.Name = "walletError"
				stat.Message = fullErrString
				statusCh <- stat
			}
		}
	}
}