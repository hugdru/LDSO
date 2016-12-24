package datastore

import (
	"errors"
	"server/datastore/metadata"
	"strconv"
)

type Auditor struct {
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
		`)`

	_, err := ds.postgres.Exec(sql, a.IdEntity)
	if err != nil {
		return err
	}

	a.SetExists()

	return err
}

func (ds *Datastore) UpdateAuditor(a *Auditor) error {

	//if !a.Exists() {
	//	return errors.New("update failed: does not exist")
	//}
	//
	//if a.Deleted() {
	//	return errors.New("update failed: marked for deletion")
	//}
	//
	//const sql = `UPDATE places4all.auditor SET (` +
	//	`` +
	//	`) = ( ` +
	//	`$1` +
	//	`) WHERE id_entity = $2`
	//
	//_, err := ds.postgres.Exec(sql, a.IdEntity)
	//return err
	return errors.New("NOT IMPLEMENTED");
}

func (ds *Datastore) SaveAuditor(a *Auditor) error {
	if a.Exists() {
		return ds.UpdateAuditor(a)
	}

	return ds.InsertAuditor(a)
}

func (ds *Datastore) UpsertAuditor(a *Auditor) error {

	//if a.Exists() {
	//	return errors.New("insert failed: already exists")
	//}
	//
	//const sql = `INSERT INTO places4all.auditor (` +
	//	`id_entity` +
	//	`) VALUES (` +
	//	`$1` +
	//	`) ON CONFLICT (id_entity) DO UPDATE SET (` +
	//	`id_entity` +
	//	`) = (` +
	//	`EXCLUDED.id_entity` +
	//	`)`
	//
	//_, err := ds.postgres.Exec(sql, a.IdEntity)
	//if err != nil {
	//	return err
	//}
	//
	//a.SetExists()
	//
	//return err
	return errors.New("NOT IMPLEMENTED");
}

func (ds *Datastore) DeleteAuditor(a *Auditor) error {

	if !a.Exists() {
		return nil
	}

	if a.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.auditor WHERE id_entity = $1`

	_, err := ds.postgres.Exec(sql, a.IdEntity)
	if err != nil {
		return err
	}

	a.SetDeleted()

	return err
}

func (ds *Datastore) GetAuditorEntity(a *Auditor) (*Entity, error) {
	return ds.GetEntityById(a.IdEntity)
}

func (ds *Datastore) GetAuditors(limit, offset int) ([]*Auditor, error) {

	rows, err := ds.postgres.Queryx(ds.postgres.Rebind(`SELECT ` +
		`auditor.id_entity ` +
		`FROM places4all.auditor ` +
		`ORDER BY auditor.id_entity DESC LIMIT ` + strconv.Itoa(limit) +
		`OFFSET ` + strconv.Itoa(offset)))

	if err != nil {
		return nil, err
	}

	auditor := make([]*Auditor, 0)
	for rows.Next() {
		a := NewAuditor(false)
		a.SetExists()
		err = rows.StructScan(a)
		if err != nil {
			return nil, err
		}
		auditor = append(auditor, a)
	}

	return auditor, err
}

func (ds *Datastore) GetAuditorByIdWithForeign(idEntity int64) (*Auditor, error) {
	return ds.getAuditorById(idEntity, true)
}

func (ds *Datastore) GetAuditorById(idEntity int64) (*Auditor, error) {
	return ds.getAuditorById(idEntity, false)
}

func (ds *Datastore) getAuditorById(idEntity int64, withForeign bool) (*Auditor, error) {

	sql := ds.postgres.Rebind(`SELECT ` +
		`auditor.id_entity ` +
		`FROM places4all.auditor ` +
		`WHERE auditor.id_entity = $1`)

	a := AAuditor(false)
	a.SetExists()

	err := ds.postgres.QueryRow(sql, idEntity).Scan(&a.IdEntity)
	if err != nil {
		return nil, err
	}

	if withForeign {
		a.Entity, err = ds.GetEntityByIdWithForeign(idEntity)
		if err != nil {
			return nil, err
		}
	}

	return &a, err
}
