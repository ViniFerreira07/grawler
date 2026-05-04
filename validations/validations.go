package validations

import(
	"errors"
)

type Validation interface {
	Validate(url string) (bool, error)
}

type ValidateNotEmptyUrl struct{}

func (ValidateNotEmptyUrl) Validate(url string) (bool, error) {
	if url == "" {
		return false, errors.New("Empty URL")
	}

	return true, nil
}

type ValidateUrlLength struct{}

func (ValidateUrlLength) Validate(url string) (bool, error) {
	if len(url) > 100 {
		return false, errors.New("Oversized URL - Length: " + string(len(url)))
	}

	return true, nil
}

func ValidateUrl(url string, validations []Validation) (bool, string) {
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
