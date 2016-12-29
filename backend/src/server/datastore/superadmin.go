package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type Superadmin struct {
	IdEntity int64 `json:"idEntity" db:"id_entity"`

	// Objects
	Entity *Entity `json:"entity,omitempty"`

	meta metadata.Metadata
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

func ASuperadmin(allocateObjects bool) Superadmin {
	superAdmin := Superadmin{}
	if allocateObjects {
		superAdmin.Entity = NewEntity(allocateObjects)
	}
	return superAdmin
}

func NewSuperadmin(allocateObjects bool) *Superadmin {
	superAdmin := ASuperadmin(allocateObjects)
	return &superAdmin
}

func (ds *Datastore) InsertSuperadmin(s *Superadmin) error {

	if s == nil {
		return errors.New("superadmin should not be nil")
	}

	if s.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.superadmin (` +
		`id_entity` +
		`) VALUES (` +
		`$1` +
		`)`

	_, err := ds.postgres.Exec(sql, s.IdEntity)
	if err != nil {
		return err
	}

	s.SetExists()

	return err
}

func (ds *Datastore) UpdateSuperadmin(s *Superadmin) error {
	//if s == nil {
	//	return errors.New("superadmin should not be nil")
	//}
	//
	//if !s.Exists() {
	//	return errors.New("update failed: does not exist")
	//}
	//
	//if s.Deleted() {
	//	return errors.New("update failed: marked for deletion")
	//}
	//
	//const sql = `UPDATE places4all.superadmin SET (` +
	//	`id_entity` +
	//	`) = ( ` +
	//	`$1` +
	//	`) WHERE id = $2`
	//
	//_, err := ds.postgres.Exec(sql, s.IdEntity, s.Id)
	//return err
	return errors.New("NOT IMPLEMENTED")
}

func (ds *Datastore) SaveSuperadmin(s *Superadmin) error {
	if s.Exists() {
		return ds.UpdateSuperadmin(s)
	}

	return ds.InsertSuperadmin(s)
}

func (ds *Datastore) UpsertSuperadmin(s *Superadmin) error {

	if s == nil {
		return errors.New("superadmin should not be nil")
	}

	if s.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.superadmin (` +
		`id_entity` +
		`) VALUES (` +
		`$1` +
		`) ON CONFLICT (id_entity) DO UPDATE SET (` +
		`id_entity` +
		`) = (` +
		`EXCLUDED.id_entity` +
		`)`

	_, err := ds.postgres.Exec(sql, s.IdEntity, s.IdEntity)
	if err != nil {
		return err
	}

	s.SetExists()

	return nil
}

func (ds *Datastore) DeleteSuperadmin(s *Superadmin) error {

	if s == nil {
		return errors.New("superadmin should not be nil")
	}

	if !s.Exists() {
		return nil
	}

	if s.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.superadmin WHERE id_entity = $1`

	_, err := ds.postgres.Exec(sql, s.IdEntity)
	if err != nil {
		return err
	}

	s.SetDeleted()

	return err
}

func (ds *Datastore) GetSuperadminEntity(s *Superadmin) (*Entity, error) {
	return ds.GetEntityById(s.IdEntity)
}

func (ds *Datastore) GetSuperadminByIdWithEntity(idEntity int64) (*Superadmin, error) {
	return ds.getSuperadminById(idEntity, true)
}

func (ds *Datastore) GetSuperadminById(idEntity int64) (*Superadmin, error) {
	return ds.getSuperadminById(idEntity, false)
}

func (ds *Datastore) getSuperadminById(idEntity int64, withForeign bool) (*Superadmin, error) {

	const sql = `SELECT ` +
		`id_entity ` +
		`FROM places4all.superadmin ` +
		`WHERE id_entity = $1`

	s := ASuperadmin(false)
	s.SetExists()

	err := ds.postgres.QueryRow(sql, idEntity).Scan(&s.IdEntity)
	if err != nil {
		return nil, err
	}

	if withForeign {
		s.Entity, err = ds.GetEntityByIdWithCountry(idEntity)
		if err != nil {
			return nil, err
		}
	}

	return &s, err

}
