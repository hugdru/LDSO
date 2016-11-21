package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type LocalAdmin struct {
	Id       int64 `json:"id" db:"id"`
	IdEntity int64 `json:"idPerson" db:"id_entity"`
	meta     metadata.Metadata
}

func (l *LocalAdmin) SetExists() {
	l.meta.Exists = true
}

func (l *LocalAdmin) SetDeleted() {
	l.meta.Deleted = true
}

func (l *LocalAdmin) Exists() bool {
	return l.meta.Exists
}

func (l *LocalAdmin) Deleted() bool {
	return l.meta.Deleted
}

func (ds *Datastore) InsertLocalAdmin(l *LocalAdmin) error {
	var err error

	if l.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.localadmin (` +
		`id_person` +
		`) VALUES (` +
		`$1` +
		`) RETURNING id`

	err = ds.postgres.QueryRow(sql, l.IdPerson).Scan(&l.Id)
	if err != nil {
		return err
	}

	l.SetExists()

	return nil
}

func (ds *Datastore) UpdateLocalAdmin(l *LocalAdmin) error {
	var err error

	if !l.Exists() {
		return errors.New("update failed: does not exist")
	}

	if l.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.localadmin SET (` +
		`id_person` +
		`) = ( ` +
		`$1` +
		`) WHERE id = $2`

	_, err = ds.postgres.Exec(sql, l.IdPerson, l.Id)
	return err
}

func (ds *Datastore) SaveLocalAdmin(l *LocalAdmin) error {
	if l.Exists() {
		return ds.UpdateLocalAdmin(l)
	}

	return ds.InsertLocalAdmin(l)
}

func (ds *Datastore) UpsertLocalAdmin(l *LocalAdmin) error {
	var err error

	if l.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sqlstr = `INSERT INTO places4all.localadmin (` +
		`id, id_person` +
		`) VALUES (` +
		`$1, $2` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_person` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_person` +
		`)`

	_, err = ds.postgres.Exec(sqlstr, l.Id, l.IdPerson)
	if err != nil {
		return err
	}

	l.SetExists()

	return nil
}

func (ds *Datastore) DeleteLocalAdmin(l *LocalAdmin) error {
	var err error

	if !l.Exists() {
		return nil
	}

	if l.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.localadmin WHERE id = $1`

	_, err = ds.postgres.Exec(sql, l.Id)
	if err != nil {
		return err
	}

	l.SetDeleted()

	return nil
}

func (ds *Datastore) GetLocalAdminPerson(l *LocalAdmin) (*Person, error) {
	return ds.GetPersonById(l.IdPerson)
}

func (ds *Datastore) GetLocalAdminById(id int64) (*LocalAdmin, error) {
	var err error

	const sql = `SELECT ` +
		`id, id_person ` +
		`FROM places4all.localadmin ` +
		`WHERE id = $1`

	l := LocalAdmin{}
	l.SetExists()

	err = ds.postgres.QueryRow(sql, id).Scan(&l.Id, &l.IdPerson)
	if err != nil {
		return nil, err
	}

	return &l, nil
}
