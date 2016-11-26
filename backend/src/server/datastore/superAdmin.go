package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type SuperAdmin struct {
	Id       int64 `json:"id" db:"id"`
	IdEntity int64 `json:"idEntity" db:"id_entity"`

	// Objects
	Entity *Entity `json:"entity,omitempty"`

	meta metadata.Metadata
}

func (s *SuperAdmin) SetExists() {
	s.meta.Exists = true
}

func (s *SuperAdmin) SetDeleted() {
	s.meta.Deleted = true
}

func (s *SuperAdmin) Exists() bool {
	return s.meta.Exists
}

func (s *SuperAdmin) Deleted() bool {
	return s.meta.Deleted
}

func ASuperAdmin(allocateObjects bool) SuperAdmin {
	superAdmin := SuperAdmin{}
	if allocateObjects {
		superAdmin.Entity = NewEntity(allocateObjects)
	}
	return superAdmin
}

func NewSuperAdmin(allocateObjects bool) *SuperAdmin {
	superAdmin := ASuperAdmin(allocateObjects)
	return &superAdmin
}

func (ds *Datastore) InsertSuperAdmin(s *SuperAdmin) error {

	if s.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.superadmin (` +
		`id_entity` +
		`) VALUES (` +
		`$1` +
		`) RETURNING id`

	res, err := ds.postgres.Exec(sql, s.IdEntity)
	if err != nil {
		return err
	}
	s.Id, err = res.LastInsertId()
	if err != nil {
		return err
	}

	s.SetExists()

	return err
}

func (ds *Datastore) UpdateSuperAdmin(s *SuperAdmin) error {

	if !s.Exists() {
		return errors.New("update failed: does not exist")
	}

	if s.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.superadmin SET (` +
		`id_entity` +
		`) = ( ` +
		`$1` +
		`) WHERE id = $2`

	_, err := ds.postgres.Exec(sql, s.IdEntity, s.Id)
	return err
}

func (ds *Datastore) SaveSuperAdmin(s *SuperAdmin) error {
	if s.Exists() {
		return ds.UpdateSuperAdmin(s)
	}

	return ds.InsertSuperAdmin(s)
}

func (ds *Datastore) UpsertSuperAdmin(s *SuperAdmin) error {

	if s.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.superadmin (` +
		`id, id_entity` +
		`) VALUES (` +
		`$1, $2` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_entity` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_entity` +
		`)`

	_, err := ds.postgres.Exec(sql, s.Id, s.IdEntity)
	if err != nil {
		return err
	}

	s.SetExists()

	return nil
}

func (ds *Datastore) DeleteSuperAdmin(s *SuperAdmin) error {

	if !s.Exists() {
		return nil
	}

	if s.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.superadmin WHERE id = $1`

	_, err := ds.postgres.Exec(sql, s.Id)
	if err != nil {
		return err
	}

	s.SetDeleted()

	return err
}

func (ds *Datastore) GetSuperAdminEntity(s *SuperAdmin) (*Entity, error) {
	return ds.GetEntityById(s.IdEntity)
}

func (ds *Datastore) GetSuperAdminById(id int64) (*SuperAdmin, error) {

	const sql = `SELECT ` +
		`id, id_entity ` +
		`FROM places4all.superadmin ` +
		`WHERE id = $1`

	s := ASuperAdmin(false)
	s.SetExists()

	err := ds.postgres.QueryRow(sql, id).Scan(&s.Id, &s.IdEntity)
	if err != nil {
		return nil, err
	}

	return &s, err
}
