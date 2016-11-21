package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type Auditor struct {
	Id       int64 `json:"id" db:"id"`
	IdEntity int64 `json:"idPerson" db:"id_entity"`
	meta     metadata.Metadata
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

func (ds *Datastore) InsertAuditor(a *Auditor) error {
	var err error

	if a.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.auditor (` +
		`id_person` +
		`) VALUES (` +
		`$1` +
		`) RETURNING id`

	err = ds.postgres.QueryRow(sql, a.IdPerson).Scan(&a.Id)
	if err != nil {
		return err
	}

	a.SetExists()

	return nil
}

func (ds *Datastore) UpdateAuditor(a *Auditor) error {
	var err error

	if !a.Exists() {
		return errors.New("update failed: does not exist")
	}

	if a.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.auditor SET (` +
		`id_person` +
		`) = ( ` +
		`$1` +
		`) WHERE id = $2`

	_, err = ds.postgres.Exec(sql, a.IdPerson, a.Id)
	return err
}

func (ds *Datastore) SaveAuditor(a *Auditor) error {
	if a.Exists() {
		return ds.UpdateAuditor(a)
	}

	return ds.InsertAuditor(a)
}

func (ds *Datastore) UpsertAuditor(a *Auditor) error {
	var err error

	if a.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.auditor (` +
		`id, id_person` +
		`) VALUES (` +
		`$1, $2` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_person` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_person` +
		`)`

	_, err = ds.postgres.Exec(sql, a.Id, a.IdPerson)
	if err != nil {
		return err
	}

	a.SetExists()

	return nil
}

func (ds *Datastore) DeleteAuditor(a *Auditor) error {
	var err error

	if !a.Exists() {
		return nil
	}

	if a.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.auditor WHERE id = $1`

	_, err = ds.postgres.Exec(sql, a.Id)
	if err != nil {
		return err
	}

	a.SetDeleted()

	return nil
}

func (ds *Datastore) GetAuditorPerson(a *Auditor) (*Person, error) {
	return ds.GetPersonById(a.IdPerson)
}

func (ds *Datastore) GetAuditorById(id int64) (*Auditor, error) {
	var err error

	const sql = `SELECT ` +
		`id, id_person ` +
		`FROM places4all.auditor ` +
		`WHERE id = $1`

	a := Auditor{}
	a.SetExists()

	err = ds.postgres.QueryRow(sql, id).Scan(&a.Id, &a.IdPerson)
	if err != nil {
		return nil, err
	}

	return &a, nil
}
