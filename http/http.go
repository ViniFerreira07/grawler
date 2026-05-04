package http

import(
	"strings"
	"time"
	"net/http"
	"bufio"
	data "go-crawler/data"
	validations "go-crawler/validations"
)

func CheckUrlStartingHttps(url string) (string) {
	if !strings.HasPrefix(url, "https://") && strings.Contains(url, ".com") {
	 	url = "https://" + url
	}

	return url
}

func GetUrlHttpBody(chFinalUrlParsingSender <-chan string, chFinalUrlParsingReceiver chan<- data.Result, tag string) {	
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
			resultString = GetStringBetweenTags(scanner.Text(), "<"+tag+">", "</"+tag+">")
		}

		chFinalUrlParsingReceiver <- data.Result{url, response.StatusCode, resultString}

		response.Body.Close()
	}
}

func GetStringBetweenTags(requestTextBody string, startTag string, endTag string) (string) {
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

func ReceiveTargetUrl(target data.Target, validationsSlice []validations.Validation, chFinalParsingSender chan<- string) {
	statusFinal, errs := validations.ValidateUrl(target.Url, validationsSlice)

	errs = func() string {
		if errs == "" {
			return "No erros"
		} else {
			return errs
		}
	}()
	
	if statusFinal {
		url := CheckUrlStartingHttps(target.Url)
		chFinalParsingSender <- url
	}
}