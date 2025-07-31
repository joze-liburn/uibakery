package sync

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"gitlab.com/joze-liburn/uibakery/shopify"
)

func TestGetRecords(t *testing.T) {
	jose := shopify.Customer{
		CreatedAt:               time.Date(2025, time.May, 29, 15, 18, 10, 0, time.UTC),
		DisplayName:             "José Letona",
		Email:                   "info@creativo3d.com",
		FirstName:               "José",
		Id:                      "gid://shopify/Customer/8450627043496",
		LastName:                "Letona",
		LegacyResourceId:        "8450627043496",
		LifetimeDuration:        "13 days",
		Locale:                  "en-GT",
		Note:                    "",
		NumberOfOrders:          "1",
		Phone:                   "+50237614657",
		ProductSubscriberStatus: "NEVER_SUBSCRIBED",
		State:                   "DISABLED",
		Tags:                    []string{"B2B", "Login with Shop", "Shop"},
		UpdatedAt:               time.Date(2025, time.June, 4, 23, 17, 3, 0, time.UTC),
		ValidEmailAddress:       true,
		VerifiedEmail:           true,
	}
	juan := shopify.Customer{
		CreatedAt:               time.Date(2022, time.September, 12, 17, 26, 50, 0, time.UTC),
		DisplayName:             "juan zuniga",
		Email:                   "jzuniga@todoslosproductosmx.com",
		FirstName:               "juan",
		Id:                      "gid://shopify/Customer/6284129960104",
		LastName:                "zuniga",
		LegacyResourceId:        "6284129960104",
		LifetimeDuration:        "over 2 years",
		Locale:                  "en-MX",
		Note:                    "",
		NumberOfOrders:          "20",
		Phone:                   "",
		ProductSubscriberStatus: "NEVER_SUBSCRIBED",
		State:                   "ENABLED",
		Tags:                    []string{"reseller"},
		UpdatedAt:               time.Date(2025, time.May, 10, 0, 25, 17, 0, time.UTC),
		ValidEmailAddress:       true,
		VerifiedEmail:           true,
	}

	joseContact := shopify.Contact{
		CreatedAt:        time.Date(2025, time.May, 29, 15, 24, 41, 0, time.UTC),
		Id:               "gid://shopify/CompanyContact/895418536",
		IsMainContact:    true,
		LifetimeDuration: "13 days",
		Locale:           "en",
		Title:            "",
		UpdatedAt:        time.Date(2025, time.May, 29, 15, 24, 41, 0, time.UTC),
		Customer:         jose,
	}
	juanContact := shopify.Contact{
		CreatedAt:        time.Date(2025, time.May, 30, 14, 43, 33, 0, time.UTC),
		Id:               "gid://shopify/CompanyContact/897155240",
		IsMainContact:    true,
		LifetimeDuration: "12 days",
		Locale:           "en-MX",
		Title:            "",
		UpdatedAt:        time.Date(2025, time.May, 30, 14, 43, 33, 0, time.UTC),
		Customer:         juan,
	}
	tests := []struct {
		name string
		data string
		want shopify.SyncMetadata
	}{
		{
			name: "cryplex",
			data: `{"company_name":"Creativo3D","company_id":"gid://shopify/Company/2185199784","reseller":true,"vendor":true,"primary_company_email":"info@creativo3d.com","company":{"results":{"contactCount":1,"createdAt":"2025-05-29T15:24:39Z","customerSince":"2025-05-29T15:24:39Z","externalId":"A0F315E9C48C","hasTimelineComment":false,"id":"gid://shopify/Company/2185199784","lifetimeDuration":"13 days","name":"Creativo3D","note":"BF9B60168A88","updatedAt":"2025-06-09T23:27:10Z","locationsCount":{"count":1},"locations":{"nodes":[{"createdAt":"2025-05-29T15:24:39Z","currency":"GTQ","defaultCursor":"eyJkaXJlY3Rpb24iOiJuZXh0IiwibGFzdF9pZCI6Mjk1ODEzMTM2OCwibGFzdF92YWx1ZSI6IjI5NTgxMzEzNjgiLCJsaW1pdCI6MSwic29ydF9kaXIiOiJhc2MiLCJzb3J0X2ZpZWxkcyI6ImlkIn0=","externalId":null,"hasTimelineComment":false,"id":"gid://shopify/CompanyLocation/2958131368","locale":"en","name":"24 Avenida 13-20","note":null,"orderCount":1,"phone":null,"taxExemptions":[],"taxRegistrationId":"1209827-2","updatedAt":"2025-05-29T15:24:41Z","billingAddress":{"address1":"24 Avenida 13-20","address2":null,"city":"Ciudad de Guatemala","companyName":"Creativo3D","country":"Guatemala","countryCode":"GT","createdAt":"2025-05-29T15:24:40Z","firstName":"José","formattedAddress":["Creativo3D","24 Avenida 13-20","Ciudad de Guatemala GUA","01007","Guatemala"],"formattedArea":"Ciudad de Guatemala GU, Guatemala","id":"gid://shopify/CompanyAddress/2581856424","lastName":"Letona","phone":"+50237614657","province":"Guatemala","recipient":"Creativo3D","updatedAt":"2025-05-29T15:24:40Z","zip":"01007","zoneCode":"GUA"}}]},"contacts":{"nodes":[{"createdAt":"2025-05-29T15:24:41Z","id":"gid://shopify/CompanyContact/895418536","isMainContact":true,"lifetimeDuration":"13 days","locale":"en","title":null,"updatedAt":"2025-05-29T15:24:41Z","customer":{"canDelete":false,"createdAt":"2025-05-29T15:18:10Z","dataSaleOptOut":false,"displayName":"José Letona","email":"info@creativo3d.com","firstName":"José","hasTimelineComment":false,"id":"gid://shopify/Customer/8450627043496","lastName":"Letona","legacyResourceId":"8450627043496","lifetimeDuration":"13 days","locale":"en-GT","multipassIdentifier":null,"note":"","numberOfOrders":"1","phone":"+50237614657","productSubscriberStatus":"NEVER_SUBSCRIBED","state":"DISABLED","tags":["B2B","Login with Shop","Shop"],"updatedAt":"2025-06-04T23:17:03Z","validEmailAddress":true,"verifiedEmail":true}}]},"metafields":{"nodes":[{"key":"primary_company_email","type":"single_line_text_field","updatedAt":"2025-06-09T23:26:37Z","value":"info@creativo3d.com"},{"key":"vendor","type":"boolean","updatedAt":"2025-06-09T23:26:37Z","value":"true"},{"key":"reseller","type":"boolean","updatedAt":"2025-06-09T23:26:37Z","value":"true"},{"key":"domains","type":"single_line_text_field","updatedAt":"2025-06-09T23:26:37Z","value":"creativo3d.com"},{"key":"vat_number","type":"single_line_text_field","updatedAt":"2025-06-09T23:26:37Z","value":"1209827-2"},{"key":"tax_setting","type":"single_line_text_field","updatedAt":"2025-05-29T15:24:40Z","value":"Don't collect tax"},{"key":"ula_document","type":"file_reference","updatedAt":"2025-05-29T15:24:40Z","value":"gid://shopify/GenericFile/34844992241832"},{"key":"last_signed_reseller_agreement_dt","type":"date","updatedAt":"2025-05-29T15:24:40Z","value":"2025-05-29"}]},"metafield":null},"name":"Creativo3D","external_id":"gid://shopify/Company/2185199784","group_id":20488928016667,"do_not_sync_to_zendesk":false,"tags":["vendor"],"vendor":true,"reseller":true,"domain_names":["creativo3d.com"],"organization_fields":{"sync_shopify_company":false},"contacts":[{"createdAt":"2025-05-29T15:24:41Z","id":"gid://shopify/CompanyContact/895418536","isMainContact":true,"lifetimeDuration":"13 days","locale":"en","title":null,"updatedAt":"2025-05-29T15:24:41Z","customer":{"canDelete":false,"createdAt":"2025-05-29T15:18:10Z","dataSaleOptOut":false,"displayName":"José Letona","email":"info@creativo3d.com","firstName":"José","hasTimelineComment":false,"id":"gid://shopify/Customer/8450627043496","lastName":"Letona","legacyResourceId":"8450627043496","lifetimeDuration":"13 days","locale":"en-GT","multipassIdentifier":null,"note":"","numberOfOrders":"1","phone":"+50237614657","productSubscriberStatus":"NEVER_SUBSCRIBED","state":"DISABLED","tags":["B2B","Login with Shop","Shop"],"updatedAt":"2025-06-04T23:17:03Z","validEmailAddress":true,"verifiedEmail":true}}],"primary_company_email":"info@creativo3d.com","vat_number":"1209827-2","tax_setting":"Don't collect tax","ula_document":"gid://shopify/GenericFile/34844992241832","last_signed_reseller_agreement_dt":"2025-05-29","shared_tickets":true,"checksum":361772},"contact":{"createdAt":"2025-05-29T15:24:41Z","id":"gid://shopify/CompanyContact/895418536","isMainContact":true,"lifetimeDuration":"13 days","locale":"en","title":null,"updatedAt":"2025-05-29T15:24:41Z","customer":{"canDelete":false,"createdAt":"2025-05-29T15:18:10Z","dataSaleOptOut":false,"displayName":"José Letona","email":"info@creativo3d.com","firstName":"José","hasTimelineComment":false,"id":"gid://shopify/Customer/8450627043496","lastName":"Letona","legacyResourceId":"8450627043496","lifetimeDuration":"13 days","locale":"en-GT","multipassIdentifier":null,"note":"","numberOfOrders":"1","phone":"+50237614657","productSubscriberStatus":"NEVER_SUBSCRIBED","state":"DISABLED","tags":["B2B","Login with Shop","Shop"],"updatedAt":"2025-06-04T23:17:03Z","validEmailAddress":true,"verifiedEmail":true}},"contact_id":"gid://shopify/CompanyContact/895418536","customer":{"canDelete":false,"createdAt":"2025-05-29T15:18:10Z","dataSaleOptOut":false,"displayName":"José Letona","email":"info@creativo3d.com","firstName":"José","hasTimelineComment":false,"id":"gid://shopify/Customer/8450627043496","lastName":"Letona","legacyResourceId":"8450627043496","lifetimeDuration":"13 days","locale":"en-GT","multipassIdentifier":null,"note":"","numberOfOrders":"1","phone":"+50237614657","productSubscriberStatus":"NEVER_SUBSCRIBED","state":"DISABLED","tags":["B2B","Login with Shop","Shop"],"updatedAt":"2025-06-04T23:17:03Z","validEmailAddress":true,"verifiedEmail":true},"customer_id":"gid://shopify/Customer/8450627043496","checksum":501036}`,
			want: shopify.SyncMetadata{
				CompanyName:         "Creativo3D",
				CompanyId:           "gid://shopify/Company/2185199784",
				Reseller:            true,
				Vendor:              true,
				PrimaryCompanyEmail: "info@creativo3d.com",
				Company: shopify.DetailCompany{
					Result: shopify.Company{
						ContactCount:     1,
						CreatedAt:        time.Date(2025, time.May, 29, 15, 24, 39, 0, time.UTC),
						CustomerSince:    time.Date(2025, time.May, 29, 15, 24, 39, 0, time.UTC),
						ExternalId:       "A0F315E9C48C",
						Id:               "gid://shopify/Company/2185199784",
						LifetimeDuration: "13 days",
						Name:             "Creativo3D",
						Note:             "BF9B60168A88",
						UpdatedAt:        time.Date(2025, time.June, 9, 23, 27, 10, 0, time.UTC),
						LocationsCount:   shopify.Count{Count: 1},
						Locations: shopify.Locations{
							Nodes: []shopify.CompanyLocation{
								{
									CreatedAt:         time.Date(2025, time.May, 29, 15, 24, 39, 0, time.UTC),
									Currency:          "GTQ",
									DefaultCursor:     "eyJkaXJlY3Rpb24iOiJuZXh0IiwibGFzdF9pZCI6Mjk1ODEzMTM2OCwibGFzdF92YWx1ZSI6IjI5NTgxMzEzNjgiLCJsaW1pdCI6MSwic29ydF9kaXIiOiJhc2MiLCJzb3J0X2ZpZWxkcyI6ImlkIn0=",
									Id:                "gid://shopify/CompanyLocation/2958131368",
									Locale:            "en",
									Name:              "24 Avenida 13-20",
									OrderCount:        1,
									TaxExemptions:     []string{},
									TaxRegistrationId: "1209827-2",
									UpdatedAt:         time.Date(2025, time.May, 29, 15, 24, 41, 0, time.UTC),
									BillingAddress: shopify.Address{
										Address1:         "24 Avenida 13-20",
										City:             "Ciudad de Guatemala",
										CompanyName:      "Creativo3D",
										Country:          "Guatemala",
										CountryCode:      "GT",
										CreatedAt:        time.Date(2025, time.May, 29, 15, 24, 40, 0, time.UTC),
										FirstName:        "José",
										FormattedAddress: []string{"Creativo3D", "24 Avenida 13-20", "Ciudad de Guatemala GUA", "01007", "Guatemala"},
										FormattedArea:    "Ciudad de Guatemala GU, Guatemala",
										Id:               "gid://shopify/CompanyAddress/2581856424",
										LastName:         "Letona",
										Phone:            "+50237614657",
										Province:         "Guatemala",
										Recipient:        "Creativo3D",
										UpdatedAt:        time.Date(2025, time.May, 29, 15, 24, 40, 0, time.UTC),
										Zip:              "01007",
										ZoneCode:         "GUA",
									},
								},
							},
						},
						Contacts: shopify.Contacts{
							Nodes: []shopify.Contact{joseContact},
						},
						Metafields: shopify.MetafieldConnection{
							Nodes: []shopify.Metafield{
								{
									Key:       "primary_company_email",
									Type:      "single_line_text_field",
									UpdatedAt: time.Date(2025, time.June, 9, 23, 26, 37, 0, time.UTC),
									Value:     "info@creativo3d.com",
								},
								{
									Key:       "vendor",
									Type:      "boolean",
									UpdatedAt: time.Date(2025, time.June, 9, 23, 26, 37, 0, time.UTC),
									Value:     "true",
								},
								{
									Key:       "reseller",
									Type:      "boolean",
									UpdatedAt: time.Date(2025, time.June, 9, 23, 26, 37, 0, time.UTC),
									Value:     "true",
								},
								{
									Key:       "domains",
									Type:      "single_line_text_field",
									UpdatedAt: time.Date(2025, time.June, 9, 23, 26, 37, 0, time.UTC),
									Value:     "creativo3d.com",
								},
								{
									Key:       "vat_number",
									Type:      "single_line_text_field",
									UpdatedAt: time.Date(2025, time.June, 9, 23, 26, 37, 0, time.UTC),
									Value:     "1209827-2",
								},
								{
									Key:       "tax_setting",
									Type:      "single_line_text_field",
									UpdatedAt: time.Date(2025, time.May, 29, 15, 24, 40, 0, time.UTC),
									Value:     "Don't collect tax",
								},
								{
									Key:       "ula_document",
									Type:      "file_reference",
									UpdatedAt: time.Date(2025, time.May, 29, 15, 24, 40, 0, time.UTC),
									Value:     "gid://shopify/GenericFile/34844992241832",
								},
								{
									Key:       "last_signed_reseller_agreement_dt",
									Type:      "date",
									UpdatedAt: time.Date(2025, time.May, 29, 15, 24, 40, 0, time.UTC),
									Value:     "2025-05-29",
								},
							},
						},
					},
					Name:                        "Creativo3D",
					ExternalId:                  "gid://shopify/Company/2185199784",
					GroupId:                     20488928016667,
					Tags:                        []string{"vendor"},
					Reseller:                    true,
					Vendor:                      true,
					DomainNames:                 []string{"creativo3d.com"},
					OrganizationFields:          map[string]any{"sync_shopify_company": false},
					PrimaryCompanyEmail:         "info@creativo3d.com",
					VatNumber:                   "1209827-2",
					TaxSettings:                 "Don't collect tax",
					UlaDocument:                 "gid://shopify/GenericFile/34844992241832",
					LastSignedResellerAgreement: "2025-05-29",
					SharedTickets:               true,
					Checksum:                    361772,
					Contacts:                    []shopify.Contact{joseContact},
				},
				Contact:    joseContact,
				ContactId:  "gid://shopify/CompanyContact/895418536",
				Customer:   jose,
				CustomerId: "gid://shopify/Customer/8450627043496",
				Checksum:   501036,
			},
		},
		{
			name: "zendesk",
			data: `{"company_name":"Isabel Castelan Romo (TodosLosProductosMX)","company_id":"gid://shopify/Company/2187165864","reseller":true,"vendor":true,"primary_company_email":"mkt@todoslosproductosmx.com","company":{"results":{"contactCount":1,"createdAt":"2025-05-30T14:43:32Z","customerSince":"2025-05-30T14:43:32Z","externalId":"D50C618A2FC5","hasTimelineComment":false,"id":"gid://shopify/Company/2187165864","lifetimeDuration":"12 days","name":"Isabel Castelan Romo (TodosLosProductosMX)","note":null,"updatedAt":"2025-06-09T23:25:12Z","locationsCount":{"count":1},"locations":{"nodes":[{"createdAt":"2025-05-30T14:43:32Z","currency":"USD","defaultCursor":"eyJkaXJlY3Rpb24iOiJuZXh0IiwibGFzdF9pZCI6Mjk2MDUyMzQzMiwibGFzdF92YWx1ZSI6IjI5NjA1MjM0MzIiLCJsaW1pdCI6MSwic29ydF9kaXIiOiJhc2MiLCJzb3J0X2ZpZWxkcyI6ImlkIn0=","externalId":null,"hasTimelineComment":false,"id":"gid://shopify/CompanyLocation/2960523432","locale":"en","name":"José María Rico 230","note":null,"orderCount":20,"phone":null,"taxExemptions":[],"taxRegistrationId":"CARI790303MR5","updatedAt":"2025-05-30T14:43:33Z","billingAddress":{"address1":"José María Rico 230","address2":null,"city":"Ciudad de México","companyName":"Isabel Castelan Romo (TodosLosProductosMX)","country":"Mexico","countryCode":"MX","createdAt":"2025-05-30T14:43:32Z","firstName":"Isabel Castelan Romo","formattedAddress":["Isabel Castelan Romo (TodosLosProductosMX)","José María Rico 230","03100 Ciudad de México DF","Mexico"],"formattedArea":"Ciudad de México CMX, Mexico","id":"gid://shopify/CompanyAddress/2583232680","lastName":"Romo","phone":"+525556888039","province":"Ciudad de México","recipient":"TodosLosProductosMX","updatedAt":"2025-05-30T14:43:32Z","zip":"03100","zoneCode":"DF"}}]},"contacts":{"nodes":[{"createdAt":"2025-05-30T14:43:33Z","id":"gid://shopify/CompanyContact/897155240","isMainContact":true,"lifetimeDuration":"12 days","locale":"en-MX","title":null,"updatedAt":"2025-05-30T14:43:33Z","customer":{"canDelete":false,"createdAt":"2022-09-12T17:26:50Z","dataSaleOptOut":false,"displayName":"juan zuniga","email":"jzuniga@todoslosproductosmx.com","firstName":"juan","hasTimelineComment":false,"id":"gid://shopify/Customer/6284129960104","lastName":"zuniga","legacyResourceId":"6284129960104","lifetimeDuration":"over 2 years","locale":"en-MX","multipassIdentifier":null,"note":"","numberOfOrders":"20","phone":null,"productSubscriberStatus":"NEVER_SUBSCRIBED","state":"ENABLED","tags":["reseller"],"updatedAt":"2025-05-10T00:25:17Z","validEmailAddress":true,"verifiedEmail":true}}]},"metafields":{"nodes":[{"key":"primary_company_email","type":"single_line_text_field","updatedAt":"2025-06-09T23:15:29Z","value":"mkt@todoslosproductosmx.com"},{"key":"vendor","type":"boolean","updatedAt":"2025-06-09T23:15:29Z","value":"true"},{"key":"reseller","type":"boolean","updatedAt":"2025-06-09T23:15:29Z","value":"true"},{"key":"domains","type":"single_line_text_field","updatedAt":"2025-06-09T23:15:29Z","value":"www.todoslosproductosmx.com"},{"key":"vat_number","type":"single_line_text_field","updatedAt":"2025-06-09T23:15:29Z","value":"CARI790303MR5"},{"key":"ula_document","type":"file_reference","updatedAt":"2025-05-30T14:43:32Z","value":"gid://shopify/GenericFile/34861424836776"},{"key":"last_signed_reseller_agreement_dt","type":"date","updatedAt":"2025-05-30T14:43:32Z","value":"2025-05-30"}]},"metafield":null},"name":"Isabel Castelan Romo (TodosLosProductosMX)","external_id":"gid://shopify/Company/2187165864","group_id":20488928016667,"do_not_sync_to_zendesk":false,"tags":["vendor"],"vendor":true,"reseller":true,"domain_names":["www.todoslosproductosmx.com"],"organization_fields":{"sync_shopify_company":false},"contacts":[{"createdAt":"2025-05-30T14:43:33Z","id":"gid://shopify/CompanyContact/897155240","isMainContact":true,"lifetimeDuration":"12 days","locale":"en-MX","title":null,"updatedAt":"2025-05-30T14:43:33Z","customer":{"canDelete":false,"createdAt":"2022-09-12T17:26:50Z","dataSaleOptOut":false,"displayName":"juan zuniga","email":"jzuniga@todoslosproductosmx.com","firstName":"juan","hasTimelineComment":false,"id":"gid://shopify/Customer/6284129960104","lastName":"zuniga","legacyResourceId":"6284129960104","lifetimeDuration":"over 2 years","locale":"en-MX","multipassIdentifier":null,"note":"","numberOfOrders":"20","phone":null,"productSubscriberStatus":"NEVER_SUBSCRIBED","state":"ENABLED","tags":["reseller"],"updatedAt":"2025-05-10T00:25:17Z","validEmailAddress":true,"verifiedEmail":true}}],"primary_company_email":"mkt@todoslosproductosmx.com","vat_number":"CARI790303MR5","ula_document":"gid://shopify/GenericFile/34861424836776","last_signed_reseller_agreement_dt":"2025-05-30","shared_tickets":true,"checksum":370920},"contact":{"createdAt":"2025-05-30T14:43:33Z","id":"gid://shopify/CompanyContact/897155240","isMainContact":true,"lifetimeDuration":"12 days","locale":"en-MX","title":null,"updatedAt":"2025-05-30T14:43:33Z","customer":{"canDelete":false,"createdAt":"2022-09-12T17:26:50Z","dataSaleOptOut":false,"displayName":"juan zuniga","email":"jzuniga@todoslosproductosmx.com","firstName":"juan","hasTimelineComment":false,"id":"gid://shopify/Customer/6284129960104","lastName":"zuniga","legacyResourceId":"6284129960104","lifetimeDuration":"over 2 years","locale":"en-MX","multipassIdentifier":null,"note":"","numberOfOrders":"20","phone":null,"productSubscriberStatus":"NEVER_SUBSCRIBED","state":"ENABLED","tags":["reseller"],"updatedAt":"2025-05-10T00:25:17Z","validEmailAddress":true,"verifiedEmail":true}},"contact_id":"gid://shopify/CompanyContact/897155240","customer":{"canDelete":false,"createdAt":"2022-09-12T17:26:50Z","dataSaleOptOut":false,"displayName":"juan zuniga","email":"jzuniga@todoslosproductosmx.com","firstName":"juan","hasTimelineComment":false,"id":"gid://shopify/Customer/6284129960104","lastName":"zuniga","legacyResourceId":"6284129960104","lifetimeDuration":"over 2 years","locale":"en-MX","multipassIdentifier":null,"note":"","numberOfOrders":"20","phone":null,"productSubscriberStatus":"NEVER_SUBSCRIBED","state":"ENABLED","tags":["reseller"],"updatedAt":"2025-05-10T00:25:17Z","validEmailAddress":true,"verifiedEmail":true},"customer_id":"gid://shopify/Customer/6284129960104","checksum":514737}`,
			want: shopify.SyncMetadata{
				CompanyName:         "Isabel Castelan Romo (TodosLosProductosMX)",
				CompanyId:           "gid://shopify/Company/2187165864",
				Reseller:            true,
				Vendor:              true,
				PrimaryCompanyEmail: "mkt@todoslosproductosmx.com",
				Company: shopify.DetailCompany{
					Result: shopify.Company{
						ContactCount:     1,
						CreatedAt:        time.Date(2025, time.May, 30, 14, 43, 32, 0, time.UTC),
						CustomerSince:    time.Date(2025, time.May, 30, 14, 43, 32, 0, time.UTC),
						ExternalId:       "D50C618A2FC5",
						Id:               "gid://shopify/Company/2187165864",
						LifetimeDuration: "12 days",
						Name:             "Isabel Castelan Romo (TodosLosProductosMX)",
						Note:             "",
						UpdatedAt:        time.Date(2025, time.June, 9, 23, 25, 12, 0, time.UTC),
						LocationsCount:   shopify.Count{Count: 1},
						Locations: shopify.Locations{
							Nodes: []shopify.CompanyLocation{
								{
									CreatedAt:          time.Date(2025, time.May, 30, 14, 43, 32, 0, time.UTC),
									Currency:           "USD",
									DefaultCursor:      "eyJkaXJlY3Rpb24iOiJuZXh0IiwibGFzdF9pZCI6Mjk2MDUyMzQzMiwibGFzdF92YWx1ZSI6IjI5NjA1MjM0MzIiLCJsaW1pdCI6MSwic29ydF9kaXIiOiJhc2MiLCJzb3J0X2ZpZWxkcyI6ImlkIn0=",
									ExternalId:         "",
									HasTimelineComment: false,
									Id:                 "gid://shopify/CompanyLocation/2960523432",
									Locale:             "en",
									Name:               "José María Rico 230",
									Note:               "",
									OrderCount:         20,
									Phone:              "",
									TaxExemptions:      []string{},
									TaxRegistrationId:  "CARI790303MR5",
									UpdatedAt:          time.Date(2025, time.May, 30, 14, 43, 33, 0, time.UTC),
									BillingAddress: shopify.Address{
										Address1:    "José María Rico 230",
										Address2:    "",
										City:        "Ciudad de México",
										CompanyName: "Isabel Castelan Romo (TodosLosProductosMX)",
										Country:     "Mexico",
										CountryCode: "MX",
										CreatedAt:   time.Date(2025, time.May, 30, 14, 43, 32, 0, time.UTC),
										FirstName:   "Isabel Castelan Romo",
										FormattedAddress: []string{
											"Isabel Castelan Romo (TodosLosProductosMX)",
											"José María Rico 230",
											"03100 Ciudad de México DF",
											"Mexico",
										},
										FormattedArea: "Ciudad de México CMX, Mexico",
										Id:            "gid://shopify/CompanyAddress/2583232680",
										LastName:      "Romo",
										Phone:         "+525556888039",
										Province:      "Ciudad de México",
										Recipient:     "TodosLosProductosMX",
										UpdatedAt:     time.Date(2025, time.May, 30, 14, 43, 32, 0, time.UTC),
										Zip:           "03100",
										ZoneCode:      "DF",
									},
								},
							},
						},
						Contacts: shopify.Contacts{
							Nodes: []shopify.Contact{juanContact},
						},
						Metafields: shopify.MetafieldConnection{
							Nodes: []shopify.Metafield{
								{
									Key:       "primary_company_email",
									Type:      "single_line_text_field",
									UpdatedAt: time.Date(2025, time.June, 9, 23, 15, 29, 0, time.UTC),
									Value:     "mkt@todoslosproductosmx.com",
								},
								{
									Key:       "vendor",
									Type:      "boolean",
									UpdatedAt: time.Date(2025, time.June, 9, 23, 15, 29, 0, time.UTC),
									Value:     "true",
								},
								{
									Key:       "reseller",
									Type:      "boolean",
									UpdatedAt: time.Date(2025, time.June, 9, 23, 15, 29, 0, time.UTC),
									Value:     "true",
								},
								{
									Key:       "domains",
									Type:      "single_line_text_field",
									UpdatedAt: time.Date(2025, time.June, 9, 23, 15, 29, 0, time.UTC),
									Value:     "www.todoslosproductosmx.com",
								},
								{
									Key:       "vat_number",
									Type:      "single_line_text_field",
									UpdatedAt: time.Date(2025, time.June, 9, 23, 15, 29, 0, time.UTC),
									Value:     "CARI790303MR5",
								},
								{
									Key:       "ula_document",
									Type:      "file_reference",
									UpdatedAt: time.Date(2025, time.May, 30, 14, 43, 32, 0, time.UTC),
									Value:     "gid://shopify/GenericFile/34861424836776",
								},
								{
									Key:       "last_signed_reseller_agreement_dt",
									Type:      "date",
									UpdatedAt: time.Date(2025, time.May, 30, 14, 43, 32, 0, time.UTC),
									Value:     "2025-05-30",
								},
							},
						},
					},
					Name:                        "Isabel Castelan Romo (TodosLosProductosMX)",
					ExternalId:                  "gid://shopify/Company/2187165864",
					GroupId:                     20488928016667,
					Tags:                        []string{"vendor"},
					Reseller:                    true,
					Vendor:                      true,
					DomainNames:                 []string{"www.todoslosproductosmx.com"},
					OrganizationFields:          map[string]any{"sync_shopify_company": false},
					PrimaryCompanyEmail:         "mkt@todoslosproductosmx.com",
					VatNumber:                   "CARI790303MR5",
					TaxSettings:                 "",
					UlaDocument:                 "gid://shopify/GenericFile/34861424836776",
					LastSignedResellerAgreement: "2025-05-30",
					SharedTickets:               true,
					Checksum:                    370920,
					Contacts:                    []shopify.Contact{juanContact},
				},
				Contact:    juanContact,
				ContactId:  "gid://shopify/CompanyContact/897155240",
				Customer:   juan,
				CustomerId: "gid://shopify/Customer/6284129960104",
				Checksum:   514737,
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			got := shopify.SyncMetadata{}
			err := json.Unmarshal([]byte(test.data), &got)
			if err != nil {
				t.Errorf("%s: TestGetRecords: error %s", test.name, err)
			}
			if df := cmp.Diff(test.want, got); df != "" {
				t.Errorf("%s: TestGetRecords: -want +got\n%s", test.name, df)
			}
		})
	}
}

/*
{
    "company_name": "Creativo3D",
    "company_id": "gid://shopify/Company/2185199784",
    "reseller": true,
    "vendor": true,
    "primary_company_email": "info@creativo3d.com",
    "company": {
        "results": {
            "contactCount": 1,
            "createdAt": "2025-05-29T15:24:39Z",
            "customerSince": "2025-05-29T15:24:39Z",
            "externalId": "A0F315E9C48C",
            "hasTimelineComment": false,
            "id": "gid://shopify/Company/2185199784",
            "lifetimeDuration": "13 days",
            "name": "Creativo3D",
            "note": "BF9B60168A88",
            "updatedAt": "2025-06-09T23:27:10Z",
            "locationsCount": {
                "count": 1
            },
            "locations": {
                "nodes": [
                    {
                    "createdAt": "2025-05-29T15:24:39Z",
                    "currency": "GTQ",
                    "defaultCursor": "eyJkaXJlY3Rpb24iOiJuZXh0IiwibGFzdF9pZCI6Mjk1ODEzMTM2OCwibGFzdF92YWx1ZSI6IjI5NTgxMzEzNjgiLCJsaW1pdCI6MSwic29ydF9kaXIiOiJhc2MiLCJzb3J0X2ZpZWxkcyI6ImlkIn0=",
                    "externalId": null,
                    "hasTimelineComment": false,
                    "id": "gid://shopify/CompanyLocation/2958131368",
                    "locale": "en",
                    "name": "24 Avenida 13-20",
                    "note": null,
                    "orderCount": 1,
                    "phone": null,
                    "taxExemptions": [],
                    "taxRegistrationId": "1209827-2",
                    "updatedAt": "2025-05-29T15:24:41Z",
                    "billingAddress": {
                        "address1": "24 Avenida 13-20",
                        "address2": null,
                        "city": "Ciudad de Guatemala",
                        "companyName": "Creativo3D",
                        "country": "Guatemala",
                        "countryCode": "GT",
                        "createdAt": "2025-05-29T15:24:40Z",
                        "firstName": "José",
                        "formattedAddress": [
                        "Creativo3D",
                        "24 Avenida 13-20",
                        "Ciudad de Guatemala GUA",
                        "01007",
                        "Guatemala"
                        ],
                        "formattedArea": "Ciudad de Guatemala GU, Guatemala",
                        "id": "gid://shopify/CompanyAddress/2581856424",
                        "lastName": "Letona",
                        "phone": "+50237614657",
                        "province": "Guatemala",
                        "recipient": "Creativo3D",
                        "updatedAt": "2025-05-29T15:24:40Z",
                        "zip": "01007",
                        "zoneCode": "GUA"
                    }
                    }
                ]
            },
            "contacts": {
                "nodes": [
                    {
                    "createdAt": "2025-05-29T15:24:41Z",
                    "id": "gid://shopify/CompanyContact/895418536",
                    "isMainContact": true,
                    "lifetimeDuration": "13 days",
                    "locale": "en",
                    "title": null,
                    "updatedAt": "2025-05-29T15:24:41Z",
                    "customer": {
                        "canDelete": false,
                        "createdAt": "2025-05-29T15:18:10Z",
                        "dataSaleOptOut": false,
                        "displayName": "José Letona",
                        "email": "info@creativo3d.com",
                        "firstName": "José",
                        "hasTimelineComment": false,
                        "id": "gid://shopify/Customer/8450627043496",
                        "lastName": "Letona",
                        "legacyResourceId": "8450627043496",
                        "lifetimeDuration": "13 days",
                        "locale": "en-GT",
                        "multipassIdentifier": null,
                        "note": "",
                        "numberOfOrders": "1",
                        "phone": "+50237614657",
                        "productSubscriberStatus": "NEVER_SUBSCRIBED",
                        "state": "DISABLED",
                        "tags": [
                        "B2B",
                        "Login with Shop",
                        "Shop"
                        ],
                        "updatedAt": "2025-06-04T23:17:03Z",
                        "validEmailAddress": true,
                        "verifiedEmail": true
                    }
                    }
                ]
            },
            "metafields": {
                "nodes": [
                    {
                        "key": "primary_company_email",
                        "type": "single_line_text_field",
                        "updatedAt": "2025-06-09T23:26:37Z",
                        "value": "info@creativo3d.com"
                    },
                    {
                        "key": "vendor",
                        "type": "boolean",
                        "updatedAt": "2025-06-09T23:26:37Z",
                        "value": "true"
                    },
                    {
                        "key": "reseller",
                        "type": "boolean",
                        "updatedAt": "2025-06-09T23:26:37Z",
                        "value": "true"
                    },
                    {
                        "key": "domains",
                        "type": "single_line_text_field",
                        "updatedAt": "2025-06-09T23:26:37Z",
                        "value": "creativo3d.com"
                    },
                    {
                        "key": "vat_number",
                        "type": "single_line_text_field",
                        "updatedAt": "2025-06-09T23:26:37Z",
                        "value": "1209827-2"
                    },
                    {
                        "key": "tax_setting",
                        "type": "single_line_text_field",
                        "updatedAt": "2025-05-29T15:24:40Z",
                        "value": "Don't collect tax"
                    },
                    {
                        "key": "ula_document",
                        "type": "file_reference",
                        "updatedAt": "2025-05-29T15:24:40Z",
                        "value": "gid://shopify/GenericFile/34844992241832"
                    },
                    {
                        "key": "last_signed_reseller_agreement_dt",
                        "type": "date",
                        "updatedAt": "2025-05-29T15:24:40Z",
                        "value": "2025-05-29"
                    }
                ]
            },
            "metafield": null
        },
        "name": "Creativo3D",
        "external_id": "gid://shopify/Company/2185199784",
        "group_id": 20488928016667,
        "do_not_sync_to_zendesk": false,
        "tags": [
            "vendor"
        ],
        "vendor": true,
        "reseller": true,
        "domain_names": [
            "creativo3d.com"
        ],
        "organization_fields": {
            "sync_shopify_company": false
        },
        "contacts": [
            {
                "createdAt": "2025-05-29T15:24:41Z",
                "id": "gid://shopify/CompanyContact/895418536",
                "isMainContact": true,
                "lifetimeDuration": "13 days",
                "locale": "en",
                "title": null,
                "updatedAt": "2025-05-29T15:24:41Z",
                "customer": {
                    "canDelete": false,
                    "createdAt": "2025-05-29T15:18:10Z",
                    "dataSaleOptOut": false,
                    "displayName": "José Letona",
                    "email": "info@creativo3d.com",
                    "firstName": "José",
                    "hasTimelineComment": false,
                    "id": "gid://shopify/Customer/8450627043496",
                    "lastName": "Letona",
                    "legacyResourceId": "8450627043496",
                    "lifetimeDuration": "13 days",
                    "locale": "en-GT",
                    "multipassIdentifier": null,
                    "note": "",
                    "numberOfOrders": "1",
                    "phone": "+50237614657",
                    "productSubscriberStatus": "NEVER_SUBSCRIBED",
                    "state": "DISABLED",
                    "tags": [
                        "B2B",
                        "Login with Shop",
                        "Shop"
                    ],
                    "updatedAt": "2025-06-04T23:17:03Z",
                    "validEmailAddress": true,
                    "verifiedEmail": true
                }
            }
        ],
        "primary_company_email": "info@creativo3d.com",
        "vat_number": "1209827-2",
        "tax_setting": "Don't collect tax",
        "ula_document": "gid://shopify/GenericFile/34844992241832",
        "last_signed_reseller_agreement_dt": "2025-05-29",
        "shared_tickets": true,
        "checksum": 361772
    },
    "contact": {
        "createdAt": "2025-05-29T15:24:41Z",
        "id": "gid://shopify/CompanyContact/895418536",
        "isMainContact": true,
        "lifetimeDuration": "13 days",
        "locale": "en",
        "title": null,
        "updatedAt": "2025-05-29T15:24:41Z",
        "customer": {
            "canDelete": false,
            "createdAt": "2025-05-29T15:18:10Z",
            "dataSaleOptOut": false,
            "displayName": "José Letona",
            "email": "info@creativo3d.com",
            "firstName": "José",
            "hasTimelineComment": false,
            "id": "gid://shopify/Customer/8450627043496",
            "lastName": "Letona",
            "legacyResourceId": "8450627043496",
            "lifetimeDuration": "13 days",
            "locale": "en-GT",
            "multipassIdentifier": null,
            "note": "",
            "numberOfOrders": "1",
            "phone": "+50237614657",
            "productSubscriberStatus": "NEVER_SUBSCRIBED",
            "state": "DISABLED",
            "tags": [
                "B2B",
                "Login with Shop",
                "Shop"
            ],
            "updatedAt": "2025-06-04T23:17:03Z",
            "validEmailAddress": true,
            "verifiedEmail": true
        }
    },
    "contact_id": "gid://shopify/CompanyContact/895418536",
    "customer": {
    "canDelete": false,
    "createdAt": "2025-05-29T15:18:10Z",
    "dataSaleOptOut": false,
    "displayName": "José Letona",
    "email": "info@creativo3d.com",
    "firstName": "José",
    "hasTimelineComment": false,
    "id": "gid://shopify/Customer/8450627043496",
    "lastName": "Letona",
    "legacyResourceId": "8450627043496",
    "lifetimeDuration": "13 days",
    "locale": "en-GT",
    "multipassIdentifier": null,
    "note": "",
    "numberOfOrders": "1",
    "phone": "+50237614657",
    "productSubscriberStatus": "NEVER_SUBSCRIBED",
    "state": "DISABLED",
    "tags": [
        "B2B",
        "Login with Shop",
        "Shop"
    ],
    "updatedAt": "2025-06-04T23:17:03Z",
    "validEmailAddress": true,
    "verifiedEmail": true
    },
    "customer_id": "gid://shopify/Customer/8450627043496",
    "checksum": 501036
}
*/
