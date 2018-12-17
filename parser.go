package main

import (
	"fmt"
	"bufio"
    "encoding/csv"
    "encoding/json"
    "container/list"
    "github.com/tealeg/xlsx"
    "os"
    "path/filepath"
    "io/ioutil"
    "net/http"
    "io"
    "log"
    "github.com/gorilla/mux"
    //"strings"
)

func ReadCsvFile(filePath string)([]map[string]interface{}) {
    // Load a csv file.
    f, _ := os.Open(filePath)
    // Create a new reader.
    r := csv.NewReader(bufio.NewReader(f))
    result, _ := r.ReadAll()
    parsedData := make([]map[string]interface{}, 0, 0)
    header_name:=result[0]
   

    for row_counter,row := range result {
        
        if(row_counter!=0){
            var singleMap = make(map[string]interface{})
            for col_counter,col:=range row{
                    singleMap[header_name[col_counter]] = col
            }
            if(len(singleMap)>0){
                
                parsedData = append(parsedData, singleMap)
            }
        }
    }
    fmt.Println("Length of parsedData:",len(parsedData))
    return parsedData
   
}

func ReadXlsxFile(filePath string)([]map[string]interface{}) {
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
            header_iterator:= header_name.Front()
            var singleMap = make(map[string]interface{})
            
            for _, cell := range row.Cells {
                if(row_counter==0){
                            text := cell.String()
                            header_name.PushBack(text)
                    }else{
                            text := cell.String()
                            singleMap[header_iterator.Value.(string)] = text
                            header_iterator = header_iterator.Next()
                    }    
            }
            if(row_counter != 0 && len(singleMap)>0){
            	
                parsedData = append(parsedData, singleMap)
            }
            
        }
    }
    fmt.Println("Length of parsedData:",len(parsedData))
    return parsedData
}




func ExcelCsvParser(blobPath string,blobExtension string){
	fmt.Println("---------------> We are in product.go")
	if blobExtension == ".csv"{
			fmt.Println("-------We are parsing an csv file.-------------")
			parsedData := ReadCsvFile(blobPath)
            parsedJson, _:= json.Marshal(parsedData)
            fmt.Println(string(parsedJson))
		}else if blobExtension == ".xlsx"{
			fmt.Println("----------------We are parsing an xlsx file.---------------")
		    parsedData := ReadXlsxFile(blobPath)
            parsedJson, _:= json.Marshal(parsedData)
            fmt.Println(string(parsedJson))
            err := ioutil.WriteFile("output.json",parsedJson, 0644)
            _ = err
        }



}

func uploadData(w http.ResponseWriter, req *http.Request){
    params := mux.Vars(req)
    file, err := os.Create("./data/"+params["fileName"])
    _, err = io.Copy(file, req.Body)
    if err!=nil{
      log.Printf("Error while Posting data")
    }
    blobPath := "./data/"+params["fileName"]
    var extension = filepath.Ext(blobPath)
    ExcelCsvParser(blobPath,extension)
}



func main(){
  router := mux.NewRouter()
  router.HandleFunc("/upload/{fileName}", uploadData).Methods("POST")
  log.Fatal(http.ListenAndServe(":8000", router))
}