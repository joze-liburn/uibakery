package cryptlex

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"
)

type (
	// https://api.cryptlex.com/v3/docs#tag/Resellers/operation/GetReseller
	Reseller struct {
		Id                   string    `json:"id"`
		Name                 string    `json:"name"`
		Description          string    `json:"description"`
		Email                string    `json:"email"`
		AllowedOrganizations int       `json:"allowedOrganizations"`
		AllowedUsers         int       `json:"allowedUsers"`
		CreatedAt            time.Time `json:"createdAt"`
		UpdatedAt            time.Time `json:"updatedAt"`
	}
)

var (
	errResellerJson = errors.New("reseller json error")
	errResellerHttp = errors.New("cryptlex service error")
)

func jsonToError(body []byte) (Reseller, error) {
	errb := struct {
		Mesage string `json:"message"`
		Code   string `json:"code,omitempty"`
	}{}
	err := json.Unmarshal(body, &errb)
	if err != nil {
		return Reseller{}, fmt.Errorf("%w: %s", errResellerJson, err)
	}
	if len(errb.Code)+len(errb.Mesage) > 0 {
		return Reseller{}, fmt.Errorf("%w: %s (%s)", errResellerHttp, errb.Code, errb.Mesage)
	}
	return Reseller{}, nil
}

func jsonToReseller(body []byte) (Reseller, error) {
	rsl := Reseller{}
	err := json.Unmarshal(body, &rsl)
	if err != nil {
		return Reseller{}, fmt.Errorf("%w: %s", errResellerJson, err)
	}
	return rsl, nil
}

func (lex *Cryptlex) RetrieveReseller(id string) (Reseller, error) {
	body, code, err := lex.Get(fmt.Sprintf("resellers/%s", id))
	if err != nil {
		return Reseller{}, err
	}
	if code != http.StatusOK {
		return jsonToError(body)
	}
	return jsonToReseller(body)
}
