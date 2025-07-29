package shopify

import (
	"encoding/json"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestFromBody(t *testing.T) {
	tests := []struct {
		name    string
		data    map[string]any
		want    CompanyConnection
		wanterr error
	}{
		{
			name: "normal",
			data: map[string]any{
				"companies": map[string]any{
					"nodes": []any{
						map[string]any{"id": "gid://shopify/Company/001"},
						map[string]any{"id": "gid://shopify/Company/002"},
						map[string]any{"id": "gid://shopify/Company/003"},
					},
				},
			},
			want: CompanyConnection{
				Nodes: []Company{
					{Id: "gid://shopify/Company/001"},
					{Id: "gid://shopify/Company/002"},
					{Id: "gid://shopify/Company/003"},
				},
			},
		},
		{
			name: "with page",
			data: map[string]any{
				"companies": map[string]any{
					"nodes": []any{
						map[string]any{"id": "gid://shopify/Company/001"},
						map[string]any{"id": "gid://shopify/Company/002"},
						map[string]any{"id": "gid://shopify/Company/003"},
					},
					"pageInfo": map[string]any{
						"endCursor":   "endCursor",
						"hasNextPage": true,
					},
				},
			},
			want: CompanyConnection{
				Nodes: []Company{
					{Id: "gid://shopify/Company/001"},
					{Id: "gid://shopify/Company/002"},
					{Id: "gid://shopify/Company/003"},
				},
				PageInfo: &PageInfo{HasNextPage: true, EndCursor: "endCursor"},
			},
		},
		{
			name: "no companies",
			data: map[string]any{
				"podjetja": map[string]any{
					"nodes": []any{
						map[string]any{"id": "gid://shopify/Company/001"},
						map[string]any{"id": "gid://shopify/Company/002"},
						map[string]any{"id": "gid://shopify/Company/003"},
					},
				},
			},
			wanterr: errInput,
		},
		{
			name: "bad companies",
			data: map[string]any{
				"companies": []string{
					"gid://shopify/Company/001",
					"gid://shopify/Company/002",
					"gid://shopify/Company/003",
				},
			},
			wanterr: errInput,
		},
		{
			name: "no nodes",
			data: map[string]any{
				"companies": map[string]any{
					"vozlišča": []any{
						map[string]any{"id": "gid://shopify/Company/001"},
						map[string]any{"id": "gid://shopify/Company/002"},
						map[string]any{"id": "gid://shopify/Company/003"},
					},
				},
			},
			want: CompanyConnection{},
		},
		{
			name: "bad nodes",
			data: map[string]any{
				"companies": map[string]any{
					"nodes": []string{
						"gid://shopify/Company/001",
						"gid://shopify/Company/002",
						"gid://shopify/Company/003",
					},
				},
			},
			wanterr: errInput,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got, err := companiesFromGQL(test.data)
			if !errors.Is(err, test.wanterr) {
				t.Errorf("%s: got error %v, want %v", test.name, err, test.wanterr)
			}
			if df := cmp.Diff(test.want, got); df != "" {
				t.Errorf("%s: -want +got\n%s", test.name, df)
			}
		})
	}
}

func TestGetCompanyDetails(t *testing.T) {
	data := map[string]any{
		"contactCount":       4,
		"createdAt":          "2024-04-24T20:15:43Z",
		"customerSince":      "2024-04-24T20:15:43Z",
		"externalId":         "LBResellerTest",
		"hasTimelineComment": false,
		"id":                 "gid://shopify/Company/309264552",
		"lifetimeDuration":   "about 1 year",
		"name":               "LightBurn Staff Internal Test Company B2B",
		"note":               nil,
		"updatedAt":          "2024-11-14T00:55:50Z",
		"locationsCount": map[string]any{
			"count": 1,
		},
		"locations": map[string]any{
			"nodes": []map[string]any{
				{
					"createdAt":          "2024-04-24T20:15:43Z",
					"currency":           "USD",
					"defaultCursor":      "eyJkaXJlY3Rpb24iOiJuZXh0IiwibGFzdF9pZCI6MTI0MDE3MDY2NCwibGFzdF92YWx1ZSI6IjEyNDAxNzA2NjQiLCJsaW1pdCI6MSwic29ydF9kaXIiOiJhc2MiLCJzb3J0X2ZpZWxkcyI6ImlkIn0=",
					"externalId":         nil,
					"hasTimelineComment": false,
					"id":                 "gid://shopify/CompanyLocation/1240170664",
					"locale":             "en",
					"name":               "10 Falls Drive",
					"note":               nil,
					"orderCount":         2,
					"phone":              nil,
					"taxExemptions":      nil,
					"taxRegistrationId":  nil,
					"updatedAt":          "2024-09-30T16:03:37Z",
					"billingAddress": map[string]any{
						"address1":         "10 Falls Drive",
						"address2":         nil,
						"city":             "Brookfield",
						"companyName":      "LightBurn Staff Internal Test Company B2B",
						"country":          "United States",
						"countryCode":      "US",
						"createdAt":        "2024-04-24T20:15:43Z",
						"firstName":        "Light",
						"formattedAddress": []string{"LightBurn Staff Internal Test Company B2B", "10 Falls Drive", "Brookfield CT 06804", "United States"},
						"formattedArea":    "Brookfield CT, United States",
						"id":               "gid://shopify/CompanyAddress/1326710952",
						"lastName":         "Burn",
						"phone":            "+17782370957",
						"province":         "Connecticut",
						"recipient":        "LB",
						"updatedAt":        "2024-04-24T20:15:43Z",
						"zip":              "06804",
						"zoneCode":         "CT",
					},
				},
			},
		},
		"contacts": map[string]any{
			"nodes": []map[string]any{
				{
					"createdAt":        "2024-04-24T20:21:00Z",
					"id":               "gid://shopify/CompanyContact/225607848",
					"isMainContact":    false,
					"lifetimeDuration": "about 1 year",
					"locale":           "en-CA",
					"title":            nil,
					"updatedAt":        "2024-05-16T13:29:07Z",
					"customer": map[string]any{
						"canDelete":               false,
						"createdAt":               "2022-11-17T18:07:36Z",
						"dataSaleOptOut":          false,
						"displayName":             "Colin Worobetz",
						"email":                   "colin@lightburnsoftware.com",
						"firstName":               "Colin",
						"hasTimelineComment":      false,
						"id":                      "gid://shopify/Customer/6384650780840",
						"lastName":                "Worobetz",
						"legacyResourceId":        "6384650780840",
						"lifetimeDuration":        "over 2 years",
						"locale":                  "en-US",
						"multipassIdentifier":     nil,
						"note":                    nil,
						"numberOfOrders":          "5",
						"phone":                   nil,
						"productSubscriberStatus": "NEVER_SUBSCRIBED",
						"state":                   "ENABLED",
						"tags": []string{
							"Login with Shop",
							"Shop",
						},
						"updatedAt":         "2025-04-02T05:03:50Z",
						"validEmailAddress": true,
						"verifiedEmail":     true,
					},
				},
				{
					"createdAt":        "2024-05-27T20:18:13Z",
					"id":               "gid://shopify/CompanyContact/332497064",
					"isMainContact":    false,
					"lifetimeDuration": "about 1 year",
					"locale":           "en-CA",
					"title":            nil,
					"updatedAt":        "2024-05-27T20:18:13Z",
					"customer": map[string]any{
						"canDelete":               false,
						"createdAt":               "2022-07-11T18:46:32Z",
						"dataSaleOptOut":          false,
						"displayName":             "Wayne Pearson",
						"email":                   "wayne@lightburnsoftware.com",
						"firstName":               "Wayne",
						"hasTimelineComment":      false,
						"id":                      "gid://shopify/Customer/6178271985832",
						"lastName":                "Pearson",
						"legacyResourceId":        "6178271985832",
						"lifetimeDuration":        "about 3 years",
						"locale":                  "en-CA",
						"multipassIdentifier":     nil,
						"note":                    "VAT Registration Number: 8942806842",
						"numberOfOrders":          "46",
						"phone":                   nil,
						"productSubscriberStatus": "NEVER_SUBSCRIBED",
						"state":                   "ENABLED",
						"tags":                    []string{"VAT:VAT8675309"},
						"updatedAt":               "2025-06-05T20:32:40Z",
						"validEmailAddress":       true,
						"verifiedEmail":           true,
					},
				},
			},
		},
		"metafields": map[string]any{
			"nodes": []map[string]any{
				{
					"key":       "do_not_sync_to_zendesk",
					"type":      "boolean",
					"updatedAt": "2024-11-14T00:55:50Z", "value": "true",
				},
			},
		},
	}

	c, _ := companyFromGQL(data)
	if c.Name != "LightBurn Staff Internal Test Company B2B" {
		t.Errorf("Name")
	}
	if c.ContactCount != 4 {
		t.Errorf("ContactCount")
	}
	if len(c.Locations.Nodes) != 1 {
		t.Errorf("len(c.Locations.Nodes)")
	}
	if c.Locations.Nodes[0].Name != "10 Falls Drive" {
		t.Errorf("len(c.Locations[0].Name])")
	}
	j, _ := json.MarshalIndent(c, "", "  ")
	t.Log(string(j))
	t.Log(c.Other)
}

/*
{
	"contactCount":4,
	"createdAt":"2024-04-24T20:15:43Z",
	"customerSince":"2024-04-24T20:15:43Z",
	"externalId":"LBResellerTest",
	"hasTimelineComment":false,
	"id":"gid://shopify/Company/309264552",
	"lifetimeDuration":"about 1 year",
	"name":"LightBurn Staff Internal Test Company B2B",
	"note":null,
	"updatedAt":"2024-11-14T00:55:50Z",
	"locationsCount":{
		"count":1
	},
	"locations":{
		"nodes":[
			{
				"createdAt":"2024-04-24T20:15:43Z",
				"currency":"USD",
				"defaultCursor":"eyJkaXJlY3Rpb24iOiJuZXh0IiwibGFzdF9pZCI6MTI0MDE3MDY2NCwibGFzdF92YWx1ZSI6IjEyNDAxNzA2NjQiLCJsaW1pdCI6MSwic29ydF9kaXIiOiJhc2MiLCJzb3J0X2ZpZWxkcyI6ImlkIn0=",
				"externalId":null,
				"hasTimelineComment":false,
				"id":"gid://shopify/CompanyLocation/1240170664",
				"locale":"en",
				"name":"10 Falls Drive",
				"note":null,
				"orderCount":2,
				"phone":null,
				"taxExemptions":[],
				"taxRegistrationId":null,
				"updatedAt":"2024-09-30T16:03:37Z",
				"billingAddress":{
					"address1":"10 Falls Drive",
					"address2":null,
					"city":"Brookfield",
					"companyName":"LightBurn Staff Internal Test Company B2B",
					"country":"United States",
					"countryCode":"US",
					"createdAt":"2024-04-24T20:15:43Z",
					"firstName":"Light",
					"formattedAddress":["LightBurn Staff Internal Test Company B2B","10 Falls Drive","Brookfield CT 06804","United States"],
					"formattedArea":"Brookfield CT, United States",
					"id":"gid://shopify/CompanyAddress/1326710952",
					"lastName":"Burn",
					"phone":"+17782370957",
					"province":"Connecticut",
					"recipient":"LB",
					"updatedAt":"2024-04-24T20:15:43Z",
					"zip":"06804",
					"zoneCode":"CT"
				}
			}
		]
	},
	"contacts":{
		"nodes":[
			{
				"createdAt":"2024-04-24T20:21:00Z",
				"id":"gid://shopify/CompanyContact/225607848",
				"isMainContact":false,
				"lifetimeDuration":"about 1 year",
				"locale":"en-CA",
				"title":null,
				"updatedAt":"2024-05-16T13:29:07Z",
				"customer":{
					"canDelete":false,
					"createdAt":"2022-11-17T18:07:36Z",
					"dataSaleOptOut":false,
					"displayName":"Colin Worobetz",
					"email":"colin@lightburnsoftware.com",
					"firstName":"Colin",
					"hasTimelineComment":false,
					"id":"gid://shopify/Customer/6384650780840",
					"lastName":"Worobetz",
					"legacyResourceId":"6384650780840",
					"lifetimeDuration":"over 2 years",
					"locale":"en-US",
					"multipassIdentifier":null,
					"note":null,
					"numberOfOrders":"5",
					"phone":null,
					"productSubscriberStatus":"NEVER_SUBSCRIBED",
					"state":"ENABLED",
					"tags":[
						"Login with Shop",
						"Shop"
					],
					"updatedAt":"2025-04-02T05:03:50Z",
					"validEmailAddress":true,
					"verifiedEmail":true
				}
			},
			{
				"createdAt":"2024-05-27T20:18:13Z",
				"id":"gid://shopify/CompanyContact/332497064",
				"isMainContact":false,
				"lifetimeDuration":"about 1 year",
				"locale":"en-CA",
				"title":null,
				"updatedAt":"2024-05-27T20:18:13Z",
				"customer":{
					"canDelete":false,
					"createdAt":"2022-07-11T18:46:32Z",
					"dataSaleOptOut":false,
					"displayName":"Wayne Pearson",
					"email":"wayne@lightburnsoftware.com",
					"firstName":"Wayne",
					"hasTimelineComment":false,
					"id":"gid://shopify/Customer/6178271985832",
					"lastName":"Pearson",
					"legacyResourceId":"6178271985832",
					"lifetimeDuration":"about 3 years",
					"locale":"en-CA",
					"multipassIdentifier":null,
					"note":"VAT Registration Number: 8942806842",
					"numberOfOrders":"46",
					"phone":null,
					"productSubscriberStatus":"NEVER_SUBSCRIBED",
					"state":"ENABLED",
					"tags":[
						"VAT:VAT8675309"
					],
					"updatedAt":"2025-06-05T20:32:40Z",
					"validEmailAddress":true,
					"verifiedEmail":true
				}
			},
			{
				"createdAt":"2024-05-27T20:18:36Z",
				"id":"gid://shopify/CompanyContact/332529832",
				"isMainContact":false,
				"lifetimeDuration":"about 1 year",
				"locale":"en-US",
				"title":null,
				"updatedAt":"2024-05-27T20:18:37Z",
				"customer":{
					"canDelete":false,
					"createdAt":"2021-09-23T17:16:48Z",
					"dataSaleOptOut":false,
					"displayName":"Adam Haile",
					"email":"adam@lightburnsoftware.com",
					"firstName":"Adam",
					"hasTimelineComment":false,
					"id":"gid://shopify/Customer/5671908376744",
					"lastName":"Haile",
					"legacyResourceId":"5671908376744",
					"lifetimeDuration":"almost 4 years",
					"locale":"en-US",
					"multipassIdentifier":null,
					"note":null,
					"numberOfOrders":"26","phone":null,"productSubscriberStatus":"NEVER_SUBSCRIBED","state":"ENABLED","tags":[],"updatedAt":"2024-04-22T14:57:13Z","validEmailAddress":true,"verifiedEmail":true}},
			{
				"createdAt":"2024-05-27T20:19:36Z",
				"id":"gid://shopify/CompanyContact/332562600",
				"isMainContact":false,
				"lifetimeDuration":"about 1 year",
				"locale":"en-CA",
				"title":null,
				"updatedAt":"2024-05-27T20:19:36Z",
				"customer":{
					"canDelete":false,
					"createdAt":"2022-03-28T14:38:20Z",
					"dataSaleOptOut":false,
					"displayName":"Alexander McKinnon",
					"email":"chuck@lightburnsoftware.com",
					"firstName":"Alexander",
					"hasTimelineComment":false,
					"id":"gid://shopify/Customer/6032336748712",
					"lastName":"McKinnon",
					"legacyResourceId":"6032336748712",
					"lifetimeDuration":"over 3 years",
					"locale":"en-CA",
					"multipassIdentifier":null,
					"note":null,
					"numberOfOrders":"4","phone":null,"productSubscriberStatus":"NEVER_SUBSCRIBED","state":"ENABLED","tags":["LBXEvent","Purchaser"],"updatedAt":"2025-01-07T20:30:11Z","validEmailAddress":true,"verifiedEmail":true}}
		]
	},
	"metafield":null
}
*/
