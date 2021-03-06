package datastore

import (
	"database/sql"
	"errors"
	"gopkg.in/guregu/null.v3/zero"
	"server/datastore/generators"
	"server/datastore/metadata"
	"time"
)

func entityVisibility(restricted bool) string {
	const auditorRestricted = "entity.id, entity.id_country, entity.name, entity.username, entity.image, entity.image_mimetype, entity.image_hash, entity.created_date "
	const auditorAll = "entity.id, entity.id_country, entity.name, entity.email, entity.username, entity.password, entity.image, entity.image_mimetype, entity.image_hash, entity.banned, entity.banned_date, entity.reason, entity.mobilephone, entity.telephone, entity.created_date "
	if restricted {
		return auditorRestricted
	}
	return auditorAll
}

type Entity struct {
	Id            int64       `json:"id" db:"id"`
	IdCountry     int64       `json:"idCountry" db:"id_country"`
	Name          string      `json:"name" db:"name"`
	Email         string      `json:"email" db:"email"`
	Username      string      `json:"username" db:"username"`
	Password      string      `json:"-" db:"password"`
	Image         []byte      `json:"-" db:"image"`
	ImageMimetype zero.String `json:"-" db:"image_mimetype"`
	ImageHash     zero.String `json:"imageLocation" db:"image_hash"`
	Banned        zero.Bool   `json:"banned" db:"banned"`
	BannedDate    zero.Time   `json:"bannedDate" db:"banned_date"`
	Reason        zero.String `json:"reason" db:"reason"`
	Mobilephone   zero.String `json:"mobilephone" db:"mobilephone"`
	Telephone     zero.String `json:"telephone" db:"telephone"`
	CreatedDate   time.Time   `json:"createdDate" db:"created_date"`

	// Objects
	Country *Country `json:"country,omitempty"`

	meta metadata.Metadata
}

func AEntity(allocateObjects bool) Entity {
	entity := Entity{}
	if allocateObjects {
		entity.Country = NewCountry(allocateObjects)
	}
	return entity
}

func NewEntity(allocateObjects bool) *Entity {
	entity := AEntity(allocateObjects)
	return &entity
}

func (p *Entity) SetExists() {
	p.meta.Exists = true
}

func (p *Entity) SetDeleted() {
	p.meta.Deleted = true
}

func (p *Entity) Exists() bool {
	return p.meta.Exists
}

func (p *Entity) Deleted() bool {
	return p.meta.Deleted
}

func (ds *Datastore) InsertEntity(e *Entity) error {
	return ds.InsertEntityTx(nil, e)
}

func (ds *Datastore) InsertEntityTx(tx *sql.Tx, e *Entity) error {

	if e == nil {
		return errors.New("entity should not be nil")
	}

	if e.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.entity (` +
		`id_country, name, email, username, password, image, image_mimetype, image_hash, banned, banned_date, reason, mobilephone, telephone, created_date` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14` +
		`) RETURNING id`

	var err error
	if tx != nil {
		err = tx.QueryRow(sql, e.IdCountry, e.Name, e.Email, e.Username, e.Password, e.Image, e.ImageMimetype, e.ImageHash, e.Banned, e.BannedDate, e.Reason, e.Mobilephone, e.Telephone, e.CreatedDate).Scan(&e.Id)
	} else {
		err = ds.postgres.QueryRow(sql, e.IdCountry, e.Name, e.Email, e.Username, e.Password, e.Image, e.ImageMimetype, e.ImageHash, e.Banned, e.BannedDate, e.Reason, e.Mobilephone, e.Telephone, e.CreatedDate).Scan(&e.Id)
	}
	if err != nil {
		return err
	}

	e.SetExists()

	return err
}

func (ds *Datastore) UpdateEntity(e *Entity) error {
	return ds.UpdateEntityTx(nil, e)
}

func (ds *Datastore) UpdateEntityTx(tx *sql.Tx, e *Entity) error {

	if e == nil {
		return errors.New("entity should not be nil")
	}

	if !e.Exists() {
		return errors.New("update failed: does not exist")
	}

	if e.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.entity SET (` +
		`id_country, name, email, username, password, image, image_mimetype, image_hash, banned, banned_date, reason, mobilephone, telephone, created_date` +
		`) = ( ` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14` +
		`) WHERE id = $15`

	var err error
	if tx != nil {
		_, err = tx.Exec(sql, e.IdCountry, e.Name, e.Email, e.Username, e.Password, e.Image, e.ImageMimetype, e.ImageHash, e.Banned, e.BannedDate, e.Reason, e.Mobilephone, e.Telephone, e.CreatedDate, e.Id)

	} else {
		_, err = ds.postgres.Exec(sql, e.IdCountry, e.Name, e.Email, e.Username, e.Password, e.Image, e.ImageMimetype, e.ImageHash, e.Banned, e.BannedDate, e.Reason, e.Mobilephone, e.Telephone, e.CreatedDate, e.Id)
	}
	return err
}

func (ds *Datastore) SaveEntity(e *Entity) error {
	return ds.SaveEntityTx(nil, e)
}

func (ds *Datastore) SaveEntityTx(tx *sql.Tx, e *Entity) error {

	if e == nil {
		return errors.New("entity should not be nil")
	}

	if e.Exists() {
		return ds.UpdateEntityTx(tx, e)
	}

	return ds.InsertEntityTx(tx, e)
}

func (ds *Datastore) UpsertEntity(e *Entity) error {

	if e == nil {
		return errors.New("entity should not be nil")
	}

	if e.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.entity (` +
		`id, id_country, name, email, username, password, image, image_mimetype, image_hash, banned, banned_date, reason, mobilephone, telephone, created_date` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14, $15` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_country, name, email, username, password, image, image_mimetype, image_hash, banned, banned_date, reason, mobilephone, telephone, created_date` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_country, EXCLUDED.name, EXCLUDED.email, EXCLUDED.username, EXCLUDED.password, EXCLUDED.image, EXCLUDED.image_mimetype, EXCLUDED.image_hash, EXCLUDED.banned, EXCLUDED.banned_date, EXCLUDED.reason, EXCLUDED.mobilephone, EXCLUDED.telephone, EXCLUDED.created_date` +
		`)`

	_, err := ds.postgres.Exec(sql, e.Id, e.IdCountry, e.Name, e.Email, e.Username, e.Password, e.Image, e.ImageMimetype, e.ImageHash, e.Banned, e.BannedDate, e.Reason, e.Mobilephone, e.Telephone, e.CreatedDate)
	if err != nil {
		return err
	}

	e.SetExists()

	return err
}

func (ds *Datastore) DeleteEntity(e *Entity) error {

	if e == nil {
		return errors.New("entity should not be nil")
	}

	if !e.Exists() {
		return nil
	}

	if e.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.entity WHERE id = $1`

	_, err := ds.postgres.Exec(sql, e.Id)
	if err != nil {
		return err
	}

	e.SetDeleted()

	return err
}

func (ds *Datastore) GetEntityCountry(e *Entity) (*Country, error) {
	return ds.GetCountryById(e.IdCountry)
}

func (ds *Datastore) CheckEntityExists(filter map[string]interface{}) error {

	where, values := generators.GenerateOrSearchClause(filter)
	query := ds.postgres.Rebind(`SELECT id ` +
		`FROM places4all.entity ` +
		where)

	e := AEntity(false)
	e.SetExists()

	var id int64
	return ds.postgres.QueryRow(query, values...).Scan(&id)
}

func (ds *Datastore) GetEntityByEmail(email string, withCountry, restricted bool) (*Entity, error) {
	filter := make(map[string]interface{})
	filter["email"] = email
	return ds.GetEntity(filter, withCountry, restricted)
}

func (ds *Datastore) GetEntityByUsername(username string, withCountry, restricted bool) (*Entity, error) {
	filter := make(map[string]interface{})
	filter["username"] = username
	return ds.GetEntity(filter, withCountry, restricted)
}
func (ds *Datastore) GetEntityByUsernamePassword(username, password string, withCountry, restricted bool) (*Entity, error) {
	filter := make(map[string]interface{})
	filter["username"] = username
	filter["password"] = password
	return ds.GetEntity(filter, withCountry, restricted)
}

func (ds *Datastore) GetEntityById(id int64, withCountry, restricted bool) (*Entity, error) {
	filter := make(map[string]interface{})
	filter["id"] = id
	return ds.GetEntity(filter, withCountry, restricted)
}

func (ds *Datastore) GetEntity(filter map[string]interface{}, withCountry, restricted bool) (*Entity, error) {
	where, values := generators.GenerateAndSearchClause(filter)
	sql := ds.postgres.Rebind(`SELECT ` +
		entityVisibility(restricted) +
		`FROM places4all.entity ` +
		where)

	e := AEntity(false)
	e.SetExists()

	err := ds.postgres.QueryRowx(sql, values...).StructScan(&e)
	if err != nil {
		return nil, err
	}

	if withCountry {
		e.Country, err = ds.GetCountryById(e.IdCountry)
		if err != nil {
			return nil, err
		}
	}

	return &e, err
}
