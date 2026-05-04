package file

import(
    "os"
    "fmt"
    "encoding/csv"
    data "go-crawler/data"
    
)

func WriteFileWithJsonResult(jsonFinalParsingResult []byte) {
	f, err := os.Create("results.json")

	if err != nil {
		panic(err)
	}

	resultCreationFile, err := f.WriteString(string(jsonFinalParsingResult))

	if err != nil {
		panic(err)
	}

	fmt.Println(resultCreationFile)
}

func ReceiveResults(chFinalUrlParsingReceiver <-chan data.Result, resultsFinalParsing *[]data.Result) {
	result := <-chFinalUrlParsingReceiver
	*resultsFinalParsing = append(*resultsFinalParsing, result)
}

func ReadFile(path string, targetList *[]data.Target) {
    file, err := os.Open(path)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer file.Close()

    reader := csv.NewReader(file)
    records, err := reader.ReadAll()
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    for _, record := range records {
        *targetList = append(*targetList, data.Target{record[0]})
    }
}