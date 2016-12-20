package datastore

import (
	"errors"
	"gopkg.in/guregu/null.v3/zero"
	"server/datastore/metadata"
	"time"
)

type Entity struct {
	Id          int64       `json:"id" db:"id"`
	IdCountry   int64       `json:"idCountry" db:"id_country"`
	Name        string      `json:"name" db:"name"`
	Email       string      `json:"email" db:"email"`
	Username    string      `json:"username" db:"username"`
	Password    string      `json:"-" db:"password"`
	Image       []byte      `json:"image" db:"image"`
	Banned      zero.Bool   `json:"banned" db:"banned"`
	BannedDate  zero.Time   `json:"bannedDate" db:"banned_date"`
	Reason      zero.String `json:"reason" db:"reason"`
	Mobilephone zero.String `json:"mobilephone" db:"mobilephone"`
	Telephone   zero.String `json:"telephone" db:"telephone"`
	CreatedDate time.Time   `json:"createdDate" db:"created_date"`

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

func (ds *Datastore) InsertEntity(p *Entity) error {

	if p.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.entity (` +
		`id_country, name, email, username, password, image, banned, banned_date, reason, mobilephone, telephone, created_date` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12` +
		`) RETURNING id`

	err := ds.postgres.QueryRow(sql, p.IdCountry, p.Name, p.Email, p.Username, p.Password, p.Image, p.Banned, p.BannedDate, p.Reason, p.Mobilephone, p.Telephone, p.CreatedDate).Scan(&p.Id)
	if err != nil {
		return err
	}

	p.SetExists()

	return err
}

func (ds *Datastore) UpdateEntity(p *Entity) error {
	var err error

	if !p.Exists() {
		return errors.New("update failed: does not exist")
	}

	if p.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.entity SET (` +
		`id_country, name, email, username, password, image, banned, banned_date, reason, mobilephone, telephone, created_date` +
		`) = ( ` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12` +
		`) WHERE id = $13`

	_, err = ds.postgres.Exec(sql, p.IdCountry, p.Name, p.Email, p.Username, p.Password, p.Image, p.Banned, p.BannedDate, p.Reason, p.Mobilephone, p.Telephone, p.CreatedDate, p.Id)
	return err
}

func (ds *Datastore) SaveEntity(p *Entity) error {
	if p.Exists() {
		return ds.UpdateEntity(p)
	}

	return ds.InsertEntity(p)
}

func (ds *Datastore) UpsertEntity(p *Entity) error {

	if p.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.entity (` +
		`id, id_country, name, email, username, password, image, banned, banned_date, reason, mobilephone, telephone, created_date` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12, $13` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_country, name, email, username, password, image, banned, banned_date, reason, mobilephone, telephone, created_date` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_country, EXCLUDED.name, EXCLUDED.email, EXCLUDED.username, EXCLUDED.password, EXCLUDED.image, EXCLUDED.banned, EXCLUDED.banned_date, EXCLUDED.reason, EXCLUDED.mobilephone, EXCLUDED.telephone, EXCLUDED.created_date` +
		`)`

	_, err := ds.postgres.Exec(sql, p.Id, p.IdCountry, p.Name, p.Email, p.Username, p.Password, p.Image, p.Banned, p.BannedDate, p.Reason, p.Mobilephone, p.Telephone, p.CreatedDate)
	if err != nil {
		return err
	}

	p.SetExists()

	return err
}

func (ds *Datastore) DeleteEntity(p *Entity) error {

	if !p.Exists() {
		return nil
	}

	if p.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.entity WHERE id = $1`

	_, err := ds.postgres.Exec(sql, p.Id)
	if err != nil {
		return err
	}

	p.SetDeleted()

	return err
}

func (ds *Datastore) GetEntityCountry(p *Entity) (*Country, error) {
	return ds.GetCountryById(p.IdCountry)
}

func (ds *Datastore) GetEntityByEmail(email string) (*Entity, error) {

	const sql = `SELECT ` +
		`id, id_country, name, email, username, password, image, banned, banned_date, reason, mobilephone, telephone, created_date ` +
		`FROM places4all.entity ` +
		`WHERE email = $1`

	p := AEntity(false)
	p.SetExists()

	err := ds.postgres.QueryRowx(sql, email).StructScan(&p)
	if err != nil {
		return nil, err
	}

	return &p, err
}

func (ds *Datastore) GetEntityById(id int64) (*Entity, error) {

	const sql = `SELECT ` +
		`id, id_country, name, email, username, password, image, banned, banned_date, reason, mobilephone, telephone, created_date ` +
		`FROM places4all.entity ` +
		`WHERE id = $1`

	p := AEntity(false)
	p.SetExists()

	err := ds.postgres.QueryRowx(sql, id).Scan(&p)
	if err != nil {
		return nil, err
	}

	return &p, err
}

func (ds *Datastore) GetEntityByUsername(username string) (*Entity, error) {
	var err error

	const sql = `SELECT ` +
		`id, id_country, name, email, username, password, image, banned, banned_date, reason, mobilephone, telephone, created_date ` +
		`FROM places4all.entity ` +
		`WHERE username = $1`

	p := AEntity(false)
	p.SetExists()

	err = ds.postgres.QueryRowx(sql, username).StructScan(&p)
	if err != nil {
		return nil, err
	}

	return &p, err
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

func (ds *Datastore) CheckEntityUsername(username string, email string) (*Entity, error)  {
	var err error

	const sql = `SELECT ` +
		`id, id_country, name, email, username, password, image, banned, banned_date, reason, mobilephone, telephone, created_date ` +
		`FROM places4all.entity ` +
		`WHERE username = $1 and email=$2`

	p := AEntity(false)
	p.SetExists()

	err = ds.postgres.QueryRowx(sql, username).StructScan(&p)
	if err != nil {
		return nil, err
	}

	return &p, err
}