package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
)


func UploadStatus(w http.ResponseWriter, r *http.Request) {
	
	fmt.Println("Uploading the file......")
	fmt.Fprintf(w, "Welcome to my Go File Parser - you are the status page")
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my Go File Parser - you are the home page.")
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Uploading File \n")

	// 1. parse input

	r.ParseMultipartForm(10 << 20)


	// 2. retrieve the file from posted form-data

	file, handler, err := r.FormFile("myFile")
	if err != nil {
		fmt.Println("Error Retrieving file from form-data")
		fmt.Println(err)
		return
	}
	defer file.Close()
	fmt.Printf("Uploaded File: %+v\n", handler.Filename)
	fmt.Printf("File Size: %+v\n", handler.Size)
	fmt.Printf("MIME Header: %+v\n", handler.Header)

	// 3. Write temporary file on our server

	tempFile, err := ioutil.TempFile("temp-files", "upload-*.csv")
	if err != nil {
		fmt.Println(err)
		return
	}

	defer tempFile.Close()

	fileBytes, err := ioutil.ReadAll(file)

	if err != nil {
		fmt.Println(err)
	}
	tempFile.Write(fileBytes)

	// 4. return whether or not this has been successful
	fmt.Fprintf(w, "Sucessfully Uploaded File \n")

}

func handleRequests(){
	http.HandleFunc("/", homePage)
	// http.HandleFunc("/status", UploadStatus)
	http.HandleFunc("/upload",uploadFile)
	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	fmt.Println("Upload Manager Loading....")
	handleRequests()
}