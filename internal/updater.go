package internal

import (
	"encoding/xml"
	"os"
	"fmt"
	"io/ioutil"
	"path/filepath"
	"net/http"
	"io"
	"strings"
	"archive/zip"
	"log"
)


// our struct which contains the complete
// array of all Users in the file
type Apps struct {
	XMLName xml.Name `xml:"Apps"`
	App    []App   `xml:"App"`
}

type App struct{
	XMLName xml.Name `xml:"App"`
	Name	string	`xml:"Name"`
	Url		string	`xml:"Url"`
	Version	string	`xml:"Version"`
}

func CheckForUpdates(){
	err := DownloadFile(downloadCachePath, downloadUrl + "AppInfo.xml")
	if err != nil {
		log.Fatal(err)
	}
}

func ReadUpdateInfo(pathToUpdateInfo string) Apps{
	var apps Apps
	xmlFile, err := os.Open(pathToUpdateInfo)
	if err != nil {
		log.Fatal(err)
	}
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)
	xml.Unmarshal(byteValue, &apps)
	return apps
}

func DownloadFile(filepath string, url string) error {

	// Create the filepath
	err := os.MkdirAll(downloadCachePath, 0644)
	if err != nil {
		log.Fatal(err)
		return err
	}
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
		return err
	}
	defer resp.Body.Close()

	// Write the body to file
	out, err := os.Create(filepath)
	if err != nil {
		log.Fatal(err)
		return err
	}
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		log.Fatal(err)
		return err
	}

	return nil
}

// Unzip will decompress a zip archive, moving all files and folders
// within the zip file (parameter 1) to an output directory (parameter 2).
func unzip(src string, dest string) ([]string, error) {

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

func CopyFile(pathSrc, pathDst string){
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
