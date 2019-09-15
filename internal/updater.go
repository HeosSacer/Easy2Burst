package internal

import (
	"archive/zip"
	"encoding/xml"
	"errors"
	"fmt"
	"github.com/cavaliercoder/grab"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

// our struct which contains the complete
// array of all Users in the file
type Apps struct {
	XMLName xml.Name `xml:"Apps"`
	App     []App    `xml:"App"`
}

type App struct {
	XMLName xml.Name `xml:"App"`
	Name    string   `xml:"Name"`
	Url     string   `xml:"Url"`
	Version string   `xml:"Version"`
}

var (
	downloader = grab.NewClient()
)

func CheckForUpdates(statusCh chan Status) {
	stat.Name = "checkingUpdates"
	statusCh <- stat
	err := DownloadFile(downloadCachePath, downloadUrl+"AppInfo.xml", statusCh)
	if err != nil {
		log.Fatal(err)
		stat.Name = "updateFailed"
		stat.Message = fmt.Sprintf("%v", err)
		statusCh <- stat
		return
	}
	newApps, err1 := ReadUpdateInfo(downloadCachePath + "AppInfo.xml")
	oldApps, err2 := ReadUpdateInfo(toolPath + "AppInfo.xml")
	if err1 != nil || err2 != nil {
		stat.Name = "updateFailed"
		stat.Message = fmt.Sprintf("%v %v", err1, err2)
		statusCh <- stat
		return
	}
	appsToUpdate, err := CompareUpdateInfo(oldApps, newApps)
	if err != nil {
		stat.Name = "updateFailed"
		stat.Message = fmt.Sprintf("%v", err)
		statusCh <- stat
		return
	}
	for _, appName := range appsToUpdate {
		DownloadFile(downloadCachePath+appName+".zip", downloadUrl+appName+".zip", statusCh)
	}
	stat.Name = "updateFinished"
	statusCh <- stat
	if NeedsJava("1.8.0") {
		stat.Name = "noJava"
		statusCh <- stat
	}
	return
}

func ReadUpdateInfo(pathToUpdateInfo string) (Apps, error) {
	var apps Apps
	xmlFile, err := os.Open(pathToUpdateInfo)
	if err != nil {
		log.Fatal(err)
		return apps, err
	}
	defer xmlFile.Close()
	byteValue, err := ioutil.ReadAll(xmlFile)
	if err != nil {
		log.Fatal(err)
		return apps, err
	}
	xml.Unmarshal(byteValue, &apps)
	return apps, nil
}

func CompareUpdateInfo(oldApps, newApps Apps) ([]string, error) {
	appsToUpdate := make([]string, 0)
	for index, newApp := range newApps.App {
		if index < len(oldApps.App) {
			if newApps.App[index].Name != oldApps.App[index].Name {
				msg := fmt.Sprintf("AppInfo.xml is not consistent: Matching %s and %s", newApps.App[index].Name, oldApps.App[index].Name)
				err := errors.New(msg)
				log.Fatal(err)
				return appsToUpdate, err
			} else {
				if newApps.App[index].Version != oldApps.App[index].Version {
					appsToUpdate = append(appsToUpdate, newApp.Name)
				}
			}
		} else {
			appsToUpdate = append(appsToUpdate, newApp.Name)
		}
	}
	return appsToUpdate, nil
}

func DownloadFile(filepath string, url string, statusCh chan Status) error {
	// Make Folder
	os.MkdirAll(downloadCachePath, os.ModePerm)
	req, err := grab.NewRequest(filepath, url)
	pathSplit := strings.Split(filepath, "/")
	filename := pathSplit[len(pathSplit)-1]
	log.Printf("Downloading [%s] from "+url+" to "+filepath, filename)
	if err != nil {
		log.Fatalf("Download failed: %v\n", err)
		stat.Name = "downloadFailed"
		stat.Message = fmt.Sprintf("%v", err)
		return err
	}
	resp := downloader.Do(req)
	// start UI Loop
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()
Loop:
	for {
		select {
		case <-t.C:
			stat.Name = "downloadMissing"
			stat.Message = fmt.Sprintf("Downloading %s", filename)
			stat.Progress = fmt.Sprintf("%.2f%%", 100*resp.Progress())
			stat.Size = fmt.Sprintf("%.2f", float64(resp.Size())/1000000)
			statusCh <- stat

		case <-resp.Done:
			// download is complete
			log.Printf("Download of %s finished after %.2v", filename, resp.Duration())
			stat.Name = "downloadFinished"
			statusCh <- stat
			break Loop
		}
	}

	// check for errors
	if err := resp.Err(); err != nil {
		log.Fatalf("Download failed: %v\n", err)
		stat.Name = "downloadFailed"
		stat.Message = fmt.Sprintf("%v", err)
	}
	return nil
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func unzip(src string, dest string, statusCh chan Status) ([]string, error) {

	var filenames []string

	r, err := zip.OpenReader(src)
	if err != nil {
		log.Fatal(err)
		return filenames, err
	}
	defer r.Close()

	for _, f := range r.File {

		rc, err := f.Open()
		if err != nil {
			log.Fatal(err)
			return filenames, err
		}
		defer rc.Close()

		// Store filename/path for returning and using later on
		fpath := filepath.Join(dest, f.Name)

		// Check for ZipSlip. More Info: http://bit.ly/2MsjAWE
		if !strings.HasPrefix(fpath, filepath.Clean(dest)+string(os.PathSeparator)) {
			return filenames, fmt.Errorf("%s: illegal file path", fpath)
		}

		filenames = append(filenames, fpath)

		if f.FileInfo().IsDir() {

			// Make Folder
			os.MkdirAll(fpath, os.ModePerm)

		} else {

			// Make File
			if err = os.MkdirAll(filepath.Dir(fpath), os.ModePerm); err != nil {
				log.Fatal(err)
				return filenames, err
			}

			outFile, err := os.OpenFile(fpath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.Mode())
			if err != nil {
				log.Fatal(err)
				return filenames, err
			}

			_, err = io.Copy(outFile, rc)

			// Close the file without defer to close before next iteration of loop
			outFile.Close()

			if err != nil {
				log.Fatal(err)
				return filenames, err
			}

		}
	}
	return filenames, nil
}

func CopyFile(pathSrc, pathDst string) {
	from, err := os.Open(pathSrc)
	if err != nil {
		log.Fatal(err)
	}
	defer from.Close()

	to, err := os.OpenFile(pathDst, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	if err != nil {
		log.Fatal(err)
	}
}

func UpdateE2B(commandCh chan string) {
	commandCh <- ""
}

func UpdateBinary(pathToBinary string) {
	time.Sleep(2 * time.Second)
	os.Remove(pathToBinary)
	CopyFile(os.Args[0], pathToBinary)
	cmd := exec.Command(pathToBinary)
	cmd.Run()
}

func NeedsJava(javaVersion string) bool {
	cmd := exec.Command("java", "-version")
	cmd.Env = append(os.Environ())
	out, err := cmd.CombinedOutput()
	currentJavaVersion := strings.Split(string(out[:]), " ")[2][1:9]
	if err != nil {
		log.Fatal(err)
	}
	javaVersionNum, _ := strconv.ParseFloat(javaVersion[0:3], 32)
	currentJavaVersionNum, _ := strconv.ParseFloat(currentJavaVersion[0:3], 32)
	if javaVersionNum > currentJavaVersionNum {
		return true
	} else {
		return false
	}
}
