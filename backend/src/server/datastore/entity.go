package datastore

import (
	"errors"
	"gopkg.in/guregu/null.v3/zero"
	"server/datastore/generators"
	"server/datastore/metadata"
	"time"
)

type Entity struct {
	Id            int64       `json:"id" db:"id"`
	IdCountry     int64       `json:"idCountry" db:"id_country"`
	Name          string      `json:"name" db:"name"`
	Email         string      `json:"email" db:"email"`
	Username      string      `json:"username" db:"username"`
	Password      string      `json:"-" db:"password"`
	Image         []byte      `json:"image" db:"image"`
	ImageMimetype zero.String `json:"-" db:"image_mimetype"`
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

	if e.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.entity (` +
		`id_country, name, email, username, password, image, image_mimetype, banned, banned_date, reason, mobilephone, telephone, created_date` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13` +
		`) RETURNING id`

	err := ds.postgres.QueryRow(sql, e.IdCountry, e.Name, e.Email, e.Username, e.Password, e.Image, e.ImageMimetype, e.Banned, e.BannedDate, e.Reason, e.Mobilephone, e.Telephone, e.CreatedDate).Scan(&e.Id)
	if err != nil {
		return err
	}

	e.SetExists()

	return err
}

func (ds *Datastore) UpdateEntity(e *Entity) error {
	var err error

	if !e.Exists() {
		return errors.New("update failed: does not exist")
	}

	if e.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.entity SET (` +
		`id_country, name, email, username, password, image, image_mimetype, banned, banned_date, reason, mobilephone, telephone, created_date` +
		`) = ( ` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13` +
		`) WHERE id = $14`

	_, err = ds.postgres.Exec(sql, e.IdCountry, e.Name, e.Email, e.Username, e.Password, e.Image, e.ImageMimetype, e.Banned, e.BannedDate, e.Reason, e.Mobilephone, e.Telephone, e.CreatedDate, e.Id)
	return err
}

func (ds *Datastore) SaveEntity(e *Entity) error {
	if e.Exists() {
		return ds.UpdateEntity(e)
	}

	return ds.InsertEntity(e)
}

func (ds *Datastore) UpsertEntity(p *Entity) error {

	if p.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.entity (` +
		`id, id_country, name, email, username, password, image, image_mimetype, banned, banned_date, reason, mobilephone, telephone, created_date` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13, $14` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_country, name, email, username, password, image, image_mimetype, banned, banned_date, reason, mobilephone, telephone, created_date` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_country, EXCLUDED.name, EXCLUDED.email, EXCLUDED.username, EXCLUDED.password, EXCLUDED.image, EXCLUDED.image_mimetype, EXCLUDED.banned, EXCLUDED.banned_date, EXCLUDED.reason, EXCLUDED.mobilephone, EXCLUDED.telephone, EXCLUDED.created_date` +
		`)`

	_, err := ds.postgres.Exec(sql, p.Id, p.IdCountry, p.Name, p.Email, p.Username, p.Password, p.Image, p.ImageMimetype, p.Banned, p.BannedDate, p.Reason, p.Mobilephone, p.Telephone, p.CreatedDate)
	if err != nil {
		return err
	}

	p.SetExists()

	return err
}

func (ds *Datastore) DeleteEntity(e *Entity) error {

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

func (ds *Datastore) GetEntityByEmail(email string) (*Entity, error) {

	const sql = `SELECT ` +
		`id, id_country, name, email, username, password, image, banned, banned_date, reason, mobilephone, telephone, created_date ` +
		`FROM places4all.entity ` +
		`WHERE email = $1`

	e := AEntity(false)
	e.SetExists()

	err := ds.postgres.QueryRowx(sql, email).StructScan(&e)
	if err != nil {
		return nil, err
	}

	return &e, err
}

func (ds *Datastore) GetEntityById(id int64) (*Entity, error) {

	const sql = `SELECT ` +
		`id, id_country, name, email, username, password, image, banned, banned_date, reason, mobilephone, telephone, created_date ` +
		`FROM places4all.entity ` +
		`WHERE id = $1`

	e := AEntity(false)
	e.SetExists()

	err := ds.postgres.QueryRowx(sql, id).StructScan(&e)
	if err != nil {
		return nil, err
	}

	return &e, err
}

func (ds *Datastore) GetEntityByUsername(username string) (*Entity, error) {
	var err error

	const sql = `SELECT ` +
		`id, id_country, name, email, username, password, image, banned, banned_date, reason, mobilephone, telephone, created_date ` +
		`FROM places4all.entity ` +
		`WHERE username = $1`

	e := AEntity(false)
	e.SetExists()

	err = ds.postgres.QueryRowx(sql, username).StructScan(&e)
	if err != nil {
		return nil, err
	}

	return &e, err
}
func (ds *Datastore) GetEntityByUsernamePassword(username string, password string) (*Entity, error) {
	var err error

	const sql = `SELECT ` +
		`id, id_country, name, email, username, password, image, banned, banned_date, reason, mobilephone, telephone, created_date ` +
		`FROM places4all.entity ` +
		`WHERE username = $1 and password = $2`

	e := AEntity(false)
	e.SetExists()

	err = ds.postgres.QueryRowx(sql, username, password).StructScan(&e)
	if err != nil {
		return nil, err
	}

	return &e, err
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

func (ds *Datastore) GetEntityWithForeign(filter map[string]interface{}) (*Entity, error) {
	where, values := generators.GenerateAndSearchClause(filter)
	sql := ds.postgres.Rebind(`SELECT ` +
		`entity.id, entity.id_country, entity.name, entity.email, entity.username, entity.password, entity.image, entity.image_mimetype, entity.banned, entity.banned_date, entity.reason, entity.mobilephone, entity.telephone, entity.created_date ` +
		`FROM places4all.entity ` +
		where)

	e := AEntity(false)

	err := ds.postgres.QueryRowx(sql, values...).StructScan(&e)
	if err != nil {
		return nil, err
	}

	e.Country, err = ds.GetCountryById(e.IdCountry)
	if err != nil {
		return nil, err
	}

	return &e, err
}

func (ds *Datastore) GetEntityByUsernamePassword(username string, password string) (*Entity, error) {
	var err error

	const sql = `SELECT ` +
		`id, id_country, name, email, username, password, image, banned, banned_date, reason, mobilephone, telephone, created_date ` +
		`FROM places4all.entity ` +
		`WHERE username = $1 and password = $2`

	p := AEntity(false)
	p.SetExists()

	err = ds.postgres.QueryRowx(sql, username,password).StructScan(&p)
	if err != nil {
		return nil, err
	}

	return &p, err
}

