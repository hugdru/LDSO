package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type Localadmin struct {
	IdEntity int64 `json:"IdEntity" db:"id_entity"`

	// Objects
	Entity *Entity `json:"entity,omitempty"`

	meta metadata.Metadata
}

func (l *Localadmin) SetExists() {
	l.meta.Exists = true
}

func (l *Localadmin) SetDeleted() {
	l.meta.Deleted = true
}

func (l *Localadmin) Exists() bool {
	return l.meta.Exists
}

func (l *Localadmin) Deleted() bool {
	return l.meta.Deleted
}

func ALocaladmin(allocateObjects bool) Localadmin {
	localAdmin := Localadmin{}
	if allocateObjects {
		localAdmin.Entity = NewEntity(true)
	}
	return localAdmin
}
func NewLocaladmin(allocateObjects bool) *Localadmin {
	localAdmin := ALocaladmin(allocateObjects)
	return &localAdmin
}

func (ds *Datastore) InsertLocaladmin(l *Localadmin) error {

	if l == nil {
		return errors.New("localadmin should not be nil")
	}

	if l.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.localadmin (` +
		`id_entity` +
		`) VALUES (` +
		`$1` +
		`)`

	_, err := ds.postgres.Exec(sql, l.IdEntity)
	if err != nil {
		return err
	}

	l.SetExists()

	return err
}

func (ds *Datastore) UpdateLocaladmin(l *Localadmin) error {

	//if l == nil {
	//	return errors.New("localadmin should not be nil")
	//}

	//if !l.Exists() {
	//	return errors.New("update failed: does not exist")
	//}
	//
	//if l.Deleted() {
	//	return errors.New("update failed: marked for deletion")
	//}
	//
	//const sql = `UPDATE places4all.localadmin SET (` +
	//	`id_entity` +
	//	`) = ( ` +
	//	`$1` +
	//	`) WHERE id = $2`
	//
	//_, err := ds.postgres.Exec(sql, l.IdEntity, l.Id)
	//return err
	return errors.New("NOT IMPLEMENTED")
}

func (ds *Datastore) SaveLocaladmin(l *Localadmin) error {
	if l.Exists() {
		return ds.UpdateLocaladmin(l)
	}

	return ds.InsertLocaladmin(l)
}

func (ds *Datastore) UpsertLocalAdmin(l *Localadmin) error {

	if l == nil {
		return errors.New("localadmin should not be nil")
	}

	if l.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.localadmin (` +
		`id_entity` +
		`) VALUES (` +
		`$1` +
		`) ON CONFLICT (id_entity) DO UPDATE SET (` +
		`id_entity` +
		`) = (` +
		`EXCLUDED.id_entity` +
		`)`

	_, err := ds.postgres.Exec(sql, l.IdEntity)
	if err != nil {
		return err
	}

	l.SetExists()

	return err
}

func (ds *Datastore) DeleteLocaladmin(l *Localadmin) error {

	if l == nil {
		return errors.New("localadmin should not be nil")
	}

	if !l.Exists() {
		return nil
	}

	if l.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.localadmin WHERE id_entity = $1`

	_, err := ds.postgres.Exec(sql, l.IdEntity)
	if err != nil {
		return err
	}

	l.SetDeleted()

	return err
}

func (ds *Datastore) GetLocaladminEntity(l *Localadmin) (*Entity, error) {
	return ds.GetEntityById(l.IdEntity)
}

func (ds *Datastore) GetLocaladminByIdWithForeign(idEntity int64) (*Localadmin, error) {
	return ds.getLocaladminById(idEntity, true)
}

func (ds *Datastore) GetLocaladminById(idEntity int64) (*Localadmin, error) {
	return ds.getLocaladminById(idEntity, false)
}

func (ds *Datastore) getLocaladminById(idEntity int64, withForeign bool) (*Localadmin, error) {

	const sql = `SELECT ` +
		`id_entity ` +
		`FROM places4all.localadmin ` +
		`WHERE id_entity = $1`

	l := ALocaladmin(false)
	l.SetExists()

	err := ds.postgres.QueryRow(sql, idEntity).Scan(&l.IdEntity)
	if err != nil {
		return nil, err
	}

	if withForeign {
		l.Entity, err = ds.GetEntityByIdWithForeign(idEntity)
		if err != nil {
			return nil, err
		}
	}

	return &l, err
}
