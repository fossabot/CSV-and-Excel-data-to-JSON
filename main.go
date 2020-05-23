package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gorilla/mux"
	"github.com/dineshsonachalam/CSV-and-Excel-data-to-JSON/parser"
)

func UploadData(w http.ResponseWriter, req *http.Request) {
	if req.Method == "GET" {
		fmt.Println("GET")
		t, _ := template.ParseFiles("./templates/index.html")
		t.Execute(w, nil)

	} else if req.Method == "POST" {
		fmt.Println("POST")
		file, handler, err := req.FormFile("uploadfile")
		defer file.Close()
		if err != nil {
			log.Printf("Error while Posting data")
			t, _ := template.ParseFiles("./templates/index.html")
			t.Execute(w, nil)
		} else {
			fmt.Println("error throws in else statement")
			fmt.Println("handler.Filename", handler.Filename)
			fmt.Printf("Type of handler.Filename:%T\n", handler.Filename)
			fmt.Println("Length:", len(handler.Filename))
			f, err := os.OpenFile("./data/"+handler.Filename, os.O_WRONLY|os.O_CREATE, 0666)
			if err != nil {
				fmt.Println("Error:", err)
				t, _ := template.ParseFiles("./templates/index.html")
				t.Execute(w, nil)
			}
			defer f.Close()
			io.Copy(f, file)
			blobPath := "./data/" + handler.Filename
			var extension = filepath.Ext(blobPath)
			parsedData := ExcelCsvParser(blobPath, extension)
			parsedJson, _ := json.Marshal(parsedData)
			fmt.Println(string(parsedJson))
			err = os.Remove(blobPath)
			if err != nil {
				fmt.Println(err.Error())
			} else {
				fmt.Println("File has been deleted successfully.")
			}
			t, _ := template.ParseFiles("./templates/index.html")
			t.Execute(w, string(parsedJson))
		}
	} else {
		log.Printf("Error while Posting data")
		t, _ := template.ParseFiles("./templates/index.html")
		t.Execute(w, nil)

	}
}

func ExcelCsvParser(blobPath string, blobExtension string) (parsedData []map[string]interface{}) {
	fmt.Println("---------------> We are in product.go")
	if blobExtension == ".csv" {
		fmt.Println("-------We are parsing an csv file.-------------")
		parsedData := parser.ReadCsvFile(blobPath)
		fmt.Printf("Type:%T\n", parsedData)
		return parsedData

	} else if blobExtension == ".xlsx" {
		fmt.Println("----------------We are parsing an xlsx file.---------------")
		parsedData := parser.ReadXlsxFile(blobPath)
		return parsedData
	} else if blobExtension == ".xls" {
		fmt.Println("----------------We are parsing an xls file.---------------")
		parsedData := parser.ReadXlsFile(blobPath)
		return parsedData
	}
	return parsedData
}

func init() {
	path := "./data"
	if _, err := os.Stat(path); os.IsNotExist(err) {
		os.Mkdir(path, 0777)
		fmt.Println("Created data directory")
	} else {
		fmt.Println("Data directory already exists")
	}
}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", UploadData)
	router.PathPrefix("/").Handler(http.FileServer(http.Dir("./templates/")))
	log.Fatal(http.ListenAndServe(":8000", router))
}
