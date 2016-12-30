package datastore

import (
	"database/sql"
	"errors"
	"server/datastore/generators"
	"server/datastore/metadata"
	"strconv"
)

func auditorVisibility(restricted bool) string {
	const auditorRestricted = "auditor.id_entity "
	const auditorAll = "auditor.id_entity "
	if restricted {
		return auditorRestricted
	}
	return auditorAll
}

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

func (a *Auditor) MustSet(idEntity int64) error {
	if idEntity != 0 {
		a.IdEntity = idEntity
	} else {
		return errors.New("idEntity must be set")
	}

	return nil
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
	return ds.InsertAuditorTx(nil, a)
}

func (ds *Datastore) InsertAuditorTx(tx *sql.Tx, a *Auditor) error {

	if a == nil {
		return errors.New("auditor should not be nil")
	}

	if a.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.auditor (` +
		`id_entity` +
		`) VALUES (` +
		`$1` +
		`)`

	var err error
	if tx != nil {
		_, err = tx.Exec(sql, a.IdEntity)

	} else {
		_, err = ds.postgres.Exec(sql, a.IdEntity)

	}
	if err != nil {
		return err
	}

	a.SetExists()

	return err
}

func (ds *Datastore) UpdateAuditor(a *Auditor) error {
	return ds.UpdateAuditorTx(nil, a)
}

func (ds *Datastore) UpdateAuditorTx(tx *sql.Tx, a *Auditor) error {

	//if a == nil {
	//	return errors.New("auditor should not be nil")
	//}
	//
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
	return errors.New("NOT IMPLEMENTED")
}

func (ds *Datastore) SaveAuditor(a *Auditor) error {
	return ds.SaveAuditorTx(nil, a)
}

func (ds *Datastore) SaveAuditorTx(tx *sql.Tx, a *Auditor) error {

	if a == nil {
		return errors.New("auditor should not be nil")
	}

	if a.Exists() {
		return ds.UpdateAuditorTx(tx, a)
	}

	return ds.InsertAuditorTx(tx, a)
}

func (ds *Datastore) UpsertAuditor(a *Auditor) error {
	//if a == nil {
	//	return errors.New("auditor should not be nil")
	//}
	//
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
	return errors.New("NOT IMPLEMENTED")
}

func (ds *Datastore) DeleteAuditor(a *Auditor) error {

	if a == nil {
		return errors.New("auditor should not be nil")
	}

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

func (ds *Datastore) GetAuditorEntity(a *Auditor, withCountry, restricted bool) (*Entity, error) {
	return ds.GetEntityById(a.IdEntity, withCountry, restricted)
}

func (ds *Datastore) GetAuditors(limit, offset int, withEntity, restricted bool) ([]*Auditor, error) {
	rows, err := ds.postgres.Queryx(ds.postgres.Rebind(`SELECT ` +
		auditorVisibility(restricted) +
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
		if withEntity {
			a.Entity, err = ds.GetEntityById(a.IdEntity, true, restricted)
			if err != nil {
				return nil, err
			}
		}
		auditor = append(auditor, a)
	}

	return auditor, err
}

func (ds *Datastore) GetAuditorById(idEntity int64, withEntity, restricted bool) (*Auditor, error) {
	filter := make(map[string]interface{})
	filter["id_entity"] = idEntity
	return ds.GetAuditor(filter, withEntity, restricted)
}

func (ds *Datastore) GetAuditor(filter map[string]interface{}, withEntity, restricted bool) (*Auditor, error) {
	where, values := generators.GenerateAndSearchClause(filter)
	sql := ds.postgres.Rebind(`SELECT ` +
		auditorVisibility(restricted) +
		`FROM places4all.auditor ` +
		where)

	a := AAuditor(false)
	a.SetExists()

	err := ds.postgres.QueryRowx(sql, values...).StructScan(&a)
	if err != nil {
		return nil, err
	}

	if withEntity {
		a.Entity, err = ds.GetEntityById(a.IdEntity, true, restricted)
		if err != nil {
			return nil, err
		}
	}

	return &a, err
}
