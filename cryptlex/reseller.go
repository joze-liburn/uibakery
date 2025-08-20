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
		Id                   string    `json:"id,omitempty"`
		Name                 string    `json:"name"`
		Description          string    `json:"description"`
		Email                string    `json:"email"`
		AllowedOrganizations int       `json:"allowedOrganizations"`
		AllowedUsers         int       `json:"allowedUsers"`
		CreatedAt            time.Time `json:"createdAt,omitempty"`
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

func jsonToResellers(body []byte) ([]Reseller, error) {
	rsl := []Reseller{}
	err := json.Unmarshal(body, &rsl)
	if err != nil {
		return []Reseller{}, fmt.Errorf("%w: %s", errResellerJson, err)
	}
	return rsl, nil
}

// https://api.cryptlex.com/v3/docs#tag/Resellers/operation/GetReseller
func (lex *Cryptlex) RetrieveReseller(id string) (Reseller, error) {
	body, code, err := lex.Get("resellers", id)
	if err != nil {
		return Reseller{}, err
	}
	if code != http.StatusOK {
		return jsonToError(body)
	}
	return jsonToReseller(body)
}

// https://api.cryptlex.com/v3/docs#tag/Resellers/operation/GetAllResellers
func (lex *Cryptlex) ListResellers(opts ...CallOptions) ([]Reseller, error) {
	criteria := CryptlexOpts{}
	if err := criteria.Apply(opts...); err != nil {
		return []Reseller{}, err
	}

	body, code, err := lex.Get(fmt.Sprintf("resellers%s", criteria.ToSearchParam()))
	if err != nil {
		return []Reseller{}, err
	}
	if code != http.StatusOK {
		_, err := jsonToError(body)
		return []Reseller{}, err
	}
	return jsonToResellers(body)
}

// https://api.cryptlex.com/v3/docs#tag/Resellers/operation/CreateReseller
func (lex *Cryptlex) CreateReseller(r Reseller) (Reseller, error) {
	up := Reseller{
		Name:                 r.Name,
		Description:          r.Description,
		Email:                r.Email,
		AllowedOrganizations: r.AllowedOrganizations,
		AllowedUsers:         r.AllowedUsers,
	}
	b, err := json.Marshal(up)
	if err != nil {
		return Reseller{}, err
	}
	rsp, code, err := lex.Post("resellers", b)
	if err != nil {
		return Reseller{}, err
	}
	if code != http.StatusOK {
		_, err := jsonToError(rsp)
		return Reseller{}, err
	}
	return jsonToReseller(rsp)
}

// https://api.cryptlex.com/v3/docs#tag/Resellers/operation/UpdateReseller
func (lex *Cryptlex) UpdateReseller(r Reseller) (Reseller, error) {
	up := Reseller{
		Name:                 r.Name,
		Description:          r.Description,
		Email:                r.Email,
		AllowedOrganizations: r.AllowedOrganizations,
		AllowedUsers:         r.AllowedUsers,
	}
	b, err := json.Marshal(up)
	if err != nil {
		return Reseller{}, err
	}
	rsp, code, err := lex.Patch("resellers", r.Id, b)
	if err != nil {
		return Reseller{}, err
	}
	if code != http.StatusOK {
		_, err := jsonToError(rsp)
		return Reseller{}, err
	}
	return jsonToReseller(rsp)
}
