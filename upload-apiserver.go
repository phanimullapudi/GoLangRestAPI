package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"encoding/csv"
	"os"
)


type ErrorMessage struct {
    ErrorStatus string `json:"errorstatus"`
	Response string `json:"errorresponse"`
}

type CsvLine struct {
    WineID string
    WineTitle string
}

func UploadStatus(w http.ResponseWriter, r *http.Request) {
	
	fmt.Println("Uploading the file......")
	fmt.Fprintf(w, "Welcome to my Go File Parser - you are the status page")
}

func homePage(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Welcome to my Go File Parser - you are the home page.")
}

func uploadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Uploading File \n")

	// To parse input

	r.ParseMultipartForm(10 << 20)

	// 2. retrieve the file from posted form-data

	file, _, err := r.FormFile("myFile")


	// Error handling of parser
	if err != nil { 
		uploadErrorMessage := ErrorMessage{
			ErrorStatus: "Error",
			Response: err.Error(),
		}
		prettyJSON, _ := json.MarshalIndent(uploadErrorMessage, "", "    ")
		fmt.Printf("%s\n", string(prettyJSON))
		return
	} else {
		uploadErrorMessage := ErrorMessage{
			ErrorStatus: "Success",
			Response: "Able to upload",
		}
		prettyJSON, _ := json.MarshalIndent(uploadErrorMessage, "", "    ")
		fmt.Printf("%s\n", string(prettyJSON))
	}
	
	defer file.Close()

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


// ParseCsv accepts a file and returns its content as a multi-dimentional type
// with lines and each column. Only parses to string type.

func ParseCsv(filename string) ([][]string, error) {

    // Open CSV file
    f, err := os.Open(filename)
    if err != nil {
        return [][]string{}, err
    }
    defer f.Close()

    // Read File into a Variable
    lines, err := csv.NewReader(f).ReadAll()
    if err != nil {
        return [][]string{}, err
    }

    return lines, nil
}

func ReadFile(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Parsing File \n")
	lines, err := ParseCsv("temp-files/upload-998258515.csv")
    if err != nil {
        panic(err)
    }

    // Loop through lines & turn into object
    for _, line := range lines {
        data := CsvLine{
            WineID: line[0],
            WineTitle: line[11],
		}

	
		var jsonData []byte
		jsonData, err := json.Marshal(data)
		if err != nil {
			log.Println(err)
		}
		//fmt.Println(string(jsonData))
		fmt.Fprintf(w, string(jsonData))
	
    }

}

// This function allows to read single item

func ReadSingleItem(w http.ResponseWriter, r *http.Request){
	
}



func handleRequests(){
	http.HandleFunc("/", homePage)
	http.HandleFunc("/status",uploadFile)
	http.HandleFunc("/getwines",ReadFile)
	http.HandleFunc("/getwines/{id}",ReadSingleItem)


	log.Fatal(http.ListenAndServe(":8081", nil))
}

func main() {
	fmt.Println("Upload Manager Loading....")
	handleRequests()
}