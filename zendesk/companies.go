package zendesk

import (
	"encoding/json"
	"fmt"
	"time"
)

// https://developer.zendesk.com/api-reference/ticketing/organizations/organizations/
type (
	OrganizationFields struct {
		MakerspaceDiscountCode string `json:"makerspace_discount_code"`
		OrgReseller            bool   `json:"org_reseller"`
		ResellerDiscountCode   string `json:"reseller_discount_code"`
		SyncShopifyCompany     bool   `json:"sync_shopify_company"`
		Other                  map[string]any
	}
	Organization struct {
		CreatedAt          *time.Time          `json:"created_at"`
		Details            string              `json:"details"`
		DomainNames        []string            `json:"domain_names"`
		ExternalId         string              `json:"external_id"`
		GroupId            int64               `json:"group_id"`
		Id                 int64               `json:"id"`
		Name               string              `json:"name"`
		Notes              string              `json:"notes"`
		OrganizationFields *OrganizationFields `json:"organization_fields"`
		SharedComments     bool                `json:"shared_comments"`
		SharedTickets      bool                `json:"shared_tickets"`
		Tags               []string            `json:"tags"`
		UpdatedAt          *time.Time          `json:"updated_at"`
		Url                string              `json:"url"`
	}
	OrganizationResult struct {
		Organizations []Organization `json:"organizations"`
		Meta          Meta           `json:"meta"`
		Links         Links          `json:"links"`
	}
)

func jsonToOrganizations(js []byte) (OrganizationResult, error) {
	var or OrganizationResult
	if err := json.Unmarshal(js, &or); err != nil {
		return or, err
	}
	return or, nil
}

func (zd *Zendesk) GetOrganizations(page int) (OrganizationResult, error) {
	opts := []GetOptions{}
	if page > 0 {
		opts = append(opts, WithPage(page))
	}
	rsp, err := zd.Get("organizations", opts...)
	if err != nil {
		return OrganizationResult{}, err
	}
	return jsonToOrganizations(rsp)
}

type OrganizationError struct {
	Organization Organization
	Err          error
}

// StreamOrganizations provides a channel with successive Organizations. Every
// organization in the channel has non-trivial ExternalId.
func (zd *Zendesk) StreamOrganizations(pageSize int, maxCount uint) <-chan OrganizationError {
	opts := []GetOptions{}
	if pageSize > 0 {
		opts = append(opts, WithPage(pageSize))
	}
	out := make(chan OrganizationError)
	go func() {
		defer close(out)
		var count uint
		rsp, geterr := zd.Get("organizations", opts...)
		for {
			if geterr != nil {
				out <- OrganizationError{Err: geterr}
				return
			}
			or, err := jsonToOrganizations(rsp)
			if err != nil {
				fmt.Printf("----- %s -----\n", zd.zdApi)
				out <- OrganizationError{Err: err}
				return
			}
			for _, org := range or.Organizations {
				if len(org.ExternalId) == 0 {
					continue
				}
				out <- OrganizationError{Organization: org}
				count++
				if count >= maxCount {
					return
				}
			}
			if !or.Meta.HasMore {
				return
			}
			rsp, geterr = zd.Get("organizations", append(opts, StartAfter(or.Meta.AfterCursor))...)
		}
	}()
	return out
}

func (zd *Zendesk) GetOrganizationsByExternalId(extId string, page int) (OrganizationResult, error) {
	//	zd := NewZendesk(ZendeskApi, "574a0f524e9d4fb15bc6f678cf67f11ef442cd285d62c6b8f28397a996b7d37a")
	rsp, err := zd.Get("organizations")
	if err != nil {
		return OrganizationResult{}, err
	}
	return jsonToOrganizations(rsp)
}
