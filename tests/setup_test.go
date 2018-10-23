package tests

import (
	"testing"
	"time"
	"reflect"
	"io/ioutil"
	"fmt"
	"strings"
	"github.com/HeosSacer/Easy2Burst/internal"
)

/*
func TestIntegrationCheckTools(t *testing.T){
	statusCh := make(chan internal.Status, 1)
	commandCh := make(chan string)
	go internal.CheckTools(statusCh, commandCh)

	Loop:
		for {
			stat := <-statusCh
			if stat.Name == "downloadMissing"{
				fmt.Printf("CheckToolsTest: %s %s (%s) \n", stat.Message, stat.Size, stat.Progress)
			}
			if stat.Name == "setupFinished" {
				break Loop
			}
			if stat.Name == "downloadFailed" {
				t.Fail()
			}
		}
}

func TestNeedsJava(t *testing.T){
	result := internal.NeedsJava("1.8.0")
	AssertEqual(t, result, false)
	result = internal.NeedsJava("1.9.0")
	AssertEqual(t, result, true)
}

func TestCheckForUpdates(t *testing.T){
	statusCh := make(chan internal.Status, 1)
	pathToXml := filepath.ToSlash(os.Getenv("APPDATA") + "/Easy2Burst/") + "AppInfo.xml"
	cleanUp := func (){
		statusCh = make(chan internal.Status, 1)
		ReplaceLine(pathToXml,"<Version>2.3", "\t\t<Version>2.2.3</Version>")
		os.RemoveAll(filepath.ToSlash(os.Getenv("APPDATA") + "/Easy2Burst/downloadCache"))
	}
	cleanUp()
	go internal.CheckForUpdates(statusCh)
	//No Updates
Loop1:
	for {
		stat := <-statusCh
		if stat.Name == "downloadMissing"{
			t.Errorf("Downloading something, despite everthing is up to date!")
		}
		if stat.Name == "updateFinished" {
			break Loop1
		}
		if stat.Name == "downloadFailed" {
			t.Errorf("Download failed!")
		}
	}
	//Force Updates
	ReplaceLine(pathToXml,"<Version>2.2","\t\t<Version>2.3.2</Version>")
	go internal.CheckForUpdates(statusCh)
Loop2:
	for {
		stat := <-statusCh
		if stat.Name == "downloadMissing"{
			fmt.Printf("CheckForUpdatesTest: %s %s (%s) \n", stat.Message, stat.Size, stat.Progress)
			if stat.Message != "Downloading BurstWallet.zip"{
				cleanUp()
				t.Errorf("Downloading something else than BurstWallet.zip")
			}
		}
		if stat.Name == "updateFinished" {
			cleanUp()
			break Loop2
		}
		if stat.Name == "downloadFailed" {
			cleanUp()
			t.Errorf("Download Failed!")
		}
	}
}

func TestCheckBurstDB(t *testing.T){
	t.Fail()
	//internal.CheckBurstDB()
}
*/
func TestStartWallet(t *testing.T){
	statusCh := make(chan internal.Status, 1)
	commandCh := make(chan string)
	go internal.StartWallet(statusCh, commandCh)
	checkArray := []bool {false, false, false, false}
	//Timeout if it takes too long
	stoppedRegular := false
	timer := time.NewTicker(30 * time.Second)
	defer timer.Stop()
	Loop:
	for{
		select {
		case <-timer.C:
			commandCh <- "stopWallet"
			break Loop
		case stat := <-statusCh:
			if stat.Name == "walletStarting" {
				checkArray[0] = true
			}
			if stat.Name == "walletStarted" {
				checkArray[1] = true
				commandCh <- "stopWallet"
			}
			if stat.Name == "walletStopping" {
				checkArray[2] = true
			}
			if stat.Name == "walletStopped" {
				checkArray[3] = true
				stoppedRegular = true
				break Loop
			}
			if stat.Name == "walletError"{
				fmt.Print(stat.Message)
				t.Fail()
			}
		default:
		}
	}
	if stoppedRegular == false{
		timer = time.NewTicker(5 * time.Second)
	Loop2:
		for {
			select {
			case <-timer.C:
				break Loop2
			case stat := <-statusCh:
				if stat.Name == "walletStopped" {
					checkArray[3] = true
					break Loop2
				}
			}
		}
	}

	checkMask :=[]bool {true, true, true, true}
	AssertEqual(t, checkArray, checkMask)
}

// AssertEqual checks if values are equal
func AssertEqual(t *testing.T, a interface{}, b interface{}) {
	if a == b {
		return
	}
	// debug.PrintStack()
	t.Errorf("Received %v (type %v), expected %v (type %v)", a, reflect.TypeOf(a), b, reflect.TypeOf(b))
}

func ReplaceLine(filepath string, existingLineContains string, replacingString string){
	input, err := ioutil.ReadFile(filepath)
	if err != nil {
		fmt.Print(err)
	}

	lines := strings.Split(string(input), "\n")

	for i, line := range lines {
		if strings.Contains(line, existingLineContains) {
			lines[i] = replacingString
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(filepath, []byte(output), 0644)
	if err != nil {
		fmt.Print(err)
	}
}