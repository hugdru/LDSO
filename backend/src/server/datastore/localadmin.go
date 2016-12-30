package datastore

import (
	"database/sql"
	"errors"
	"server/datastore/generators"
	"server/datastore/metadata"
	"strconv"
)

func localadminVisibility(restricted bool) string {
	const localadminRestricted = "localadmin.id_entity "
	const localadminAll = "localadmin.id_entity "
	if restricted {
		return localadminRestricted
	}
	return localadminAll
}

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

func (l *Localadmin) MustSet(idEntity int64) error {
	if idEntity != 0 {
		l.IdEntity = idEntity
	} else {
		return errors.New("idEntity must be set")
	}

	return nil
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
	return ds.InsertLocaladminTx(nil, l)
}

func (ds *Datastore) InsertLocaladminTx(tx *sql.Tx, l *Localadmin) error {

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

	var err error
	if tx != nil {
		_, err = tx.Exec(sql, l.IdEntity)

	} else {
		_, err = ds.postgres.Exec(sql, l.IdEntity)
	}
	if err != nil {
		return err
	}

	l.SetExists()

	return err
}

func (ds *Datastore) UpdateLocaladmin(l *Localadmin) error {
	return ds.UpdateLocaladminTx(nil, l)
}

func (ds *Datastore) UpdateLocaladminTx(tx *sql.Tx, l *Localadmin) error {

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
	return ds.SaveLocaladminTx(nil, l)
}

func (ds *Datastore) SaveLocaladminTx(tx *sql.Tx, l *Localadmin) error {

	if l == nil {
		return errors.New("localadmin should not be nil")
	}

	if l.Exists() {
		return ds.UpdateLocaladminTx(tx, l)
	}

	return ds.InsertLocaladminTx(tx, l)
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

func (ds *Datastore) GetLocaladminEntity(l *Localadmin, withCountry, restricted bool) (*Entity, error) {
	return ds.GetEntityById(l.IdEntity, withCountry, restricted)
}

func (ds *Datastore) GetLocaladmins(limit, offset int, withEntity, restricted bool) ([]*Localadmin, error) {
	rows, err := ds.postgres.Queryx(ds.postgres.Rebind(`SELECT ` +
		localadminVisibility(restricted) +
		`FROM places4all.localadmin ` +
		`ORDER BY localadmin.id_entity DESC LIMIT ` + strconv.Itoa(limit) +
		`OFFSET ` + strconv.Itoa(offset)))

	if err != nil {
		return nil, err
	}

	localadmin := make([]*Localadmin, 0)
	for rows.Next() {
		c := NewLocaladmin(false)
		c.SetExists()
		err = rows.StructScan(c)
		if err != nil {
			return nil, err
		}
		if withEntity {
			c.Entity, err = ds.GetEntityById(c.IdEntity, true, restricted)
			if err != nil {
				return nil, err
			}
		}
		localadmin = append(localadmin, c)
	}

	return localadmin, err
}

func (ds *Datastore) GetLocaladminById(idEntity int64, withEntity, restricted bool) (*Localadmin, error) {
	filter := make(map[string]interface{})
	filter["id_entity"] = idEntity
	return ds.GetLocaladmin(filter, withEntity, restricted)
}

func (ds *Datastore) GetLocaladmin(filter map[string]interface{}, withEntity, restricted bool) (*Localadmin, error) {
	where, values := generators.GenerateAndSearchClause(filter)
	sql := ds.postgres.Rebind(`SELECT ` +
		localadminVisibility(restricted) +
		`FROM places4all.localadmin ` +
		where)

	a := ALocaladmin(false)
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
