package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type Auditor struct {
	Id       int64 `json:"id" db:"id"`
	IdEntity int64 `json:"IdEntity" db:"id_entity"`

	// Objects
	Entity *Entity `json:"entity,omitempty"`

	meta metadata.Metadata
}

func (a *Auditor) SetExists() {
	a.meta.Exists = true
}

func (a *Auditor) SetDeleted() {
	a.meta.Deleted = true
}

func (a *Auditor) Exists() bool {
	return a.meta.Exists
}

func (a *Auditor) Deleted() bool {
	return a.meta.Deleted
}

func AAuditor(allocateObjects bool) Auditor {
	auditor := Auditor{}
	if allocateObjects {
		auditor.Entity = NewEntity(true)
	}
	return auditor
}
func NewAuditor(allocateObjects bool) *Auditor {
	auditor := AAuditor(allocateObjects)
	return &auditor
}

func (ds *Datastore) InsertAuditor(a *Auditor) error {

	if a.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.auditor (` +
		`id_entity` +
		`) VALUES (` +
		`$1` +
		`) RETURNING id`

	res, err := ds.postgres.Exec(sql, a.IdEntity)
	if err != nil {
		return err
	}
	a.Id, err = res.LastInsertId()
	if err != nil {
		return err
	}

	a.SetExists()

	return err
}

func (ds *Datastore) UpdateAuditor(a *Auditor) error {

	if !a.Exists() {
		return errors.New("update failed: does not exist")
	}

	if a.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.auditor SET (` +
		`id_entity` +
		`) = ( ` +
		`$1` +
		`) WHERE id = $2`

	_, err := ds.postgres.Exec(sql, a.IdEntity, a.Id)
	return err
}

func (ds *Datastore) SaveAuditor(a *Auditor) error {
	if a.Exists() {
		return ds.UpdateAuditor(a)
	}

	return ds.InsertAuditor(a)
}

func (ds *Datastore) UpsertAuditor(a *Auditor) error {

	if a.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.auditor (` +
		`id, id_entity` +
		`) VALUES (` +
		`$1, $2` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_entity` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_entity` +
		`)`

	_, err := ds.postgres.Exec(sql, a.Id, a.IdEntity)
	if err != nil {
		return err
	}

	a.SetExists()

	return err
}

func (ds *Datastore) DeleteAuditor(a *Auditor) error {

	if !a.Exists() {
		return nil
	}

	if a.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.auditor WHERE id = $1`

	_, err := ds.postgres.Exec(sql, a.Id)
	if err != nil {
		return err
	}

	a.SetDeleted()

	return err
}

func (ds *Datastore) GetAuditorEntity(a *Auditor) (*Entity, error) {
	return ds.GetEntityById(a.IdEntity)
}

func (ds *Datastore) GetAuditorById(id int64) (*Auditor, error) {
	var err error

	const sql = `SELECT ` +
		`id, id_entity ` +
		`FROM places4all.auditor ` +
		`WHERE id = $1`

	a := AAuditor(false)
	a.SetExists()

	err = ds.postgres.QueryRowx(sql, id).Scan(&a.Id, &a.IdEntity)
	if err != nil {
		return nil, err
	}

	return &a, err
}
