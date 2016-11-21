package datastore

import (
	"database/sql"
	"errors"
	"github.com/lib/pq"
	"server/datastore/metadata"
	"time"
)

type Person struct {
	Id          int64          `json:"id" db:"id"`
	IdCountry   int64          `json:"idCountry" db:"id_country"`
	Name        string         `json:"name" db:"name"`
	Email       string         `json:"email" db:"email"`
	Username    string         `json:"username" db:"username"`
	Password    string         `json:"password" db:"password"`
	ImageUrl    sql.NullString `json:"imageUrl" db:"image_url"`
	Banned      pq.NullTime    `json:"banned" db:"banned"`
	Reason      sql.NullString `json:"reason" db:"reason"`
	Mobilephone sql.NullString `json:"mobilephone" db:"mobilephone"`
	Telephone   sql.NullString `json:"telephone" db:"telephone"`
	Created     *time.Time     `json:"created" db:"created"`
	meta        metadata.Metadata
}

func (p *Person) SetExists() {
	p.meta.Exists = true
}

func (p *Person) SetDeleted() {
	p.meta.Deleted = true
}

func (p *Person) Exists() bool {
	return p.meta.Exists
}

func (p *Person) Deleted() bool {
	return p.meta.Deleted
}

func (ds *Datastore) InsertPerson(p *Person) error {
	var err error

	if p.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.person (` +
		`id_country, name, email, username, password, image_url, banned, reason, mobilephone, telephone, created` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11` +
		`) RETURNING id`

	err = ds.postgres.QueryRow(sql, p.IdCountry, p.Name, p.Email, p.Username, p.Password, p.ImageUrl, p.Banned, p.Reason, p.Mobilephone, p.Telephone, p.Created).Scan(&p.Id)
	if err != nil {
		return err
	}

	p.SetExists()

	return nil
}

func (ds *Datastore) UpdatePerson(p *Person) error {
	var err error

	if !p.Exists() {
		return errors.New("update failed: does not exist")
	}

	if p.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.person SET (` +
		`id_country, name, email, username, password, image_url, banned, reason, mobilephone, telephone, created` +
		`) = ( ` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11` +
		`) WHERE id = $12`

	_, err = ds.postgres.Exec(sql, p.IdCountry, p.Name, p.Email, p.Username, p.Password, p.ImageUrl, p.Banned, p.Reason, p.Mobilephone, p.Telephone, p.Created, p.Id)
	return err
}

func (ds *Datastore) SavePerson(p *Person) error {
	if p.Exists() {
		return ds.UpdatePerson(p)
	}

	return ds.InsertPerson(p)
}

func (ds *Datastore) UpsertPerson(p *Person) error {
	var err error

	if p.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.person (` +
		`id, id_country, name, email, username, password, image_url, banned, reason, mobilephone, telephone, created` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_country, name, email, username, password, image_url, banned, reason, mobilephone, telephone, created` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_country, EXCLUDED.name, EXCLUDED.email, EXCLUDED.username, EXCLUDED.password, EXCLUDED.image_url, EXCLUDED.banned, EXCLUDED.reason, EXCLUDED.mobilephone, EXCLUDED.telephone, EXCLUDED.created` +
		`)`

	_, err = ds.postgres.Exec(sql, p.Id, p.IdCountry, p.Name, p.Email, p.Username, p.Password, p.ImageUrl, p.Banned, p.Reason, p.Mobilephone, p.Telephone, p.Created)
	if err != nil {
		return err
	}

	p.SetExists()

	return nil
}

func (ds *Datastore) DeletePerson(p *Person) error {
	var err error

	if !p.Exists() {
		return nil
	}

	if p.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.person WHERE id = $1`

	_, err = ds.postgres.Exec(sql, p.Id)
	if err != nil {
		return err
	}

	p.SetDeleted()

	return nil
}

func (ds *Datastore) GetPersonCountry(p *Person) (*Country, error) {
	return ds.GetCountryById(p.IdCountry)
}

func (ds *Datastore) GetPersonByEmail(email string) (*Person, error) {
	var err error

	const sql = `SELECT ` +
		`id, id_country, name, email, username, password, image_url, banned, reason, mobilephone, telephone, created ` +
		`FROM places4all.person ` +
		`WHERE email = $1`

	p := Person{}
	p.SetExists()

	err = ds.postgres.QueryRow(sql, email).Scan(&p.Id, &p.IdCountry, &p.Name, &p.Email, &p.Username, &p.Password, &p.ImageUrl, &p.Banned, &p.Reason, &p.Mobilephone, &p.Telephone, &p.Created)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (ds *Datastore) GetPersonById(id int64) (*Person, error) {
	var err error

	const sql = `SELECT ` +
		`id, id_country, name, email, username, password, image_url, banned, reason, mobilephone, telephone, created ` +
		`FROM places4all.person ` +
		`WHERE id = $1`

	p := Person{}
	p.SetExists()

	err = ds.postgres.QueryRow(sql, id).Scan(&p.Id, &p.IdCountry, &p.Name, &p.Email, &p.Username, &p.Password, &p.ImageUrl, &p.Banned, &p.Reason, &p.Mobilephone, &p.Telephone, &p.Created)
	if err != nil {
		return nil, err
	}

	return &p, nil
}

func (ds *Datastore) GetPersonByUsername(username string) (*Person, error) {
	var err error

	const sql = `SELECT ` +
		`id, id_country, name, email, username, password, image_url, banned, reason, mobilephone, telephone, created ` +
		`FROM places4all.person ` +
		`WHERE username = $1`

	p := Person{}
	p.SetExists()

	err = ds.postgres.QueryRow(sql, username).Scan(&p.Id, &p.IdCountry, &p.Name, &p.Email, &p.Username, &p.Password, &p.ImageUrl, &p.Banned, &p.Reason, &p.Mobilephone, &p.Telephone, &p.Created)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
