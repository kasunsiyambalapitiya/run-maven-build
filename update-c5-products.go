package main

import (
	"os"
	"net/http"
	"io"
	"os/exec"
	"fmt"
)

func main() {
	downloadPOM("http://10.100.1.85:8484/pom.xml", "pom.xml")
	err := updateDistribution()
	if err != nil {
		fmt.Print(err)
	}
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
