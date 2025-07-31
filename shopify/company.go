package shopify

import (
	"errors"
	"fmt"
	"reflect"
	"slices"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/oussama4/gopify"
)

//
// after   String               The elements that come after the specified
//                              cursor.
// before  String               The elements that come before the specified
//                              cursor.
// first   Int                  The first n elements from the paginated list.
// last    Int                  The last n elements from the paginated list.
// query   String               A filter made up of terms, connectives,
//                              modifiers, and comparators. You can apply one or
//                              more filters to a query. Learn more about
//                              Shopify API search syntax.
// reverse Boolean (false)      Reverse the order of the underlying list.
// sortKey CompanySortKeys (ID) Sort the underlying list by the given key.

var (
	errInput   = errors.New("invalid input structure")
	errShopify = errors.New("shopify service error")
)

type (
	Node      string
	Companies struct {
		Companies CompanyConnection
	}

	CompanyConnection struct {
		Nodes    []Company `json:"nodes"`
		PageInfo *PageInfo `json:"pageInfo"`
	}

	// https://shopify.dev/docs/api/admin-graphql/2024-10/objects/CompanyLocation
	CompanyLocation struct {
		CreatedAt          time.Time      `json:"createdAt"`
		Currency           string         `json:"currency"`
		DefaultCursor      string         `json:"defaultCursor"`
		ExternalId         string         `json:"externalId"`
		HasTimelineComment bool           `json:"hasTimelineComment"`
		Id                 string         `json:"id"`
		Locale             string         `json:"locale"`
		Name               string         `json:"name"`
		Note               string         `json:"note"`
		OrderCount         int            `json:"orderCount"`
		Phone              string         `json:"phone"`
		TaxExemptions      []string       `json:"taxExemptions"`
		TaxRegistrationId  string         `json:"taxRegistrationId"`
		UpdatedAt          time.Time      `json:"updatedAt"`
		BillingAddress     Address        `json:"billingAddress"`
		Other              map[string]any `mapstructure:",remain"`
	}

	Locations struct {
		Nodes    []CompanyLocation `json:"nodes"`
		PageInfo *PageInfo         `json:"pageInfo"`
	}

	// https://shopify.dev/docs/api/admin-graphql/2024-10/objects/CompanyAddress
	Address struct {
		Address1         string    `json:"address1"`
		Address2         string    `json:"address2"`
		City             string    `json:"city"`
		CompanyName      string    `json:"companyName"`
		Country          string    `json:"country"`
		CountryCode      string    `json:"countryCode"`
		CreatedAt        time.Time `json:"createdAt"`
		FirstName        string    `json:"firstName"`
		FormattedAddress []string  `json:"formattedAddress"`
		FormattedArea    string    `json:"formattedArea"`
		Id               string    `json:"id"`
		LastName         string    `json:"lastName"`
		Phone            string    `json:"phone"`
		Province         string    `json:"province"`
		Recipient        string    `json:"recipient"`
		UpdatedAt        time.Time `json:"updatedAt"`
		Zip              string    `json:"zip"`
		ZoneCode         string    `json:"zoneCode"`
	}

	// https://shopify.dev/docs/api/admin-graphql/2024-10/objects/Company
	Company struct {
		ContactCount       int                 `json:"contactCount"`
		CreatedAt          time.Time           `json:"createdAt"`
		CustomerSince      time.Time           `json:"customerSince"`
		ExternalId         string              `json:"externalId"`
		HasTimelineComment bool                `json:"hasTimelineComment"`
		Id                 string              `json:"id"`
		LifetimeDuration   string              `json:"lifetimeDuration"`
		Name               string              `json:"name"`
		Note               string              `json:"note"`
		UpdatedAt          time.Time           `json:"updatedAt"`
		LocationsCount     Count               `json:"locationsCount"`
		Locations          Locations           `json:"locations"`
		Contacts           Contacts            `json:"contacts"`
		Other              map[string]any      `mapstructure:",remainOther"`
		Metafields         MetafieldConnection `json:"metafields"`
		Metafield          Metafield           `json:"metafield"`
	}

	/// DetailCompany is used to pass information on company over the
	// database based queue (compatibility with UI Bakery code).
	DetailCompany struct {
		Result                      Company        `json:"results"`
		Name                        string         `json:"name"`
		ExternalId                  string         `json:"external_id"`
		GroupId                     int            `json:"group_id"`
		NoSyncZendesk               bool           `json:"do_not_sync_to_zendesk"`
		Tags                        []string       `json:"tags"`
		Vendor                      bool           `json:"vendor"`
		Reseller                    bool           `json:"reseller"`
		DomainNames                 []string       `json:"domain_names"`
		OrganizationFields          map[string]any `json:"organization_fields"`
		Contacts                    []Contact      `json:"contacts"`
		PrimaryCompanyEmail         string         `json:"primary_company_email"`
		VatNumber                   string         `json:"vat_number,omitempty"`
		TaxSettings                 string         `json:"tax_setting,omitempty"`
		UlaDocument                 string         `json:"ula_document,omitempty"`
		LastSignedResellerAgreement string         `json:"last_signed_reseller_agreement_dt"`
		SharedTickets               bool           `json:"shared_tickets"`
		Checksum                    int64          `json:"checksum"`
	}

	SyncMetadata struct {
		CompanyName         string        `json:"company_name"`
		CompanyId           string        `json:"company_id"`
		Reseller            bool          `json:"reseller"`
		Vendor              bool          `json:"vendor"`
		PrimaryCompanyEmail string        `json:"primary_company_email"`
		Company             DetailCompany `json:"company"`
		Contact             Contact       `json:"contact"`
		ContactId           string        `json:"contact_id"`
		Customer            Customer      `json:"customer"`
		CustomerId          string        `json:"customer_id"`
		Checksum            int64         `json:"checksum"`
		NoSyncZendesk       bool          `json:"do_not_sync_to_zendesk"`
		// Unknown             jsontext.Value `json:",unknown"` // TBD: GO 1.25
	}
)

// GetMetafield returns the value of a given metafield. Second result indicates
// if the metafield was found.
func (c *Company) GetMetafield(key string) (string, bool) {
	for _, mf := range c.Metafields.Nodes {
		if mf.Key == key {
			return mf.Value, true
		}
	}
	return "", false
}

func stringToDateTimeHook(f reflect.Type, t reflect.Type, data any) (any, error) {
	if t == reflect.TypeOf(time.Time{}) && f == reflect.TypeOf("") {
		return time.Parse(time.RFC3339, data.(string))
	}

	return data, nil
}

// companiesFromGQL creates a list of companies' ids from the API result. This
// is expected to cointain "companies" on the top level map.
func companiesFromGQL(data map[string]any) (CompanyConnection, error) {
	var c Companies
	m := mapstructure.Metadata{Keys: []string{}}

	config := mapstructure.DecoderConfig{
		DecodeHook: stringToDateTimeHook,
		Metadata:   &m,
		Result:     &c,
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return CompanyConnection{}, fmt.Errorf("%w: %v", errInput, err)
	}
	err = decoder.Decode(data)
	if err != nil {
		return CompanyConnection{}, fmt.Errorf("%w: %v", errInput, err)
	}
	if slices.Contains(m.Unset, "Companies") {
		return CompanyConnection{}, fmt.Errorf("%w: missing companies", errInput)
	}
	return c.Companies, nil
}

// companyFromGQL extract Company data from the blob.
func companyFromGQL(data map[string]any) (Company, error) {
	var c Company

	config := mapstructure.DecoderConfig{
		DecodeHook: stringToDateTimeHook,
		Result:     &c,
	}
	decoder, err := mapstructure.NewDecoder(&config)
	if err != nil {
		return Company{}, fmt.Errorf("%w: %v", errInput, err)
	}
	err = decoder.Decode(data)
	if err != nil {
		return Company{}, fmt.Errorf("%w: %v", errInput, err)
	}
	return c, nil
}

// GetCompaniesIds obtains list of client ids from Shopify.
func (spfy *ShopifyOp) GetCompaniesIds(limit int) (CompanyConnection, error) {
	query := `
query Companies($pgSize: Int!, $cursor: String) {
  companies (first: $pgSize, after: $cursor) {
    nodes {
      id
    }
    pageInfo {
      hasNextPage
      endCursor
    }
  }
}
`
	gqlResult, err := spfy.client.Graphql(query, map[string]any{"pgSize": limit})
	if err != nil {
		return CompanyConnection{}, fmt.Errorf("%w: %v", errShopify, err)
	}

	return companiesFromGQL(gqlResult)
}

func GetCompanyDetails(client *gopify.Client) (Company, error) {
	query := `
query Company ($queryValue: ID!) {
    company(id: $queryValue) {
        contactCount
        createdAt
        customerSince
        externalId
        hasTimelineComment
        id
        lifetimeDuration
        name
        note
        updatedAt
        locationsCount {
            count
        }
        locations(first: 10) {
            nodes {
                createdAt
                currency
                defaultCursor
                externalId
                hasTimelineComment
                id
                locale
                name
                note
                orderCount
                phone
                taxExemptions
                taxRegistrationId
                updatedAt
                billingAddress {
                    address1
                    address2
                    city
                    companyName
                    country
                    countryCode
                    createdAt
                    firstName
                    formattedAddress
                    formattedArea
                    id
                    lastName
                    phone
                    province
                    recipient
                    updatedAt
                    zip
                    zoneCode
                }
            }
        }
        contacts(first: 100) {
            nodes {
                createdAt
                id
                isMainContact
                lifetimeDuration
                locale
                title
                updatedAt
                customer {
                    canDelete
                    createdAt
                    dataSaleOptOut
                    displayName
                    email
                    firstName
                    hasTimelineComment
                    id
                    lastName
                    legacyResourceId
                    lifetimeDuration
                    locale
                    multipassIdentifier
                    note
                    numberOfOrders
                    phone
                    productSubscriberStatus
                    state
                    tags
                    updatedAt
                    validEmailAddress
                    verifiedEmail
                }
            }
        }
        metafields(first: 10) {
            nodes {
                key
                type
                updatedAt
                value
            }
        }
        metafield(key: "vendor") {
            key
            type
            updatedAt
            value
        }
    }
}
`
	gqlResult, err := client.Graphql(query, nil)
	if err != nil {
		return Company{}, fmt.Errorf("%w: %v", errShopify, err)
	}

	return companyFromGQL(gqlResult)
}

// InitArrays fixes default array fields, from nil to empty array (except for
// OrganizationFields where it is populated withe defaults).
// Fields affected: Tags, DomainNames, OrganizationFields.
// If a field is not nil, it is not overwritten.
func (from DetailCompany) InitArrays() DetailCompany {
	if from.Tags == nil {
		from.Tags = []string{}
	}
	if from.DomainNames == nil {
		from.DomainNames = []string{}
	}
	if from.OrganizationFields == nil {
		from.OrganizationFields = map[string]any{"sync_shopify_company": false}
	}
	return from
}

// ToDetail makes baroque version of the Company structure for exchange through
// the database queue.
// This is the core body of actionShopifyGetCompanyData
func (c *Company) ToDetail() DetailCompany {
	ret := DetailCompany{}.InitArrays()
	ret.Result = *c
	ret.Name = c.Name
	ret.ExternalId = c.ExternalId
	ret.GroupId = 0
	ret.NoSyncZendesk = false
	ret.Vendor = false
	ret.Reseller = false
	ret.Contacts = c.Contacts.Nodes
	return ret
}
