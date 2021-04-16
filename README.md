# WineParser

## _The Simple Wine Data Paser, Ever_

WineParser is a GoLang RestAPI service which process, parse and show data about Wines.

## Features
- Rest API Driven
- Able to upload very large CSV files and parse them
- Read multiple records at once.
- Search and find Wine record by WineID
- Custom upload Wine Record using POST/PUT method

## Tech
WineParser is developed in Go using Gorrilla Mux and Logrus

![Gorilla Logo](https://cloud-cdn.questionable.services/gorilla-icon-64.png)![Go Process Logo](http://i.imgur.com/hTeVwmJ.png)

## Functions and Structure
- The Code is made of 3 REST API calls which uses respective functions - ReadAllItems, ReadSingleItem and CreateNewItem
- In main function we are reading the files from static-files folder and uploading into a Wine object struct. This processing is done by ReadCSV and ParseCSV functions.
- The main function calls the ReadCSV function with file name and location which internally calls ParseCSV.
- The ParseCSV file open the file using OS package and loads the files into variable using CSV package [variable name - lines]
- The variable - lines is used in ReadCSV function to process and store it as WineRecord struct.
- All other Rest API calls use the WineRecord struct to read the data and update the data.
- Every REST API when called - the details of the call is logged into STDOUT. We are also logging any sucess or error to STDOUT.
- I have used both Gorrila mux and also Logrus packages.
- Gorrila mux is used instead of standard http router because we can easily process any query parameters.

## Deviations from the ask

1. The original CSV file link given has some junk and missing characters in the first and last time. I had manually add Id in the first line and also remove the last line in the original file.
2. I was unable to develop Metric output which process and shows output for every minute.  
3. I was also unable to use same path /wine for both ReadALL and PUT calls. So i had to changes /wines to process ReadAllItems and /wine with PUT for CreateNewItem.
4. Also i was unable to read specific columns from the struct like Id and WineTitle for ReadAllItems function.





