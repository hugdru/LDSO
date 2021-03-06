package datastore

import (
	"database/sql"
	"errors"
	"gopkg.in/guregu/null.v3/zero"
	"server/datastore/generators"
	"server/datastore/metadata"
	"strconv"
	"time"
)

type Audit struct {
	Id           int64       `json:"id" db:"id"`
	IdProperty   int64       `json:"idProperty" db:"id_property"`
	IdAuditor    int64       `json:"idAuditor" db:"id_auditor"`
	IdTemplate   int64       `json:"idTemplate" db:"id_template"`
	Rating       zero.Int    `json:"rating" db:"rating"`
	Observation  zero.String `json:"observation" db:"observation"`
	CreatedDate  time.Time   `json:"createdDate" db:"created_date"`
	FinishedDate zero.Time   `json:"finishedDate" db:"finished_date"`

	meta metadata.Metadata
}

func (a *Audit) SetExists() {
	a.meta.Exists = true
}

func (a *Audit) SetDeleted() {
	a.meta.Deleted = true
}

func (a *Audit) Exists() bool {
	return a.meta.Exists
}

func (a *Audit) MustSet(idProperty, idAuditor, idTemplate int64) error {

	if idProperty != 0 {
		a.IdProperty = idProperty
	} else {
		return errors.New("idProperty must be set")
	}
	if idAuditor != 0 {
		a.IdAuditor = idAuditor
	} else {
		return errors.New("idAuditor must be set")
	}
	if idTemplate != 0 {
		a.IdTemplate = idTemplate
	} else {
		return errors.New("idTemplate must be set")
	}

	return nil
}

func (a *Audit) AllSetIfNotEmptyOrNil(
	idProperty, idAuditor, idTemplate int64,
	rating int64, observation string) error {

	if idProperty != 0 {
		a.IdProperty = idProperty
	}

	if idAuditor != 0 {
		a.IdAuditor = idAuditor
	}

	if idTemplate != 0 {
		a.IdTemplate = idTemplate
	}

	return a.OptionalSetIfNotEmptyOrNil(rating, observation)
}

func (a *Audit) OptionalSetIfNotEmptyOrNil(rating int64, observation string) error {

	if rating != 0 {
		a.Rating = zero.IntFrom(rating)
	}

	if observation != "" {
		a.Observation = zero.StringFrom(observation)
	}

	return nil
}

func (a *Audit) Deleted() bool {
	return a.meta.Deleted
}

func AAudit(allocateObjects bool) Audit {
	audit := Audit{}
	//if allocateObjects {
	//}
	return audit
}
func NewAudit(allocateObjects bool) *Audit {
	audit := AAudit(allocateObjects)
	return &audit
}

func (ds *Datastore) InsertAudit(a *Audit) error {
	return ds.InsertAuditTx(nil, a)
}

func (ds *Datastore) InsertAuditTx(tx *sql.Tx, a *Audit) error {

	if a == nil {
		return errors.New("audit should not be nil")
	}

	if a.Exists() {
		return errors.New("insert failed: already exists")
	}

	sql := ds.postgres.Rebind(`INSERT INTO places4all.audit (` +
		`id_property, id_auditor, id_template, rating, observation, created_date, finished_date` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7` +
		`) RETURNING id`)

	var err error
	if tx != nil {
		err = tx.QueryRow(sql, a.IdProperty, a.IdAuditor, a.IdTemplate, a.Rating, a.Observation, a.CreatedDate, a.FinishedDate).Scan(&a.Id)
	} else {
		err = ds.postgres.QueryRow(sql, a.IdProperty, a.IdAuditor, a.IdTemplate, a.Rating, a.Observation, a.CreatedDate, a.FinishedDate).Scan(&a.Id)
	}
	if err != nil {
		return err
	}

	a.SetExists()

	return nil
}

func (ds *Datastore) UpdateAudit(a *Audit) error {
	return ds.UpdateAuditTx(nil, a)
}

func (ds *Datastore) UpdateAuditTx(tx *sql.Tx, a *Audit) error {

	if a == nil {
		return errors.New("audit should not be nil")
	}

	if !a.Exists() {
		return errors.New("update failed: does not exist")
	}

	if a.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.audit SET (` +
		`id_property, id_auditor, id_template, rating, observation, created_date, finished_date` +
		`) = (` +
		`$1, $2, $3, $4, $5, $6, $7` +
		`) WHERE id = $8`

	var err error
	if tx != nil {
		_, err = tx.Exec(sql, a.IdProperty, a.IdAuditor, a.IdTemplate, a.Rating, a.Observation, a.CreatedDate, a.FinishedDate, a.Id)
	} else {
		_, err = ds.postgres.Exec(sql, a.IdProperty, a.IdAuditor, a.IdTemplate, a.Rating, a.Observation, a.CreatedDate, a.FinishedDate, a.Id)
	}
	return err
}

func (ds *Datastore) SaveAudit(a *Audit) error {
	return ds.SaveAuditTx(nil, a)
}

func (ds *Datastore) SaveAuditTx(tx *sql.Tx, a *Audit) error {

	if a == nil {
		return errors.New("audit should not be nil")
	}

	if a.Exists() {
		return ds.UpdateAuditTx(tx, a)
	}

	return ds.InsertAuditTx(tx, a)
}

func (ds *Datastore) UpsertAudit(a *Audit) error {

	if a == nil {
		return errors.New("audit should not be nil")
	}

	if a.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.audit (` +
		`id, id_property, id_auditor, id_template, rating, observation, created_date, finished_date` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_property, id_auditor, id_template, rating, observation, created_date, finished_date` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_property, EXCLUDED.id_auditor, EXCLUDED.id_template, EXCLUDED.rating, EXCLUDED.observation, EXCLUDED.created_date, EXCLUDED.finished_date` +
		`)`

	_, err := ds.postgres.Exec(sql, a.Id, a.IdProperty, a.IdAuditor, a.IdTemplate, a.Rating, a.Observation, a.CreatedDate, a.FinishedDate)
	if err != nil {
		return err
	}

	a.SetExists()

	return err
}

func (ds *Datastore) DeleteAudit(a *Audit) error {

	if a == nil {
		return errors.New("audit should not be nil")
	}

	if !a.Exists() {
		return nil
	}

	if a.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.audit WHERE id = $1`

	_, err := ds.postgres.Exec(sql, a.Id)
	if err != nil {
		return err
	}

	a.SetDeleted()

	return err
}

func (ds *Datastore) GetAuditAuditor(a *Audit, withEntity, restricted bool) (*Auditor, error) {
	return ds.GetAuditorById(a.IdAuditor, withEntity, restricted)
}

func (ds *Datastore) GetAuditProperty(a *Audit, withEntity, restricted bool) (*Property, error) {
	return ds.GetPropertyByIdWithAddressTagsOwners(a.IdProperty, withEntity, restricted)
}

func (ds *Datastore) GetAuditTemplate(a *Audit) (*Template, error) {
	return ds.GetTemplateById(a.IdTemplate)
}

func (ds *Datastore) GetAuditById(id int64) (*Audit, error) {

	const sql = `SELECT ` +
		`id, id_property, id_auditor, id_template, rating, observation, created_date, finished_date ` +
		`FROM places4all.audit ` +
		`WHERE id = $1`

	a := AAudit(true)
	a.SetExists()

	err := ds.postgres.QueryRowx(sql, id).StructScan(&a)
	if err != nil {
		return nil, err
	}

	return &a, err
}
func (ds *Datastore) GetAudits(limit, offset int, filter map[string]interface{}) ([]*Audit, error) {

	where, values := generators.GenerateAndSearchClause(filter)

	sql := ds.postgres.Rebind(`SELECT ` +
		`id, id_property, id_auditor, id_template, rating, observation, created_date, finished_date ` +
		`FROM places4all.audit ` +
		where +
		`ORDER BY audit.finished_date DESC LIMIT ` + strconv.Itoa(limit) +
		` OFFSET ` + strconv.Itoa(offset))

	rows, err := ds.postgres.Queryx(sql, values...)
	if err != nil {
		return nil, err
	}

	audit := make([]*Audit, 0)
	for rows.Next() {
		a := NewAudit(false)
		a.SetExists()
		err = rows.StructScan(a)
		if err != nil {
			return nil, err
		}
		audit = append(audit, a)
	}

	return audit, err
}

func (ds *Datastore) CheckTemplateUsed(idTemplate int64) (bool, error) {
	const query = `SELECT id ` +
		`FROM places4all.audit ` +
		`WHERE id_template = $1 LIMIT 1`

	var idAudit int64
	err := ds.postgres.QueryRow(query, idTemplate).Scan(&idAudit)
	if err == sql.ErrNoRows {
		return false, nil
	} else if err != nil {
		return false, err
	}

	return true, nil
}

func (ds *Datastore) DeleteAuditById(id int64) error {

	const sql = `DELETE FROM places4all.audit WHERE id = $1`

	_, err := ds.postgres.Exec(sql, id)
	if err != nil {
		return err
	}

	return err
}
