package cryptlex

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"strings"
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

func (lex *Cryptlex) geturl(api ...string) string {
	sch, host, ok := strings.Cut(lex.host, "://")
	if !ok {
		host = sch
		sch = "https"
	}
	tail := strings.Join(api, "/")
	return fmt.Sprintf("%s://%s/v%d/%s", sch, host, lex.major, tail)
}

func (lex *Cryptlex) Get(api ...string) ([]byte, int, error) {
	url := lex.geturl(api...)
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

func (lex *Cryptlex) postpatch(method string, payload []byte, api ...string) ([]byte, int, error) {
	url := lex.geturl(api...)
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", lex.token))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, http.StatusOK, err
	}
	defer res.Body.Close()

	response, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, http.StatusOK, err
	}
	return response, res.StatusCode, nil
}

func (lex *Cryptlex) Patch(api string, id string, payload []byte) ([]byte, int, error) {
	return lex.postpatch(http.MethodPatch, payload, api, id)
}

func (lex *Cryptlex) Post(api string, payload []byte) ([]byte, int, error) {
	return lex.postpatch(http.MethodPost, payload, api)
}
