package parser

import (
	"container/list"
	"fmt"

	"github.com/extrame/xls"
	"github.com/tealeg/xlsx"
)

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

func ReadXlsFile(filePath string) []map[string]interface{} {
	parsedData := make([]map[string]interface{}, 0, 0)
	if xlFile, err := xls.Open(filePath, "utf-8"); err == nil {
		total_sheets := xlFile.NumSheets()
		for sheetCounter := 0; sheetCounter < total_sheets; sheetCounter++ {
			if sheet := xlFile.GetSheet(sheetCounter); sheet != nil {
				header_name := list.New()
				for rowCounter := 0; rowCounter <= (int(sheet.MaxRow)); rowCounter++ {
					row := sheet.Row(rowCounter)
					header_iterator := header_name.Front()
					var singleMap = make(map[string]interface{})
					for colCounter := 0; colCounter < (int(row.LastCol())); colCounter++ {
						if rowCounter == 0 {
							text := row.Col(colCounter)
							header_name.PushBack(text)
						} else {
							text := row.Col(colCounter)
							singleMap[header_iterator.Value.(string)] = text
							header_iterator = header_iterator.Next()
						}
					}
					if rowCounter != 0 && len(singleMap) > 0 {
						parsedData = append(parsedData, singleMap)
					}
				}
			}
		}
	}
	fmt.Println("Length of parsedData:", len(parsedData))
	return parsedData
}
