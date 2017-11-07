package main

import (
	"os"
	"net/http"
	"io"
)

func main() {
	// get the pom from server
	// use os package to run the maven build at the location where maven(headless) resides
	//
	downloadPOM("http://10.100.1.85:8484/pom.xml", "pom.xml")

}

func downloadPOM(URL, location string) error {
	// Create the file in the file system
	out, err := os.Create(location)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get pom.xml
	response, err := http.Get(URL)
	if err != nil {
		return err
	}
	defer response.Body.Close()

	// Writes the contents of response's to the created file
	_, err = io.Copy(out, response.Body)
	if err != nil {
		return err
	}

	return nil
}
