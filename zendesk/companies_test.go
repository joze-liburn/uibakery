package zendesk

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
)

const (
	api_organizations_organizations = `"organizations":[
		{
			"url":"https://lightburnsoftware.zendesk.com/api/v2/organizations/1.json",
			"id":1,
			"name":"LightBurn Software",
			"shared_tickets":true,
			"shared_comments":true,
			"external_id":null,
			"created_at":"2023-09-19T16:26:39Z",
			"updated_at":"2024-05-01T16:40:10Z",
			"domain_names":["digifable.com","lightburnsoftware.com","millmage.com"],
			"details":"Staff Organization",
			"notes":"",
			"group_id":11,
			"tags":["lightburn-staff-email"],
			"organization_fields":{
				"makerspace_discount_code":null,
				"org_reseller":false,
				"reseller_discount_code":null,
				"sync_shopify_company":false
			}
		},
		{
			"url":"https://lightburnsoftware.zendesk.com/api/v2/organizations/2.json",
			"id":2,
			"name":"3D Flying Bear",
			"shared_tickets":false,
			"shared_comments":false,
			"external_id":null,
			"created_at":"2023-11-17T21:32:53Z",
			"updated_at":"2024-02-08T22:05:24Z",
			"domain_names":["3dflyingbear.com"],
			"details":"",
			"notes":"",
			"group_id":11,
			"tags":["reseller","vendor"],
			"organization_fields":{
				"makerspace_discount_code":null,
				"org_reseller":true,
				"reseller_discount_code":"Fly1ngB3@r-Di$tro",
				"sync_shopify_company":false
			}
		}
	]`
	api_organizations_meta = `"meta": {
		"has_more":true,
		"after_cursor":"eyJvIjoiaWQiLCJ2IjoiYVJzYTdacWlFZ0FBIn0=",
		"before_cursor":"eyJvIjoiaWQiLCJ2IjoiYVJ2OUd6VTBFUUFBIn0="
	}`
	api_organizations_links = `"links":{
		"prev":"https://lightburnsoftware.zendesk.com/api/v2/organizations.json?page%5Bbefore%5D=eyJvIjoiaWQiLCJ2IjoiYVJ2OUd6VTBFUUFBIn0%3D&page%5Bsize%5D=2",
		"next":"https://lightburnsoftware.zendesk.com/api/v2/organizations.json?page%5Bafter%5D=eyJvIjoiaWQiLCJ2IjoiYVJzYTdacWlFZ0FBIn0%3D&page%5Bsize%5D=2"
	}`
)

func Time(year int, mon time.Month, day, h, m, s int, loc *time.Location) *time.Time {
	tm := time.Date(year, mon, day, h, m, s, 0, loc)
	return &tm
}

func TestGetOrganizations(t *testing.T) {
	tests := []struct {
		name string
		data string
		want OrganizationResult
	}{
		{
			name: "normal",
			data: "{" + api_organizations_organizations + "," + api_organizations_meta + "," + api_organizations_links + "}",
			want: OrganizationResult{
				Organizations: []Organization{
					{
						Url:                "https://lightburnsoftware.zendesk.com/api/v2/organizations/1.json",
						Id:                 1,
						Name:               "LightBurn Software",
						SharedTickets:      true,
						SharedComments:     true,
						CreatedAt:          Time(2023, time.September, 19, 16, 26, 39, time.UTC),
						UpdatedAt:          Time(2024, time.May, 1, 16, 40, 10, time.UTC),
						DomainNames:        []string{"digifable.com", "lightburnsoftware.com", "millmage.com"},
						Details:            "Staff Organization",
						Notes:              "",
						GroupId:            11,
						Tags:               []string{"lightburn-staff-email"},
						OrganizationFields: &OrganizationFields{},
					},
					{
						Url:            "https://lightburnsoftware.zendesk.com/api/v2/organizations/2.json",
						Id:             2,
						Name:           "3D Flying Bear",
						SharedTickets:  false,
						SharedComments: false,
						CreatedAt:      Time(2023, time.November, 17, 21, 32, 53, time.UTC),
						UpdatedAt:      Time(2024, time.February, 8, 22, 05, 24, time.UTC),
						DomainNames:    []string{"3dflyingbear.com"},
						Details:        "",
						Notes:          "",
						GroupId:        11,
						Tags:           []string{"reseller", "vendor"},
						OrganizationFields: &OrganizationFields{
							OrgReseller:          true,
							ResellerDiscountCode: "Fly1ngB3@r-Di$tro",
						},
					},
				},
				Meta: Meta{
					HasMore:      true,
					BeforeCursor: "eyJvIjoiaWQiLCJ2IjoiYVJ2OUd6VTBFUUFBIn0=",
					AfterCursor:  "eyJvIjoiaWQiLCJ2IjoiYVJzYTdacWlFZ0FBIn0=",
				},
				Links: Links{
					Prev: "https://lightburnsoftware.zendesk.com/api/v2/organizations.json?page%5Bbefore%5D=eyJvIjoiaWQiLCJ2IjoiYVJ2OUd6VTBFUUFBIn0%3D&page%5Bsize%5D=2",
					Next: "https://lightburnsoftware.zendesk.com/api/v2/organizations.json?page%5Bafter%5D=eyJvIjoiaWQiLCJ2IjoiYVJzYTdacWlFZ0FBIn0%3D&page%5Bsize%5D=2",
				},
			},
		},
		{
			name: "organizations",
			data: "{" + api_organizations_organizations + "}",
			want: OrganizationResult{
				Organizations: []Organization{
					{
						Url:                "https://lightburnsoftware.zendesk.com/api/v2/organizations/1.json",
						Id:                 1,
						Name:               "LightBurn Software",
						SharedTickets:      true,
						SharedComments:     true,
						CreatedAt:          Time(2023, time.September, 19, 16, 26, 39, time.UTC),
						UpdatedAt:          Time(2024, time.May, 1, 16, 40, 10, time.UTC),
						DomainNames:        []string{"digifable.com", "lightburnsoftware.com", "millmage.com"},
						Details:            "Staff Organization",
						Notes:              "",
						GroupId:            11,
						Tags:               []string{"lightburn-staff-email"},
						OrganizationFields: &OrganizationFields{},
					},
					{
						Url:            "https://lightburnsoftware.zendesk.com/api/v2/organizations/2.json",
						Id:             2,
						Name:           "3D Flying Bear",
						SharedTickets:  false,
						SharedComments: false,
						CreatedAt:      Time(2023, time.November, 17, 21, 32, 53, time.UTC),
						UpdatedAt:      Time(2024, time.February, 8, 22, 05, 24, time.UTC),
						DomainNames:    []string{"3dflyingbear.com"},
						Details:        "",
						Notes:          "",
						GroupId:        11,
						Tags:           []string{"reseller", "vendor"},
						OrganizationFields: &OrganizationFields{
							OrgReseller:          true,
							ResellerDiscountCode: "Fly1ngB3@r-Di$tro",
						},
					},
				},
			},
		},
		{
			name: "meta",
			data: "{" + api_organizations_meta + "}",
			want: OrganizationResult{
				Meta: Meta{
					HasMore:      true,
					BeforeCursor: "eyJvIjoiaWQiLCJ2IjoiYVJ2OUd6VTBFUUFBIn0=",
					AfterCursor:  "eyJvIjoiaWQiLCJ2IjoiYVJzYTdacWlFZ0FBIn0=",
				},
			},
		},
		{
			name: "links",
			data: "{" + api_organizations_links + "}",
			want: OrganizationResult{
				Links: Links{
					Prev: "https://lightburnsoftware.zendesk.com/api/v2/organizations.json?page%5Bbefore%5D=eyJvIjoiaWQiLCJ2IjoiYVJ2OUd6VTBFUUFBIn0%3D&page%5Bsize%5D=2",
					Next: "https://lightburnsoftware.zendesk.com/api/v2/organizations.json?page%5Bafter%5D=eyJvIjoiaWQiLCJ2IjoiYVJzYTdacWlFZ0FBIn0%3D&page%5Bsize%5D=2",
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := jsonToOrganizations([]byte(test.data))
			if err != nil {
				t.Errorf("%s: got error %s", test.name, err)
			}
			if df := cmp.Diff(test.want, got); df != "" {
				t.Errorf("%s: -want +got\n%s", test.name, df)
			}
		})
	}
}
