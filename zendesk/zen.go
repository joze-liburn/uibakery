package zendesk

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

type (
	Zendesk struct {
		zdProtocol string
		zdHost     string
		zdApi      string
		token      string
	}

	GetUrl struct {
		page  int
		after string
		extId string
		name  string
	}

	GetOptions func(*GetUrl) error
)

func NewZendesk(zdHost string, token string) *Zendesk {
	return &Zendesk{zdProtocol: "https://", zdHost: zdHost, zdApi: "api/v2", token: token}
}

func (gu *GetUrl) Url(endpoint string) string {
	url := endpoint
	params := []string{}
	if gu.page > 0 {
		params = append(params, fmt.Sprintf("page[size]=%d", gu.page))
		if gu.after != "" {
			params = append(params, fmt.Sprintf("page[after]=%s", gu.after))
		}
	}
	if len(gu.extId) > 0 {
		params = append(params, fmt.Sprintf("external_id=%s", gu.extId))
	}
	if len(gu.name) > 0 {
		params = append(params, fmt.Sprintf("name=%s", gu.name))
	}
	if len(params) > 0 {
		url = url + "?" + strings.Join(params, "&")
	}
	return url
}

func WithPage(ps int) GetOptions {
	return func(gu *GetUrl) error {
		gu.page = ps
		return nil
	}
}

func StartAfter(af string) GetOptions {
	return func(gu *GetUrl) error {
		gu.after = af
		return nil
	}
}

func ByExternalId(eid string) GetOptions {
	return func(gu *GetUrl) error {
		gu.extId = eid
		return nil
	}
}

func ByName(nm string) GetOptions {
	return func(gu *GetUrl) error {
		gu.name = nm
		return nil
	}
}

func geturl(protocol, host, api, endpoint string, opt ...GetOptions) string {
	gu := &GetUrl{}
	for _, f := range opt {
		if err := f(gu); err != nil {
			return ""
		}
	}
	path, _ := url.JoinPath(protocol, host, api, endpoint)
	return gu.Url(path)
}

func (zd *Zendesk) Get(api string, opt ...GetOptions) ([]byte, error) {
	url := geturl(zd.zdProtocol, zd.zdHost, zd.zdApi, api, opt...)
	method := http.MethodGet
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", zd.token))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	return body, nil
}

func (zd *Zendesk) Put(api string, payload []byte, opt ...GetOptions) ([]byte, error) {
	url := geturl(zd.zdProtocol, zd.zdHost, zd.zdApi, api, opt...)
	method := http.MethodPut
	req, _ := http.NewRequest(method, url, bytes.NewBuffer(payload))
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", zd.token))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return []byte{}, err
	}
	defer res.Body.Close()

	response, err := io.ReadAll(res.Body)
	if err != nil {
		return []byte{}, err
	}
	return response, nil
}
