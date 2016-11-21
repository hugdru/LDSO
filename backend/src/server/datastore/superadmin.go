package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type Superadmin struct {
	Id       int64 `json:"id" db:"id"`
	IdEntity int64 `json:"idPerson" db:"id_entity"`
	meta     metadata.Metadata
}

func (s *Superadmin) SetExists() {
	s.meta.Exists = true
}

func (s *Superadmin) SetDeleted() {
	s.meta.Deleted = true
}

func (s *Superadmin) Exists() bool {
	return s.meta.Exists
}

func (s *Superadmin) Deleted() bool {
	return s.meta.Deleted
}

func (ds *Datastore) InsertSuperadmin(s *Superadmin) error {
	var err error

	if s.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.superadmin (` +
		`id_person` +
		`) VALUES (` +
		`$1` +
		`) RETURNING id`

	err = ds.postgres.QueryRow(sql, s.IdPerson).Scan(&s.Id)
	if err != nil {
		return err
	}

	s.SetExists()

	return nil
}

func (ds *Datastore) UpdateSuperadmin(s *Superadmin) error {
	var err error

	if !s.Exists() {
		return errors.New("update failed: does not exist")
	}

	if s.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.superadmin SET (` +
		`id_person` +
		`) = ( ` +
		`$1` +
		`) WHERE id = $2`

	_, err = ds.postgres.Exec(sql, s.IdPerson, s.Id)
	return err
}

func (ds *Datastore) SaveSuperadmin(s *Superadmin) error {
	if s.Exists() {
		return ds.UpdateSuperadmin(s)
	}

	return ds.InsertSuperadmin(s)
}

func (ds *Datastore) UpsertSuperadmin(s *Superadmin) error {
	var err error

	if s.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.superadmin (` +
		`id, id_person` +
		`) VALUES (` +
		`$1, $2` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_person` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_person` +
		`)`

	_, err = ds.postgres.Exec(sql, s.Id, s.IdPerson)
	if err != nil {
		return err
	}

	s.SetExists()

	return nil
}

func (ds *Datastore) DeleteSuperadmin(s *Superadmin) error {
	var err error

	if !s.Exists() {
		return nil
	}

	if s.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.superadmin WHERE id = $1`

	_, err = ds.postgres.Exec(sql, s.Id)
	if err != nil {
		return err
	}

	s.SetDeleted()

	return nil
}

func (ds *Datastore) GetSuperadminPerson(s *Superadmin) (*Person, error) {
	return ds.GetPersonById(s.IdPerson)
}

func (ds *Datastore) GetSuperadminById(id int64) (*Superadmin, error) {
	var err error

	const sql = `SELECT ` +
		`id, id_person ` +
		`FROM places4all.superadmin ` +
		`WHERE id = $1`

	s := Superadmin{}
	s.SetExists()

	err = ds.postgres.QueryRow(sql, id).Scan(&s.Id, &s.IdPerson)
	if err != nil {
		return nil, err
	}

	return &s, nil
}
