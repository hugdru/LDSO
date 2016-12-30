package datastore

import (
	"database/sql"
	"errors"
	"server/datastore/generators"
	"server/datastore/metadata"
	"strconv"
)

func superadminVisibility(restricted bool) string {
	const superadminRestricted = "superadmin.id_entity "
	const superadminAll = "superadmin.id_entity "
	if restricted {
		return superadminRestricted
	}
	return superadminAll
}

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

func (s *Superadmin) MustSet(idEntity int64) error {
	if idEntity != 0 {
		s.IdEntity = idEntity
	} else {
		return errors.New("idEntity must be set")
	}

	return nil
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
	return ds.InsertSuperadminTx(nil, s)
}

func (ds *Datastore) InsertSuperadminTx(tx *sql.Tx, s *Superadmin) error {

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

	var err error

	if tx != nil {
		_, err = tx.Exec(sql, s.IdEntity)

	} else {
		_, err = ds.postgres.Exec(sql, s.IdEntity)
	}
	if err != nil {
		return err
	}

	s.SetExists()

	return err
}
func (ds *Datastore) UpdateSuperadmin(s *Superadmin) error {
	return ds.UpdateSuperadminTx(nil, s)
}

func (ds *Datastore) UpdateSuperadminTx(tx *sql.Tx, s *Superadmin) error {
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
	return ds.SaveSuperadminTx(nil, s)
}

func (ds *Datastore) SaveSuperadminTx(tx *sql.Tx, s *Superadmin) error {

	if s == nil {
		return errors.New("superadmin should not be nil")
	}

	if s.Exists() {
		return ds.UpdateSuperadminTx(tx, s)
	}

	return ds.InsertSuperadminTx(tx, s)
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

func (ds *Datastore) GetSuperadminEntity(s *Superadmin, withCountry, restricted bool) (*Entity, error) {
	return ds.GetEntityById(s.IdEntity, withCountry, restricted)
}

func (ds *Datastore) GetSuperadmins(limit, offset int, withEntity, restricted bool) ([]*Superadmin, error) {
	rows, err := ds.postgres.Queryx(ds.postgres.Rebind(`SELECT ` +
		superadminVisibility(restricted) +
		`FROM places4all.superadmin ` +
		`ORDER BY superadmin.id_entity DESC LIMIT ` + strconv.Itoa(limit) +
		`OFFSET ` + strconv.Itoa(offset)))

	if err != nil {
		return nil, err
	}

	superadmin := make([]*Superadmin, 0)
	for rows.Next() {
		c := NewSuperadmin(false)
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
		superadmin = append(superadmin, c)
	}

	return superadmin, err
}

func (ds *Datastore) GetSuperadminById(idEntity int64, withEntity, restricted bool) (*Superadmin, error) {
	filter := make(map[string]interface{})
	filter["id_entity"] = idEntity
	return ds.GetSuperadmin(filter, withEntity, restricted)
}

func (ds *Datastore) GetSuperadmin(filter map[string]interface{}, withEntity, restricted bool) (*Superadmin, error) {
	where, values := generators.GenerateAndSearchClause(filter)
	sql := ds.postgres.Rebind(`SELECT ` +
		superadminVisibility(restricted) +
		`FROM places4all.superadmin ` +
		where)

	a := ASuperadmin(false)
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
