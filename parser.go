package main

import (
	"bufio"
	"container/list"
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"html/template"
	"github.com/gorilla/mux"
	"github.com/tealeg/xlsx"
	//"strings"
)

func ReadCsvFile(filePath string) []map[string]interface{} {
	// Load a csv file.
	f, _ := os.Open(filePath)
	// Create a new reader.
	r := csv.NewReader(bufio.NewReader(f))
	result, _ := r.ReadAll()
	parsedData := make([]map[string]interface{}, 0, 0)
	header_name := result[0]

	for row_counter, row := range result {

		if row_counter != 0 {
			var singleMap = make(map[string]interface{})
			for col_counter, col := range row {
				singleMap[header_name[col_counter]] = col
			}
			if len(singleMap) > 0 {

				parsedData = append(parsedData, singleMap)
			}
		}
	}
	fmt.Println("Length of parsedData:", len(parsedData))
	return parsedData

}

func ReadXlsxFile(filePath string) []map[string]interface{} {
	xlFile, err := xlsx.OpenFile(filePath)
	if err != nil {
		fmt.Println("Error reading the file")
	}

	parsedData := make([]map[string]interface{}, 0, 0)
	header_name := list.New()
	// sheet
	for _, sheet := range xlFile.Sheets {
		// rows
		for row_counter, row := range sheet.Rows {

			// column
			header_iterator := header_name.Front()
			var singleMap = make(map[string]interface{})

			for _, cell := range row.Cells {
				if row_counter == 0 {
					text := cell.String()
					header_name.PushBack(text)
				} else {
					text := cell.String()
					singleMap[header_iterator.Value.(string)] = text
					header_iterator = header_iterator.Next()
				}
			}
			if row_counter != 0 && len(singleMap) > 0 {

				parsedData = append(parsedData, singleMap)
			}

		}
	}
	fmt.Println("Length of parsedData:", len(parsedData))
	return parsedData
}

func ExcelCsvParser(blobPath string, blobExtension string) (parsedData []map[string]interface{}) {
	fmt.Println("---------------> We are in product.go")
	if blobExtension == ".csv" {
		fmt.Println("-------We are parsing an csv file.-------------")
		parsedData := ReadCsvFile(blobPath)
		fmt.Printf("Type:%T\n", parsedData)
		return parsedData

	} else if blobExtension == ".xlsx" {
		fmt.Println("----------------We are parsing an xlsx file.---------------")
		parsedData := ReadXlsxFile(blobPath)
		return parsedData
	}
	return parsedData
}

func uploadData(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		fmt.Println("GET")
        t, _ := template.ParseFiles("./templates/index.html")
        t.Execute(w, nil)

	} else if req.Method == "POST"{
		fmt.Println("POST")
		file, handler, err := req.FormFile("uploadfile")
		if err != nil {
			log.Printf("Error while Posting data")
		}
		defer file.Close()
		f, err := os.OpenFile("./data/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
        if err != nil {
            fmt.Println(err)
            return
        }
        defer f.Close()
        io.Copy(f, file)
		blobPath := "./data/" + handler.Filename
		var extension = filepath.Ext(blobPath)
		parsedData := ExcelCsvParser(blobPath, extension)
		parsedJson, _ := json.Marshal(parsedData)
		fmt.Println(string(parsedJson))
		err = os.Remove(blobPath)
		if(err!=nil){
			fmt.Println(err.Error())
		}else{
			fmt.Println("File has been deleted successfully.")
		}
		t, _ := template.ParseFiles("./templates/index.html")
        t.Execute(w, string(parsedJson))
	} else {
	
        fmt.Println("Unknown HTTP " + req.Method + "  Method")
    
	}
	
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", uploadData)
	//http.FileServer(http.Dir("./templates"))
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./templates/")))
	log.Fatal(http.ListenAndServe(":8000", router))
}
