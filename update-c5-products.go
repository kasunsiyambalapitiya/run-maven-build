package main

import (
	"os"
	"net/http"
	"io"
	"os/exec"
	"fmt"
	"archive/zip"
	"path/filepath"
)

func main() {
	// Download the pom from server
	downloadPOM("http://10.100.1.85:8484/pom.xml", "pom.xml")
	// Extract the given distribution
	extractDistribution("./c5-custom-product-5.3.0.zip", "./")
	// Update the distribution
	err := updateDistribution()
	if err != nil {
		fmt.Print(err)
	}

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

	// Todo change the m2 location
	return nil
}

// Downloads the pom from server
func downloadPOM(URL, location string) error {
	// Create the file in the file system
	out, err := os.Create(location)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get pom.xml from server
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
