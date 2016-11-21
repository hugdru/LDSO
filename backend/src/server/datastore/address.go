package datastore

import (
	"database/sql"
	"errors"
	"server/datastore/metadata"
)

type Address struct {
	Id           int64           `json:"id" db:"id"`
	IdCountry    int64           `json:"idCountry" db:"id_country"`
	AddressLine1 string          `json:"addressLine1" db:"address_line1"`
	AddressLine2 sql.NullString  `json:"addressLine2" db:"address_line2"`
	AddressLine3 sql.NullString  `json:"addressLine3" db:"address_line3"`
	TownCity     sql.NullString  `json:"townCity" db:"town_city"`
	County       sql.NullString  `json:"county" db:"county"`
	Postcode     sql.NullString  `json:"postcode" db:"postcode"`
	Latitude     sql.NullFloat64 `json:"latitude" db:"latitude"`
	Longitude    sql.NullFloat64 `json:"longitude" db:"longitude"`
	meta         metadata.Metadata
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

func (ds *Datastore) InsertAddress(a *Address) error {
	var err error

	if a.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.address (` +
		`id_country, address_line1, address_line2, address_line3, town_city, county, postcode, latitude, longitude` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9` +
		`) RETURNING id`

	err = ds.postgres.QueryRow(sql, a.IdCountry, a.AddressLine1, a.AddressLine2, a.AddressLine3, a.TownCity, a.County, a.Postcode, a.Latitude, a.Longitude).Scan(&a.Id)
	if err != nil {
		return err
	}

	a.SetExists()

	return nil
}

func (ds *Datastore) UpdateAddress(a *Address) error {
	var err error

	if !a.Exists() {
		return errors.New("update failed: does not exist")
	}

	if a.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.address SET (` +
		`id_country, address_line1, address_line2, address_line3, town_city, county, postcode, latitude, longitude` +
		`) = ( ` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9` +
		`) WHERE id = $10`

	_, err = ds.postgres.Exec(sql, a.IdCountry, a.AddressLine1, a.AddressLine2, a.AddressLine3, a.TownCity, a.County, a.Postcode, a.Latitude, a.Longitude, a.Id)
	return err
}

func (ds *Datastore) SaveAddress(a *Address) error {
	if a.Exists() {
		return ds.UpdateAddress(a)
	}

	return ds.InsertAddress(a)
}

func (ds *Datastore) UpsertAddress(a *Address) error {
	var err error

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

	_, err = ds.postgres.Exec(sql, a.Id, a.IdCountry, a.AddressLine1, a.AddressLine2, a.AddressLine3, a.TownCity, a.County, a.Postcode, a.Latitude, a.Longitude)
	if err != nil {
		return err
	}

	a.SetExists()

	return nil
}

func (ds *Datastore) DeleteAddress(a *Address) error {
	var err error

	if !a.Exists() {
		return nil
	}

	if a.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.address WHERE id = $1`

	_, err = ds.postgres.Exec(sql, a.Id)
	if err != nil {
		return err
	}

	a.SetDeleted()

	return nil
}

func (ds *Datastore) GetAddressCountry(a *Address) (*Country, error) {
	return ds.GetCountryById(a.IdCountry)
}

func (ds *Datastore) GetAddressById(id int64) (*Address, error) {
	var err error

	const sql = `SELECT ` +
		`id, id_country, address_line1, address_line2, address_line3, town_city, county, postcode, latitude, longitude ` +
		`FROM places4all.address ` +
		`WHERE id = $1`

	a := Address{}
	a.SetExists()

	err = ds.postgres.QueryRow(sql, id).Scan(&a.Id, &a.IdCountry, &a.AddressLine1, &a.AddressLine2, &a.AddressLine3, &a.TownCity, &a.County, &a.Postcode, &a.Latitude, &a.Longitude)
	if err != nil {
		return nil, err
	}

	return &a, nil
}
