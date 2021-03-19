package main

import (
	"fmt"
	"log"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"encoding/csv"
	"os"
	"github.com/gorilla/mux"
)


type ErrorMessage struct {
    ErrorStatus string `json:"errorstatus"`
	Response string `json:"errorresponse"`
}

type CsvLine struct {
    WineID string
    WineTitle string
}

type WineRecord struct {
	Id int
	Country string
	Description string
	points string
	price string
	province string
	region_1 string
	region_2 string
	taster_name string
	taster_twitter_handle string
	title string
	variety string
	winery string
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

func ReadAllItems(w http.ResponseWriter, r *http.Request) {
	fmt.Println("Pasring all Records Method Called")
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
	vars := mux.Vars(r)
    key := vars["id"]

    fmt.Fprintf(w, "Key: " + key)

	fmt.Println("Single Item Fetch Method Called")
	lines, err := ParseCsv("temp-files/upload-998258515.csv")
    if err != nil {
        panic(err)
    }

	for _, line := range lines {
		if line[0] == key {
			data := CsvLine {
				WineID: line[0],
				WineTitle: line[11],
			}
			jsonData, _ := json.Marshal(data)
			fmt.Fprintf(w, string(jsonData))
		}

	}
}


func createNewItem(w http.ResponseWriter, r *http.Request) {

    fmt.Println("Create New Item Method Called")
	
	var winerecord WineRecord
	err := json.NewDecoder(r.Body).Decode(&winerecord)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
	fmt.Fprintf(w, "Winerecord: %+v", winerecord)

}

func handleRequests(){

	myRouter := mux.NewRouter().StrictSlash(true)

	myRouter.HandleFunc("/", homePage)
	myRouter.HandleFunc("/status",uploadFile)
	myRouter.HandleFunc("/wines",ReadAllItems)
	myRouter.HandleFunc("/wine", createNewItem).Methods("POST")
	myRouter.HandleFunc("/wine/{id}",ReadSingleItem)
	
	log.Fatal(http.ListenAndServe(":8081", myRouter))
}

func main() {
	fmt.Println("Upload Manager Loading....")
	handleRequests()
}