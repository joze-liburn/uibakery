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

// DetailsToSync expands a given list of companies (in DetailCompany format) so
// that for every contact and every destination it creates a database record.
func DetailsToSync(companies []shopify.DetailCompany, destinations []string, submitter string, batch time.Time) <-chan lbqueue.LogSyncRecord {
	out := make(chan lbqueue.LogSyncRecord)
	base := lbqueue.LogSyncRecord{
		Submitter:        submitter,
		BatchId:          &batch,
		RecordType:       "customer",
		BulkSubmit:       true,
		SubmissionStatus: "NEW",
	}
	go func() {
		defer close(out)
		for _, company := range companies {
			template := shopify.SyncMetadata{
				CompanyName:         company.Name,
				CompanyId:           company.ExternalId,
				Reseller:            company.Reseller,
				Vendor:              company.Vendor,
				PrimaryCompanyEmail: company.PrimaryCompanyEmail,
				Company:             company,
			}
			base.CompanyId = company.ExternalId
			base.SyncMdChecksum = company.Checksum
			for _, contact := range company.Contacts {
				dbRow := template
				dbRow.ContactId = contact.Id
				dbRow.Contact = contact
				dbRow.Customer = contact.Customer
				dbRow.CustomerId = contact.Customer.Id
				//		elt.NoSyncZendesk = cnt.NoSyncZendesk

				json, jerr := json.Marshal(dbRow)
				if jerr != nil {
					return
				}
				next := base
				next.SyncMetadata = string(json)
				next.RecordId = contact.Customer.Id
				for _, dst := range destinations {
					next.DestinationName = dst
					next.LbCreate = time.Now()
					out <- next
				}
			}
		}
	}()
	return out
}
