package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type Country struct {
	Id   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
	Iso2 string `json:"iso2" db:"iso2"`

	meta metadata.Metadata
}

func (c *Country) SetExists() {
	c.meta.Exists = true
}

func (c *Country) SetDeleted() {
	c.meta.Deleted = true
}

func (c *Country) Exists() bool {
	return c.meta.Exists
}

func (c *Country) Deleted() bool {
	return c.meta.Deleted
}

func ACountry(allocateObjects bool) Country {
	return Country{}
}

func NewCountry(allocateObjects bool) *Country {
	country := ACountry(allocateObjects)
	return &country
}

func (ds *Datastore) InsertCountry(c *Country) error {

	if c.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.country (` +
		`name, iso2` +
		`) VALUES (` +
		`$1, $2` +
		`) RETURNING id`

	err := ds.postgres.QueryRow(sql, c.Name, c.Iso2).Scan(&c.Id)
	if err != nil {
		return err
	}

	c.SetExists()

	return err
}

func (ds *Datastore) UpdateCountry(c *Country) error {

	if !c.Exists() {
		return errors.New("update failed: does not exist")
	}

	if c.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.country SET (` +
		`name, iso2` +
		`) = ( ` +
		`$1, $2` +
		`) WHERE id = $3`

	_, err := ds.postgres.Exec(sql, c.Name, c.Iso2, c.Id)
	return err
}

func (ds *Datastore) SaveCountry(c *Country) error {
	if c.Exists() {
		return ds.UpdateCountry(c)
	}

	return ds.InsertCountry(c)
}

func (ds *Datastore) UpsertCountry(c *Country) error {

	if c.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.country (` +
		`id, name, iso2` +
		`) VALUES (` +
		`$1, $2, $3` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, name, iso2` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.name, EXCLUDED.iso2` +
		`)`

	_, err := ds.postgres.Exec(sql, c.Id, c.Name, c.Iso2)
	if err != nil {
		return err
	}

	c.SetExists()

	return err
}

func (ds *Datastore) DeleteCountry(c *Country) error {

	if !c.Exists() {
		return nil
	}

	if c.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.country WHERE id = $1`

	_, err := ds.postgres.Exec(sql, c.Id)
	if err != nil {
		return err
	}

	c.SetDeleted()

	return err
}

func (ds *Datastore) GetCountryByName(name string) (*Country, error) {

	const sql = `SELECT ` +
		`id, name, iso2 ` +
		`FROM places4all.country ` +
		`WHERE name = $1`

	c := Country{}
	c.SetExists()

	err := ds.postgres.QueryRow(sql, name).Scan(&c.Id, &c.Name, &c.Iso2)
	if err != nil {
		return nil, err
	}

	return &c, err
}

func (ds *Datastore) GetCountryById(id int64) (*Country, error) {

	const sql = `SELECT ` +
		`id, name, iso2 ` +
		`FROM places4all.country ` +
		`WHERE id = $1`

	c := ACountry(true)
	c.SetExists()

	err := ds.postgres.QueryRow(sql, id).Scan(&c.Id, &c.Name, &c.Iso2)
	if err != nil {
		return nil, err
	}

	return &c, err
}

func (ds *Datastore) GetCountries() ([]*Country, error) {

	const sql = `SELECT ` +
		`id, name, iso2 ` +
		`FROM places4all.country`

	rows, err := ds.postgres.Queryx(sql)
	if err != nil {
		return nil, err
	}

	countries := make([]*Country, 0)
	for rows.Next() {
		c := NewCountry(true)
		c.SetExists()
		err = rows.StructScan(c)
		if err != nil {
			return nil, err
		}
		countries = append(countries, c)
	}

	return countries, err
}
