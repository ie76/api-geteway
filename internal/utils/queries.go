package utils

import (
	"assignment/internal/errors"
	"net/url"
)

func AddQueryParams(baseURL string, newParams map[string]string) (string, *errors.Error) {
	u, err := url.Parse(baseURL)
	if err != nil {
		return "", errors.New(errors.ErrParseUrl)
	}

	params := u.Query()

	for key, value := range newParams {
		params.Add(key, value)
	}

	u.RawQuery = params.Encode()
	return u.String(), nil
}
