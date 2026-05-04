package data

type Target struct {
	Url string
}

type Result struct {
	Url string `json:"url"`
	StatusCode int `json:"statusCode"`
	Result string `json:"result"`
}