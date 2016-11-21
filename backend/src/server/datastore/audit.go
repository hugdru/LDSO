package datastore

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"server/datastore/metadata"
	"time"
)

type Audit struct {
	Id          int64          `json:"id" db:"id"`
	IdProperty  int64          `json:"id_property" db:"id_property"`
	IdAuditor   int64          `json:"id_auditor" db:"id_auditor"`
	IdTemplate  int64          `json:"id_template" db:"id_template"`
	Rating      sql.NullInt64  `json:"rating" db:"rating"`
	Observation sql.NullString `json:"observation" db:"observation"`
	Created     *time.Time     `json:"created" db:"created"`
	Finished    pq.NullTime    `json:"finished" db:"finished"`
	meta        metadata.Metadata
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

func (ds *Datastore) InsertAudit(a *Audit) error {
	var err error

	if a.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.audit (` +
		`id_property, id_auditor, id_template, rating, observation, created, finished` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7` +
		`) RETURNING id`

	err = ds.postgres.QueryRow(sql, a.IdProperty, a.IdAuditor, a.IdTemplate, a.Rating, a.Observation, a.Created, a.Finished).Scan(&a.Id)
	if err != nil {
		return err
	}

	a.SetExists()

	return nil
}

func (ds *Datastore) UpdateAudit(a *Audit) error {
	var err error

	if !a.Exists() {
		return errors.New("update failed: does not exist")
	}

	if a.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.audit SET (` +
		`id_property, id_auditor, id_template, rating, observation, created, finished` +
		`) = ( ` +
		`$1, $2, $3, $4, $5, $6, $7` +
		`) WHERE id = $8`

	_, err = ds.postgres.Exec(sql, a.IdProperty, a.IdAuditor, a.IdTemplate, a.Rating, a.Observation, a.Created, a.Finished, a.Id)
	return err
}

func (ds *Datastore) SaveAudit(a *Audit) error {
	if a.Exists() {
		return ds.UpdateAudit(a)
	}

	return ds.InsertAudit(a)
}

func (ds *Datastore) UpsertAudit(a *Audit) error {
	var err error

	if a.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.audit (` +
		`id, id_property, id_auditor, id_template, rating, observation, created, finished` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_property, id_auditor, id_template, rating, observation, created, finished` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_property, EXCLUDED.id_auditor, EXCLUDED.id_template, EXCLUDED.rating, EXCLUDED.observation, EXCLUDED.created, EXCLUDED.finished` +
		`)`

	_, err = ds.postgres.Exec(sql, a.Id, a.IdProperty, a.IdAuditor, a.IdTemplate, a.Rating, a.Observation, a.Created, a.Finished)
	if err != nil {
		return err
	}

	a.SetExists()

	return nil
}

func (ds *Datastore) DeleteAudit(a *Audit) error {
	var err error

	if !a.Exists() {
		return nil
	}

	if a.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.audit WHERE id = $1`

	_, err = ds.postgres.Exec(sql, a.Id)
	if err != nil {
		return err
	}

	a.SetDeleted()

	return nil
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
	var err error

	const sql = `SELECT ` +
		`id, id_property, id_auditor, id_template, rating, observation, created, finished ` +
		`FROM places4all.audit ` +
		`WHERE id = $1`

	a := Audit{}
	a.SetExists()

	err = ds.postgres.QueryRow(sql, id).Scan(&a.Id, &a.IdProperty, &a.IdAuditor, &a.IdTemplate, &a.Rating, &a.Observation, &a.Created, &a.Finished)
	if err != nil {
		return nil, err
	}

	return &a, nil
}
