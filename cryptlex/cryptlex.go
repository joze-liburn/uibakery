package cryptlex

import (
	"fmt"
	"io"
	"net/http"
)

type (
	Cryptlex struct {
		major uint
		host  string
		token string
	}
)

func NewCryptlex(major uint, host string, secret string) *Cryptlex {
	return &Cryptlex{major: major, host: host, token: secret}
}

func (lex *Cryptlex) Get(api string) ([]byte, int, error) {
	url := fmt.Sprintf("https://%s/v%d/%s", lex.host, lex.major, api)
	method := http.MethodGet
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", lex.token))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, http.StatusOK, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, http.StatusOK, err
	}
	return body, res.StatusCode, nil
}
