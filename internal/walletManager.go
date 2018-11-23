package internal

import (
	"bufio"
	"context"
	"fmt"
	"io"
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
	ctx := context.WithValue(context.Background(), "language", "javaw")
	cmd := exec.CommandContext(ctx,"javaw", "-Ddev=true", "-cp", "burst.jar;conf", "brs.Burst")
	cmd.Dir = burstCmdPath
	cmd.Env = append(os.Environ())
	cmd.Env = append(cmd.Env)
	stdout, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}
	stdin, _ := cmd.StdinPipe()
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(1 * time.Second)
	err = cmd.Start()
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(stdout)
	go monitorWallet(statusCh, commandCh, cmd, reader, stdin)
}

func monitorWallet(statusCh chan Status, commandCh chan string, cmd *exec.Cmd, reader *bufio.Reader, stdin io.WriteCloser) {
	go func() {
		for {
			msg1 := <-commandCh
			fmt.Printf("Wallet received %s", msg1)
			if msg1 == "stopWallet" {
				stdin.Write([]byte("shutdown\r\n"))
				stdin.Close()
				return
			}
		}
	}()
	fullErrString := ""

	for {
		out, _, err := reader.ReadLine()
		if err != nil{
			if strings.Contains(err.Error(), "read |0: file already closed") || err == io.EOF{
				stat.Name = "walletStopped"
				statusCh <- stat
				return
			}
			log.Print(err)
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
		if strings.Contains(scannerText, "received command: >shutdown<") {
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
