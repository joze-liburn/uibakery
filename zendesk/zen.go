package zendesk

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
)

type (
	Zendesk struct {

		zdApi string
		token string
	}

	GetUrl struct {
		page  int
		after string
		extId int
		name  string
	}

	GetOptions func(*GetUrl) error
)

func NewZendesk(zdAddr string, token string) *Zendesk {
	return &Zendesk{zdApi: zdAddr, token: token}
}

func (gu *GetUrl) Url(endpoint string) string {
	url := endpoint
	if gu.page > 0 {
		url = fmt.Sprintf("%s&page[size]=%d", url, gu.page)
		if gu.after != "" {
			url = fmt.Sprintf("%s&page[after]=%s", url, gu.after)
		}
	}
	if gu.extId > 0 {
		url = fmt.Sprintf("%s&external_id=%d", url, gu.extId)
	}
	if len(gu.name) > 0 {
		url = fmt.Sprintf("%s&name=%s", url, gu.name)
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

func ByExternalId(eid int) GetOptions {
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

func geturl(base, endpoint string, opt ...GetOptions) string {
	gu := &GetUrl{}
	for _, f := range opt {
		if err := f(gu); err != nil {
			return ""
		}
	}
	path, _ := url.JoinPath(base, endpoint)
	return gu.Url(path)
}

func (zd *Zendesk) Get(api string, opt ...GetOptions) ([]byte, error) {
	url := geturl(zd.zdApi, api, opt...)
	method := http.MethodGet
	req, _ := http.NewRequest(method, url, nil)
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Basic %s", zd.token))

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
	defer res.Body.Close()

	body, err := io.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return []byte{}, err
	}
	return body, nil
}
