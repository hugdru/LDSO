package datastore

import (
	"errors"
	"server/datastore/metadata"
	"strconv"

)

type Auditor struct {
	Id       int64 `json:"id" db:"id"`
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
		`) RETURNING id`

	err := ds.postgres.QueryRow(sql, a.IdEntity).Scan(&a.Id)
	if err != nil {
		return err
	}

	a.SetExists()

	return err
}

func (ds *Datastore) UpdateAuditor(a *Auditor) error {

	if !a.Exists() {
		return errors.New("update failed: does not exist")
	}

	if a.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.auditor SET (` +
		`id_entity` +
		`) = ( ` +
		`$1` +
		`) WHERE id = $2`

	_, err := ds.postgres.Exec(sql, a.IdEntity, a.Id)
	return err
}

func (ds *Datastore) SaveAuditor(a *Auditor) error {
	if a.Exists() {
		return ds.UpdateAuditor(a)
	}

	return ds.InsertAuditor(a)
}

func (ds *Datastore) UpsertAuditor(a *Auditor) error {

	if a.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.auditor (` +
		`id, id_entity` +
		`) VALUES (` +
		`$1, $2` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_entity` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_entity` +
		`)`

	_, err := ds.postgres.Exec(sql, a.Id, a.IdEntity)
	if err != nil {
		return err
	}

	a.SetExists()

	return err
}

func (ds *Datastore) DeleteAuditor(a *Auditor) error {

	if !a.Exists() {
		return nil
	}

	if a.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.auditor WHERE id = $1`

	_, err := ds.postgres.Exec(sql, a.Id)
	if err != nil {
		return err
	}

	a.SetDeleted()

	return err
}

func (ds *Datastore) GetAuditorEntity(a *Auditor) (*Entity, error) {
	return ds.GetEntityById(a.IdEntity)
}

func (ds *Datastore) GetAuditorById(id int64) (*Auditor, error) {
	var err error


	const sql = `SELECT
		auditor.id, auditor.id_entity,
		entity.id, entity.id_country, entity.name, entity.email,
		entity.username, entity.password, entity.image, entity.banned, entity.banned_date,
		entity.reason, entity.mobilephone, entity.telephone, entity.created_date,
		country.id, country.name, country.iso2
		FROM places4all.auditor
		JOIN places4all.entity  on entity.id   = auditor.id_entity
		JOIN places4all.country on country.id =  entity.id_country
		WHERE auditor.id = $1 `

	a := AAuditor(true)
	a.SetExists()

	err = ds.postgres.QueryRow(sql, id).Scan(
		&a.Id, &a.IdEntity,
		&a.Entity.Id, &a.Entity.IdCountry, &a.Entity.Name, &a.Entity.Email,
		&a.Entity.Username, &a.Entity.Password, &a.Entity.Image, &a.Entity.Banned, &a.Entity.BannedDate,
		&a.Entity.Reason, &a.Entity.Mobilephone, &a.Entity.Telephone, &a.Entity.CreatedDate,
		&a.Entity.Country.Id, &a.Entity.Country.Name, &a.Entity.Country.Iso2,
	)

	if err != nil {
		return nil, err
	}

	return &a, err
}
func (ds *Datastore) GetAuditors(limit, offset int) ([]*Auditor, error) {
//nao retor a info da entity
	rows, err := ds.postgres.Queryx(`SELECT ` +
		`auditor.id, auditor.id_entity `+ //, entity.email, entity.username, entity.password, entity.created_date ` +
		`FROM places4all.auditor `+
		`JOIN places4all.entity on entity.id = auditor.id_entity `+
//		`JOIN places4all.country on country.id = entity.id_country `+
		`ORDER BY auditor.id DESC LIMIT ` + strconv.Itoa(limit) +
		`OFFSET ` + strconv.Itoa(offset))

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