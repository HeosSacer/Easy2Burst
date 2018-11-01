package internal

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

type walletCom struct {
	//Commands
	start string
	stop  string
	//Messages
	starting     string
	started      string
	stopping     string
	loadingChain string
}

var (
	WalletCommand = walletCom{
		start: "startWallet",
		stop:  "stopWallet"}
	WalletMessage = walletCom{
		starting:     "walletStarting",
		stopping:     "walletStopping",
		loadingChain: "walletLoadingChain"}
)

func StartWallet(statusCh chan Status, commandCh chan string) {
	if _, err := os.Stat(burstCmdPath + "/brs.log"); !os.IsNotExist(err) {
		err2 := os.Remove(burstCmdPath + "/brs.log")
		if err2 != nil {
			log.Print("Attempted to start wallet. Wallet still running!")
			stat.Name = "walletStillRunning"
			statusCh <- stat
		}
	}
	cmd := exec.Command("javaw", "-cp", "burst.jar;conf", "brs.Burst")
	cmd.Dir = burstCmdPath
	cmd.Env = append(os.Environ())
	stdout, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	//brsLog, err := os.OpenFile(burstCmdPath + "/brs.log", os.O_RDONLY, 0644)
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	//reader := bufio.NewScanner(stdout)
	reader := bufio.NewReader(stdout)
	go monitorWallet(statusCh, commandCh, cmd, reader)
}

func monitorWallet(statusCh chan Status, commandCh chan string, cmd *exec.Cmd, reader *bufio.Reader) {
	walletIO := make(chan string)
	go func(walletIO chan string) {
		defer cmd.Process.Signal(os.Interrupt)
		err := cmd.Wait()
		log.Print("Wallet stopped.")
		walletIO <- "walletStopped"
		if err != nil {
			log.Fatal(err)
		}
		return
	}(walletIO)
	go func() {
		for {
			msg1 := <-commandCh
			fmt.Printf("Wallet received %s", msg1)
			if msg1 == "stopWallet" {
				timer := time.NewTicker(5 * time.Second)
				stat.Name = "walletStopping"
				statusCh <- stat
				for {
					select {
					case <-timer.C:
						cmd.Process.Kill() //really bad, don't do that!
						stat.Name = "walletStopped"
						statusCh <- stat
						return
					default:
						cmd.Process.Signal(os.Interrupt) //should not work on windows, but sometimes does?
					}
				}
			}
		}
	}()
	fullErrString := ""

	for {
		select {
		case msg1 := <-walletIO:
			if msg1 == "walletStopped" {
				stat.Name = "walletStopped"
				statusCh <- stat
				return
			}
		default:
			out, _, err := reader.ReadLine()
			if err != nil {
				log.Fatal(err)
			}
			scannerText := string(out)
			if scannerText != "" {
				fmt.Printf(scannerText + "\n")
			}
			if strings.Contains(scannerText, "brs.db.sql.Db") {
				stat.Name = "walletStarting"
				statusCh <- stat
			}
			if strings.Contains(scannerText, "started successfully.") {
				stat.Name = "walletStarted"
				statusCh <- stat
			}
			if strings.Contains(scannerText, "Shutting down...") {
				stat.Name = "walletStopping"
				statusCh <- stat
			}
			if strings.Contains(scannerText, "brs.statistics.StatisticsManagerImpl - handling") {
				stat.Name = "walletLoadingChain"
				statusCh <- stat
			}
			if strings.Contains(scannerText, "[SEVERE]") {
				fullErrString = ""
				fullErrString = fullErrString + scannerText
			}
			if strings.Contains(fullErrString, "[SEVERE]") {
				fullErrString = fullErrString + scannerText
			}
			if strings.Contains(scannerText, "[INFO]") && strings.Contains(fullErrString, "[SEVERE]") {
				stat.Name = "walletError"
				stat.Message = fullErrString
				statusCh <- stat
			}
		}
	}
}
