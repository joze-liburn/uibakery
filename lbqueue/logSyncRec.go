package lbqueue

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
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

	// LbDb is a database handle.
	LbDb struct {
		opened bool
		user   string
		addr   string
		db     *sql.DB
	}

	// QueueCount hholds a count of records by destination, submission status
	// and claim status.
	QueueCount struct {
		DestinationName  string
		SubmissionStatus string
		Claimed          bool
		Count            uint
	}

	// Claim is a claim id and a count (of matching records).
	Claim struct {
		Id    string
		Count uint
	}
)

// GetTime converts database nullable timestamp null values into Go's nil.
func GetTime(tm sql.NullTime) *time.Time {
	if !tm.Valid {
		return nil
	}
	v, err := tm.Value()
	if err != nil {
		// TODO: do something. Yell.
		return nil
	}
	res := v.(time.Time)
	return &res
}

// GetTime converts database nullable string null values into Go's empty strings.
func GetString(s sql.NullString) string {
	if !s.Valid {
		return ""
	}
	v, err := s.Value()
	if err != nil {
		// TODO: do something. Yell.
		return ""
	}
	return v.(string)
}

// IsOpened verifies the database is opened.
func (db *LbDb) IsOpened() bool {
	return db.opened
}

// Open opens a database.
func (db *LbDb) Open(user string, secret string, host string, port uint, database string) error {
	var err error
	db.user = user
	db.addr = host
	db.db, err = sql.Open("pgx", fmt.Sprintf("user=%s password=%s host=%s port=%d database=%s", user, secret, host, port, database)) // "lb_ap_uibakery@35.196.117.104:5432/lightburn")
	if err != nil {
		return err
	}
	db.opened = true
	return nil
}

// ClaimRecords claims up to max unclaimed records from the database, and
// returns a claim id.
func (db *LbDb) ClaimRecords(max uint) (string, int, error) {
	if max == 0 {
		return "", 0, nil
	}
	guid, err := uuid.NewRandom()
	if err != nil {
		return "", 0, err
	}
	res, err := db.db.Exec(`with batch as (
    select
        lb_id
    from
        public.log_sync_record
    where
        (ps_guid is null or ps_guid = '' or ps_guid = '?')
        and (submission_status in ('new', 'pending')
             or (submission_status = 'grief' and lb_update_ts < now() - interval '15 minutes'))
        and destination_name = 'zendesk'
  	order by lb_id
    limit $2
)
update
    public.log_sync_record lsr
set
    ps_guid = $1
from
    batch
where
    lsr.lb_id = batch.lb_id;
`, guid.String(), max)
	if err != nil {
		return "", 0, fmt.Errorf("%w -- checked a secret?: ", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return "", 0, err
	}
	return guid.String(), int(rows), err
}

// UnclaimRecords claims up to max unclaimed records from the database, and
// returns a claim id.
func (db *LbDb) UnclaimRecords(claim string) (int, error) {
	res, err := db.db.Exec(`update
    public.log_sync_record lsr
set
    ps_guid = null
where
    ps_guid = $1
`, claim)
	if err != nil {
		return 0, fmt.Errorf("%w -- checked a secret?: ", err)
	}
	rows, err := res.RowsAffected()
	if err != nil {
		return 0, err
	}
	return int(rows), err
}

// GetClaimedRecords reads claimed records from the database. Parameter claim
// should be a value retured by ClaimRecords (but feeding some random values
// will simply return an empty set).
func (db *LbDb) GetClaimedRecords(claim string) ([]LogSyncRecord, error) {
	rows, err := db.db.Query(`select
    lb_id
    , submitter
    , batch_id
    , destination_name
    , record_type
    , record_id
    , company_id
    , submission_status
    , lb_create_ts
    , lb_update_ts
    , lb_complete_ts
    , sync_metadata
    , ps_guid
    , bulk_submit
    , sync_md_checksum
from
    public.log_sync_record
where
    ps_guid = $1`, claim)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []LogSyncRecord{}
	for rows.Next() {
		row := LogSyncRecord{}
		var (
			batch_id   sql.NullTime
			updated    sql.NullTime
			completed  sql.NullTime
			company_id sql.NullString
			ps_guid    sql.NullString
		)
		err := rows.Scan(&row.LbId,
			&row.Submitter,
			&batch_id,
			&row.DestinationName,
			&row.RecordType,
			&row.RecordId,
			&company_id,
			&row.SubmissionStatus,
			&row.LbCreate,
			&updated,
			&completed,
			&row.SyncMetadata,
			&ps_guid,
			&row.BulkSubmit,
			&row.SyncMdChecksum)
		if err != nil {
			return result, err
		}
		row.BatchId = GetTime(batch_id)
		row.LbUpdate = GetTime(updated)
		row.LbComplete = GetTime(completed)
		row.PsGuid = GetString(ps_guid)
		row.CompanyId = GetString(company_id)
		result = append(result, row)
	}
	return result, nil
}

// GetQueueStats obtains a count of records by status, destination, and claim
// status (not by individual claims).
func (db *LbDb) GetQueueStats() ([]QueueCount, error) {
	rows, err := db.db.Query(`select
    submission_status
    , destination_name
    , case when ps_guid is null or ps_guid = '' or ps_guid = '?' then 0 else 1 end claimed
    , count(1)
from
    log_sync_record
group by
    submission_status
    , destination_name
    , case when ps_guid is null or ps_guid = '' or ps_guid = '?' then 0 else 1 end`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []QueueCount{}
	for rows.Next() {
		row := QueueCount{}
		var (
			claimed int
			count   uint
		)
		err := rows.Scan(&row.SubmissionStatus,
			&row.DestinationName,
			&claimed,
			&count)
		if err != nil {
			return result, err
		}
		row.Count = count
		row.Claimed = claimed == 1
		result = append(result, row)
	}
	return result, nil
}

// GetQueueStats obtains a count of records by status, destination, and claim
// status (not by individual claims).
func (db *LbDb) ListClaims(status *string) ([]Claim, error) {
	rows, err := db.db.Query(`select
    ps_guid
    , count(1)
from
    log_sync_record
where
    ps_guid is not null and ps_guid <> '' and ps_guid <> '?'
    and submission_status is not distinct from coalesce($1, submission_status)
group by
    ps_guid`, status)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := []Claim{}
	for rows.Next() {
		var claim Claim
		err := rows.Scan(&claim.Id, &claim.Count)
		if err != nil {
			return result, err
		}
		result = append(result, claim)
	}
	return result, nil
}
