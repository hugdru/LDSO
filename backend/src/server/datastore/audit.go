package datastore

import (
	"errors"
	"gopkg.in/guregu/null.v3/zero"
	"server/datastore/metadata"
	"time"
)

type Audit struct {
	Id           int64       `json:"id" db:"id"`
	IdProperty   int64       `json:"id_property" db:"id_property"`
	IdAuditor    int64       `json:"id_auditor" db:"id_auditor"`
	IdTemplate   int64       `json:"id_template" db:"id_template"`
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

	if a.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.audit (` +
		`id_property, id_auditor, id_template, rating, observation, created_date, finished_date` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7` +
		`) RETURNING id`

	err := ds.postgres.QueryRow(sql, a.IdProperty, a.IdAuditor, a.IdTemplate, a.Rating, a.Observation, a.CreatedDate, a.FinishedDate).Scan(&a.Id)
	if err != nil {
		return err
	}

	a.SetExists()

	return nil
}

func (ds *Datastore) UpdateAudit(a *Audit) error {

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

	_, err := ds.postgres.Exec(sql, a.IdProperty, a.IdAuditor, a.IdTemplate, a.Rating, a.Observation, a.CreatedDate, a.FinishedDate, a.Id)
	return err
}

func (ds *Datastore) SaveAudit(a *Audit) error {
	if a.Exists() {
		return ds.UpdateAudit(a)
	}

	return ds.InsertAudit(a)
}

func (ds *Datastore) UpsertAudit(a *Audit) error {

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

func (ds *Datastore) GetAuditAuditor(a *Audit) (*Auditor, error) {
	return ds.GetAuditorById(a.IdAuditor)
}

func (ds *Datastore) GetAuditProperty(a *Audit) (*Property, error) {
	return ds.GetPropertyById(a.IdProperty)
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
