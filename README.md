# WineParser
## _The Simple Wine Data Paser, Ever_


WineParser is a GoLang RestAPI service which process, parse and show data about Wines.


## Features

- Process CSV files and stores in the system using File upload function
- Pares the CSV files and shows the WineID and WineTitle in JSON form



## Tech

WineParser is developed in Go.


## Functions and Structure

- Function UploadFile is used to upload the CSV file in your folder structure temp-files
- Function ReadFile is used to Parse CSV file and shows the WineID and WineTitle in form of JSON.
- Function ParseCsv is used by ReadFile to parse the CSV.

> We need to manually create temp-files folder if you are running this code locally. All the uploaded files will be stored here. 




