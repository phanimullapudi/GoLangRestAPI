package main

import (
	"encoding/csv"
	"encoding/json"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

type StatusMessage struct {
	Status 			string `json:"status,omitempty"`
	Msg    			string `json:"msg,omitempty"`
	Timestamp	    string `json:"ts,omitempty"`
}

type WineRecord struct {
	Id                  int    `json:"id,omitempty"`
	Country             string `json:"country,omitempty"`
	Description         string `json:"description,omitempty"`
	Designation         string `json:"designation,omitempty"`
	Points              string `json:"points,omitempty"`
	Price               string `json:"price,omitempty"`
	Province            string `json:"province,omitempty"`
	Region1             string `json:"region1,omitempty"`
	Region2             string `json:"region2,omitempty"`
	TasterName          string `json:"taster_name,omitempty"`
	TasterTwitterHandle string `json:"taster_twitter_handle,omitempty"`
	Title               string `json:"title,omitempty"`
	Variety             string `json:"variety,omitempty"`
	Winery              string `json:"winery,omitempty"`
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

func ReadFromCSV(fileName string) (map[int]*WineRecord, error) {
	log.Info("Starting loading csv")
	lines, err := ParseCsv(fileName)
	if err != nil {
		log.Error("Failed To Parse the csv file")
		return nil, err
	}
	items := make(map[int]*WineRecord)

	// Loop through lines & turn into object
	for _, line := range lines[1:] {
		id, err := strconv.Atoi(line[0])
		if err != nil {
			log.Error("Failed to parse")
		}

		data := &WineRecord{
			Id:                  id,
			Country:             line[1],
			Description:         line[2],
			Designation:         line[3],
			Points:              line[4],
			Price:               line[5],
			Province:            line[6],
			Region1:             line[7],
			Region2:             line[8],
			TasterName:          line[9],
			TasterTwitterHandle: line[10],
			Title:               line[11],
			Variety:             line[12],
			Winery:              line[13],
		}
		items[id] = data
		if id > latestRecordId {
			latestRecordId = id
		}
	}
	return items, nil
}

// This function allows to read single item
func ReadSingleItem(writer http.ResponseWriter, reader *http.Request) {
	log.Info("Reading single item called")

	vars := mux.Vars(reader)
	key,_ := strconv.Atoi(vars["id"])
	
	if val, ok := globalList[key]; ok {
		json.NewEncoder(writer).Encode(val)
		log.Info("Record found")
	}else {
		status := StatusMessage{
			Msg: "Value Not Found",
			Timestamp: time.Now().Format(time.RFC3339),
		}
		json.NewEncoder(writer).Encode(status)
		log.Error("Record not found")
	}
}

func ReadAllItems(writer http.ResponseWriter,reader *http.Request) {
	log.Info("Reading all items called")

	var wineRecords []*WineRecord

	for _, line := range globalList{
		wineRecords = append(wineRecords, line)
	}
	json.NewEncoder(writer).Encode(wineRecords)
}

func CreateNewItem(w http.ResponseWriter, r *http.Request) {
	log.Info("Create New Item Method Called")

	reqBody, _ := ioutil.ReadAll(r.Body)
	var wineRecord WineRecord
	err := json.Unmarshal(reqBody, &wineRecord)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		log.Error("Record upload failed")
		return
	}

	latestRecordId++
	wineRecord.Id = latestRecordId
	globalList[wineRecord.Id] = &wineRecord
	json.NewEncoder(w).Encode(wineRecord)
}

var globalList = map[int]*WineRecord{}
var latestRecordId = -1
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
	myRouter.HandleFunc("/wine", CreateNewItem).Methods("PUT")
	myRouter.HandleFunc("/wine/{id}", ReadSingleItem)

	log.Fatal(http.ListenAndServe(":8081", myRouter))
}
