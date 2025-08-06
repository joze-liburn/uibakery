package lbqueue

import (
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type (
	// LogSyncRecord represents a row in a database and encapsuates
	// (far beyond minimum) parameters for a single object to a single sink
	// synchronization (source is always Shopify).
	LogSyncRecord struct {
		LbId             int32      `json:"lb_id,omitempty"`
		Submitter        string     `json:"submitter,omitempty"`
		BatchId          *time.Time `json:"batch_id,omitempty"`
		DestinationName  string     `json:"destination_name,omitempty"`
		RecordType       string     `json:"record_type,omitempty"`
		RecordId         string     `json:"record_id,omitempty"`
		CompanyId        string     `json:"company_id,omitempty"`
		SubmissionStatus string     `json:"submission_status,omitempty"`
		LbCreate         time.Time  `json:"lb_create_ts,omitempty"`
		LbUpdate         *time.Time `json:"lb_update_ts,omitempty"`
		LbComplete       *time.Time `json:"lb_complete_ts,omitempty"`
		SyncMetadata     string     `json:"sync_metadata,omitempty"`
		PsGuid           string     `json:"ps_guid,omitempty"`
		BulkSubmit       bool       `json:"bulk_submit,omitempty"`
		SyncMdChecksum   int64      `json:"sync_md_checksum,omitempty"`
	}
)
