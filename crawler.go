package main

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
	"flag"
	"encoding/csv"
)

type Target struct {
	url string
}

type Result struct {
	Url string `json:"url"`
	StatusCode int `json:"statusCode"`
	Result string `json:"result"`
}

type Validation interface {
	Validate(url string) (bool, error)
}

type validateNotEmptyUrl struct{}

func (validateNotEmptyUrl) Validate(url string) (bool, error) {
	if url == "" {
		return false, errors.New("Empty URL")
	}

	return true, nil
}

type validateUrlLength struct{}

func (validateUrlLength) Validate(url string) (bool, error) {
	if len(url) > 100 {
		return false, errors.New("Oversized URL - Length: " + strconv.Itoa(len(url)))
	}

	return true, nil
}

func validateUrl(url string, validations []Validation) (bool, string) {
	var errs string
	statusFinal := true

	for _, exec := range validations {
		if status, e := exec.Validate(url); !status {
			errs = errs + e.Error() + "; "
			statusFinal = false
		}
	}

	return statusFinal, errs
}

func checkUrlStartingHttps(url string) (string) {
	if !strings.HasPrefix(url, "https://") && strings.Contains(url, ".com") {
	 	url = "https://" + url
	}

	return url
}

func getUrlHttpBody(chFinalUrlParsingSender <-chan string, chFinalUrlParsingReceiver chan<- Result, tag string) {	
	for url := range chFinalUrlParsingSender {
		client := &http.Client{
			Timeout: 3 * time.Second,
		}

		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0")

		response, err := client.Do(req)

		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(response.Body)

		if !scanner.Scan() {
			continue
		}

		resultString := ""

		for i := 0; scanner.Scan() && resultString == ""; i++ {
			resultString = getStringBetweenTags(scanner.Text(), "<"+tag+">", "</"+tag+">")
		}

		chFinalUrlParsingReceiver <- Result{url, response.StatusCode, resultString}

		response.Body.Close()
	}
}

func urlProcessWorker(validations []Validation, chFinalUrlParsingSender <-chan Target, chFinalUrlParsingReceiver chan<- Result) {	
	
}

func getStringBetweenTags(requestTextBody string, startTag string, endTag string) (string) {
	startTagIndex := strings.Index(requestTextBody, startTag)

	if startTagIndex == -1 {
		return ""
	}

	startTagIndex += len(startTag)

	endTagIndex := strings.Index(requestTextBody[startTagIndex:], endTag)

	if endTagIndex == -1 {
		return ""
	}

	endTagIndex = startTagIndex + endTagIndex

	if endTagIndex < startTagIndex {
		return ""
	}

	return requestTextBody[startTagIndex:endTagIndex]
}

func receiveTargetUrl(target Target, validations []Validation, chFinalParsingSender chan<- string) {
	statusFinal, errs := validateUrl(target.url, validations)

	errs = func() string {
		if errs == "" {
			return "No erros"
		} else {
			return errs
		}
	}()
	
	if statusFinal {
		url := checkUrlStartingHttps(target.url)
		chFinalParsingSender <- url
	}
}

func writeFileWithJsonResult(jsonFinalParsingResult []byte) {
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

func receiveResults(chFinalUrlParsingReceiver <-chan Result, resultsFinalParsing *[]Result) {
	result := <-chFinalUrlParsingReceiver
	*resultsFinalParsing = append(*resultsFinalParsing, result)
}

func readFile(path string, targetList *[]Target) {
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
        *targetList = append(*targetList, Target{record[0]})
    }
}

func main() {
	var wg sync.WaitGroup
	var resultsFinalParsing []Result
	var targetList []Target

	validations := []Validation{validateNotEmptyUrl{} ,validateUrlLength{}}

	workersFlag := flag.Int("workers", 30, "Numbers of workers for URL parsing")
	fileFlag := flag.String("file", "url.csv", "File with the URLs")
	tagFlag := flag.String("tag", "title", "Tag to search in the URLs")

	flag.Parse()

	readFile(*fileFlag, &targetList)

	var numJobs = len(targetList)
	numWorkers := *workersFlag  
	tag := *tagFlag

	chFinalUrlParsingSender := make(chan string, numJobs)
	chFinalUrlParsingReceiver := make(chan Result, numJobs)	

	for i := 0; i <= numWorkers; i++ {
		wg.Go(func() {
			getUrlHttpBody(chFinalUrlParsingSender, chFinalUrlParsingReceiver, tag)
		}) 
	}

	for i := 0; i < numJobs; i++ {
		receiveTargetUrl(targetList[i], validations, chFinalUrlParsingSender)
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
	writeFileWithJsonResult(jsonFinalParsingResult)
}