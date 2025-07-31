package sync

import (
	json "encoding/json" // TBD: "encoding/json/v2" when go 1.25
	"time"

	"gitlab.com/joze-liburn/uibakery/lbqueue"
	"gitlab.com/joze-liburn/uibakery/shopify"
)

// GetSyncMetadata takes the record from the database and returs embedded sync
// data.
func GetSyncMetadata(from lbqueue.LogSyncRecord) (shopify.SyncMetadata, error) {
	ret := shopify.SyncMetadata{}
	err := json.Unmarshal([]byte(from.SyncMetadata), &ret)

	return ret, err
}

// For every contact...
func DetailToSync(company shopify.DetailCompany, destination string, submitter string, batch *time.Time) ([]lbqueue.LogSyncRecord, error) {
	ret := []lbqueue.LogSyncRecord{}
	base := lbqueue.LogSyncRecord{
		Submitter:       submitter,
		DestinationName: destination,
		BatchId:         batch,
		CompanyId:       company.ExternalId,
	}
	template := shopify.SyncMetadata{
		CompanyName:         company.Name,
		CompanyId:           company.ExternalId,
		Reseller:            company.Reseller,
		Vendor:              company.Vendor,
		PrimaryCompanyEmail: company.PrimaryCompanyEmail,
		Company:             company,
	}
	for _, cnt := range company.Contacts {
		elt := template
		elt.ContactId = cnt.Id
		elt.Contact = cnt
		elt.Customer = cnt.Customer
		elt.CustomerId = cnt.Customer.Id
		//		elt.NoSyncZendesk = cnt.NoSyncZendesk
		json, err := json.Marshal(elt)
		if err != nil {
			return []lbqueue.LogSyncRecord{}, err
		}
		next := base
		next.SyncMetadata = string(json)
		ret = append(ret, next)
	}
	return ret, nil
}
