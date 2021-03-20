package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
	"time"
)

type StatusMessage struct {
	Status 			string `json:"status"`
	Msg    			string `json:"msg,omitempty"`
	Timestamp	    string `json:"ts"`
}

type WineRecord struct {
	Id                  string
	Country             string
	Description         string
	Designation			string
	points              string
	price               string
	province            string
	region1             string
	region2             string
	tasterName          string
	tasterTwitterHandle string
	title               string
	variety             string
	winery              string
}

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

func ReadAllItems(writer http.ResponseWriter,reader *http.Request) {
	log.Info("Reading all items")

	var wineRecords []*WineRecord
	for _, line := range globalList{
		wineRecords = append(wineRecords, line)
	}
	json.NewEncoder(writer).Encode(wineRecords)
}

func ReadFromCSV(fileName string) (map[string]*WineRecord, error) {
	log.Info("Starting loading csv")
	lines, err := ParseCsv(fileName)
	if err != nil {
		log.Error("Failed To Parse the csv file")
		return nil, err
	}
	items := make(map[string]*WineRecord)

	// Loop through lines & turn into object
	for index, line := range lines {
		if index == 0 {
			continue
		}
		data := &WineRecord{
			Id:                  line[0],
			Country:             line[1],
			Description:         line[2],
			Designation: 		 line[3],
			points:              line[4],
			price:               line[5],
			province:            line[6],
			region1:             line[7],
			region2:             line[8],
			tasterName:          line[9],
			tasterTwitterHandle: line[10],
			title:               line[11],
			variety:             line[12],
			winery:              line[13],
		}
		items[line[0]] = data
	}
	return items, nil
}

// This function allows to read single item
func ReadSingleItem(writer http.ResponseWriter, reader *http.Request) {
	vars := mux.Vars(reader)
	key := vars["id"]
	if val, ok := globalList[key]; ok {
		json.NewEncoder(writer).Encode(val)
	}else {
		json.NewEncoder(writer).Encode("")
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

var globalList = map[string]*WineRecord{}
func main() {
	log.SetFormatter(&log.JSONFormatter{})

	status := StatusMessage{
		Status: "ok",
		Timestamp: time.Now().Format(time.RFC3339),
	}

	var err error
	globalList, err = ReadFromCSV("static-files/file.csv")
	if err != nil {
		log.Error("failed to parse data")
		status.Status = "error"
		status.Msg = err.Error()
	}

	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/status", func(writer http.ResponseWriter, request *http.Request) {
		_ = json.NewEncoder(writer).Encode(status)
	})
	myRouter.HandleFunc("/wines", ReadAllItems)
	myRouter.HandleFunc("/wine", createNewItem).Methods("POST")
	myRouter.HandleFunc("/wine/{id}", ReadSingleItem)

	log.Fatal(http.ListenAndServe(":8081", myRouter))
}
