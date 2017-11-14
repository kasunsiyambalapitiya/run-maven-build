package main

import (
	"os"
	"net/http"
	"io"
	"os/exec"
	"fmt"
	"archive/zip"
	"path/filepath"
	"github.com/mholt/archiver"
	"path"
	"time"
	"strconv"
)

func main() {
	// Check for existence of maven distribution and if it doesn't download and extract maven
	if _,err:=os.Stat("./apache-maven-3.5.2");os.IsNotExist(err){
		// Download maven distribution
		err = downloadContent("http://10.100.1.85:8484/apache-maven-3.5.2-bin.zip", "apache-maven-3.5.2-bin.zip")
		if err != nil {
			fmt.Print(err)
		}
		// Extract maven distribution
		err = extractDistribution("./apache-maven-3.5.2-bin.zip", "./")
		if err != nil {
			fmt.Print(err)
		}
	}

	// Download the pom from server
	err := downloadContent("http://10.100.1.85:8484/pom.xml", "pom.xml")
	if err != nil {
		fmt.Print(err)
	}
	// Extract the given distribution
	err = extractDistribution("./c5-custom-product-5.4.0.zip", "./")
	if err != nil {
		fmt.Print(err)
	}
	// Update the distribution
	err = updateDistribution()
	if err != nil {
		fmt.Print(err)
	}

	timestamp := strconv.FormatInt(time.Now().Unix(), 10)
	// Create the updated zip
	createArchive("./c5-custom-product-5.4.0", timestamp)

	// Delete temp directories/files
	deleteTempFiles()
}

// Extracts the given distribution
func extractDistribution(distributionLocation, destination string) error {
	zipReader, err := zip.OpenReader(distributionLocation)
	if err != nil {
		return err
	}
	defer zipReader.Close()

	// Extract each file to the given destination
	for _, file := range zipReader.File {
		extractFile(file, destination)
	}

	return nil
}

// Extracts each file to the given distribution to the given location
func extractFile(file *zip.File, destination string) error {
	fileContent, error := file.Open()
	if error != nil {
		return error
	}
	defer fileContent.Close()

	// Create the file path
	path := filepath.Join(destination, file.Name)

	if file.FileInfo().IsDir() {
		err := os.MkdirAll(path, 0777)
		if err != nil {
			return err
		}
	} else {
		// Create parent directories if any
		err := os.MkdirAll(filepath.Dir(path), 0777)
		if err != nil {
			return err
		}
		// Open the file for writing
		openedFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE, 0777)
		if err != nil {
			return err
		}

		// Write the content to the opened file
		_, err = io.Copy(openedFile, fileContent)
		if err != nil {
			return err
		}
	}
	return nil
}

// Downloads content from server
func downloadContent(URL, location string) error {
	// Create the file in the file system
	out, err := os.Create(location)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get content from server
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Writes the contents of response to the created file
	_, err = io.Copy(out, response.Body)
	if err != nil {
		return err
	}

	return nil
}

// Use os package to run the maven build using location where maven(headless) resides
func updateDistribution() error {
	command := exec.Command("apache-maven-3.5.2/bin/mvn", "clean", "install")
	command.Dir = "."
	output, err := command.Output()
	if err != nil {
		return err
	}
	fmt.Printf("%s", output)
	return nil
}

// Creates the updated archive
func createArchive(destination, timestamp string) {
	archiver.Zip.Make(path.Join(destination+"-"+timestamp+".zip"), []string{"c5-custom-product-5.4.0"})
}

// Deletes temp directories/files
func deleteTempFiles() {
	os.RemoveAll("./c5-custom-product-5.4.0")
	os.RemoveAll("./target")
	//os.RemoveAll("./apache-maven-3.5.2")
	os.Remove("./pom.xml")
	os.Remove("apache-maven-3.5.2-bin.zip")
}
