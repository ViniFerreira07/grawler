package main

import (
	"encoding/json"
	"fmt"
	"sync"
	"flag"
	data "go-crawler/data"
	validations "go-crawler/validations"
	file "go-crawler/file"
	http "go-crawler/http"
)

func main() {
	var wg sync.WaitGroup
	var resultsFinalParsing []data.Result
	var targetList []data.Target

	validations := []validations.Validation{validations.ValidateNotEmptyUrl{} ,validations.ValidateUrlLength{}}

	workersFlag := flag.Int("workers", 30, "Numbers of workers for URL parsing")
	fileFlag := flag.String("file", "url.csv", "File with the URLs")
	tagFlag := flag.String("tag", "title", "Tag to search in the URLs")

	flag.Parse()

	file.ReadFile(*fileFlag, &targetList)

	var numJobs = len(targetList)
	numWorkers := *workersFlag  
	tag := *tagFlag

	chFinalUrlParsingSender := make(chan string, numJobs)
	chFinalUrlParsingReceiver := make(chan data.Result, numJobs)	

	for i := 0; i <= numWorkers; i++ {
		wg.Go(func() {
			http.GetUrlHttpBody(chFinalUrlParsingSender, chFinalUrlParsingReceiver, tag)
		}) 
	}

	for i := 0; i < numJobs; i++ {
		http.ReceiveTargetUrl(targetList[i], validations, chFinalUrlParsingSender)
	}

	close(chFinalUrlParsingSender)

	go func() {
		wg.Wait()
		close(chFinalUrlParsingReceiver)
	}()

	for result := range chFinalUrlParsingReceiver { 
		fmt.Println(result)
		resultsFinalParsing = append(resultsFinalParsing, result)
	}

	jsonFinalParsingResult, _ := json.MarshalIndent(resultsFinalParsing, " ", " ")
	file.WriteFileWithJsonResult(jsonFinalParsingResult)
}