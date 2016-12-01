package datastore

import (
	"errors"
	"gopkg.in/guregu/null.v3/zero"
	"server/datastore/metadata"
)

type Address struct {
	Id           int64       `json:"id" db:"id"`
	IdCountry    int64       `json:"idCountry" db:"id_country"`
	AddressLine1 string      `json:"addressLine1" db:"address_line1"`
	AddressLine2 zero.String `json:"addressLine2" db:"address_line2"`
	AddressLine3 zero.String `json:"addressLine3" db:"address_line3"`
	TownCity     zero.String `json:"townCity" db:"town_city"`
	County       zero.String `json:"county" db:"county"`
	Postcode     zero.String `json:"postcode" db:"postcode"`
	Latitude     zero.String `json:"latitude" db:"latitude"`
	Longitude    zero.String `json:"longitude" db:"longitude"`

	// Objects
	Country *Country `json:"country,omitempty"`

	meta metadata.Metadata
}

func (a *Address) SetExists() {
	a.meta.Exists = true
}

func (a *Address) SetDeleted() {
	a.meta.Deleted = true
}

func (a *Address) Exists() bool {
	return a.meta.Exists
}

func (a *Address) Deleted() bool {
	return a.meta.Deleted
}

func AAddress(allocateObjects bool) Address {
	address := Address{}
	if allocateObjects {
		address.Country = NewCountry(allocateObjects)
	}
	return address
}
func NewAddress(allocateObjects bool) *Address {
	address := AAddress(allocateObjects)
	return &address
}

func (ds *Datastore) InsertAddress(a *Address) error {

	if a.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.address (` +
		`id_country, address_line1, address_line2, address_line3, town_city, county, postcode, latitude, longitude` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9` +
		`) RETURNING id`

	err := ds.postgres.QueryRow(sql, a.IdCountry, a.AddressLine1, a.AddressLine2, a.AddressLine3, a.TownCity, a.County, a.Postcode, a.Latitude, a.Longitude).Scan(&a.Id)
	if err != nil {
		return err
	}

	a.SetExists()

	return err
}

func (ds *Datastore) UpdateAddress(a *Address) error {

	if !a.Exists() {
		return errors.New("update failed: does not exist")
	}

	if a.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.address SET (` +
		`id_country, address_line1, address_line2, address_line3, town_city, county, postcode, latitude, longitude` +
		`) = (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9` +
		`) WHERE id = $10`

	_, err := ds.postgres.Exec(sql, a.IdCountry, a.AddressLine1, a.AddressLine2, a.AddressLine3, a.TownCity, a.County, a.Postcode, a.Latitude, a.Longitude, a.Id)
	return err
}

func (ds *Datastore) SaveAddress(a *Address) error {
	if a.Exists() {
		return ds.UpdateAddress(a)
	}

	return ds.InsertAddress(a)
}

func (ds *Datastore) UpsertAddress(a *Address) error {

	if a.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.address (` +
		`id, id_country, address_line1, address_line2, address_line3, town_city, county, postcode, latitude, longitude` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9, $10` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_country, address_line1, address_line2, address_line3, town_city, county, postcode, latitude, longitude` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_country, EXCLUDED.address_line1, EXCLUDED.address_line2, EXCLUDED.address_line3, EXCLUDED.town_city, EXCLUDED.county, EXCLUDED.postcode, EXCLUDED.latitude, EXCLUDED.longitude` +
		`)`

	_, err := ds.postgres.Exec(sql, a.Id, a.IdCountry, a.AddressLine1, a.AddressLine2, a.AddressLine3, a.TownCity, a.County, a.Postcode, a.Latitude, a.Longitude)
	if err != nil {
		return err
	}

	a.SetExists()

	return err
}

func (ds *Datastore) DeleteAddress(a *Address) error {

	if !a.Exists() {
		return nil
	}

	if a.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.address WHERE id = $1`

	_, err := ds.postgres.Exec(sql, a.Id)
	if err != nil {
		return err
	}

	a.SetDeleted()

	return err
}

func (ds *Datastore) GetAddressCountry(a *Address) (*Country, error) {
	return ds.GetCountryById(a.IdCountry)
}

func (ds *Datastore) GetAddressById(id int64) (*Address, error) {

	const sql = `SELECT
		address.id, address.id_country,
		address.address_line1, address.address_line2, address.address_line3,
		address.town_city, address.county, address.postcode,
		address.latitude, address.longitude,
		country.id, country.name, country.iso2
		FROM places4all.address
		JOIN places4all.country ON country.id = address.id_country
		WHERE address.id = $1`

	a := NewAddress(true)
	a.SetExists()

	err := ds.postgres.QueryRow(sql, id).Scan(
		&a.Id, &a.IdCountry,
		&a.AddressLine1, &a.AddressLine2, &a.AddressLine3,
		&a.TownCity, &a.County, &a.Postcode,
		&a.Latitude, &a.Longitude,
		&a.Country.Id, &a.Country.Name, &a.Country.Iso2)
	if err != nil {
		return nil, err
	}

	return a, err
}
