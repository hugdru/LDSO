package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type LocalAdmin struct {
	Id       int64 `json:"id" db:"id"`
	IdEntity int64 `json:"IdEntity" db:"id_entity"`

	// Objects
	Entity *Entity `json:"entity,omitempty"`

	meta metadata.Metadata
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

func ALocalAdmin(allocateObjects bool) LocalAdmin {
	localAdmin := LocalAdmin{}
	if allocateObjects {
		localAdmin.Entity = NewEntity(true)
	}
	return localAdmin
}
func NewLocalAdmin(allocateObjects bool) *LocalAdmin {
	localAdmin := ALocalAdmin(allocateObjects)
	return &localAdmin
}

func (ds *Datastore) InsertLocalAdmin(l *LocalAdmin) error {

	if l.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.localadmin (` +
		`id_entity` +
		`) VALUES (` +
		`$1` +
		`) RETURNING id`

	res, err := ds.postgres.Exec(sql, l.IdEntity)
	if err != nil {
		return err
	}
	l.Id, err = res.LastInsertId()
	if err != nil {
		return err
	}

	l.SetExists()

	return err
}

func (ds *Datastore) UpdateLocalAdmin(l *LocalAdmin) error {

	if !l.Exists() {
		return errors.New("update failed: does not exist")
	}

	if l.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.localadmin SET (` +
		`id_entity` +
		`) = ( ` +
		`$1` +
		`) WHERE id = $2`

	_, err := ds.postgres.Exec(sql, l.IdEntity, l.Id)
	return err
}

func (ds *Datastore) SaveLocalAdmin(l *LocalAdmin) error {
	if l.Exists() {
		return ds.UpdateLocalAdmin(l)
	}

	return ds.InsertLocalAdmin(l)
}

func (ds *Datastore) UpsertLocalAdmin(l *LocalAdmin) error {

	if l.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.localadmin (` +
		`id, id_entity` +
		`) VALUES (` +
		`$1, $2` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_entity` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_entity` +
		`)`

	_, err := ds.postgres.Exec(sql, l.Id, l.IdEntity)
	if err != nil {
		return err
	}

	l.SetExists()

	return err
}

func (ds *Datastore) DeleteLocalAdmin(l *LocalAdmin) error {

	if !l.Exists() {
		return nil
	}

	if l.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.localadmin WHERE id = $1`

	_, err := ds.postgres.Exec(sql, l.Id)
	if err != nil {
		return err
	}

	l.SetDeleted()

	return err
}

func (ds *Datastore) GetLocalAdminEntity(l *LocalAdmin) (*Entity, error) {
	return ds.GetEntityById(l.IdEntity)
}

func (ds *Datastore) GetLocalAdminById(id int64) (*LocalAdmin, error) {

	const sql = `SELECT ` +
		`id, id_entity ` +
		`FROM places4all.localadmin ` +
		`WHERE id = $1`

	l := ALocalAdmin(false)
	l.SetExists()

	err := ds.postgres.QueryRow(sql, id).Scan(&l.Id, &l.IdEntity)
	if err != nil {
		return nil, err
	}

	return &l, err
}
