package datastore

import (
	"database/sql"
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

func (a *Address) MustSet(idCountry int64, addressLine1 string) error {

	if idCountry != 0 {
		a.IdCountry = idCountry
	} else {
		return errors.New("idCountry must be set")
	}
	if addressLine1 != "" {
		a.AddressLine1 = addressLine1
	} else {
		return errors.New("addressLine1 must be set")
	}

	return nil
}

func (a *Address) AllSetIfNotEmptyOrNil(idCountry int64, addressLine1 string,
	addressLine2 string, addressLine3 string, townCity string,
	county string, postcode string, latitude string, longitude string) error {
	if idCountry != 0 {
		a.IdCountry = idCountry
	}

	if addressLine1 != "" {
		a.AddressLine1 = addressLine1
	}

	return a.OptionalSetIfNotEmptyOrNil(addressLine2, addressLine3, townCity,
		county, postcode, latitude, longitude)
}

func (a *Address) OptionalSetIfNotEmptyOrNil(
	addressLine2 string, addressLine3 string, townCity string,
	county string, postcode string, latitude string, longitude string) error {

	if addressLine2 != "" {
		a.AddressLine2 = zero.StringFrom(addressLine2)
	}

	if addressLine3 != "" {
		a.AddressLine3 = zero.StringFrom(addressLine3)
	}

	if townCity != "" {
		a.TownCity = zero.StringFrom(townCity)
	}

	if county != "" {
		a.County = zero.StringFrom(county)
	}

	if postcode != "" {
		a.Postcode = zero.StringFrom(postcode)
	}

	if latitude != "" {
		a.Latitude = zero.StringFrom(latitude)
	}

	if longitude != "" {
		a.Longitude = zero.StringFrom(longitude)
	}

	return nil
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

func (ds *Datastore) InsertAddressTx(tx *sql.Tx, a *Address) error {

	if a == nil {
		return errors.New("address should not be nil")
	}

	if a.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.address (` +
		`id_country, address_line1, address_line2, address_line3, town_city, county, postcode, latitude, longitude` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6, $7, $8, $9` +
		`) RETURNING id`

	var err error
	if tx != nil {
		err = tx.QueryRow(sql, a.IdCountry, a.AddressLine1, a.AddressLine2, a.AddressLine3, a.TownCity, a.County, a.Postcode, a.Latitude, a.Longitude).Scan(&a.Id)
	} else {
		err = ds.postgres.QueryRow(sql, a.IdCountry, a.AddressLine1, a.AddressLine2, a.AddressLine3, a.TownCity, a.County, a.Postcode, a.Latitude, a.Longitude).Scan(&a.Id)
	}
	if err != nil {
		return err
	}

	a.SetExists()

	return err
}

func (ds *Datastore) InsertAddress(a *Address) error {
	return ds.InsertAddressTx(nil, a)
}

func (ds *Datastore) UpdateAddress(a *Address) error {
	return ds.UpdateAddressTx(nil, a)
}

func (ds *Datastore) UpdateAddressTx(tx *sql.Tx, a *Address) error {
	if a == nil {
		return errors.New("address should not be nil")
	}

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

	var err error
	if tx != nil {
		_, err = tx.Exec(sql, a.IdCountry, a.AddressLine1, a.AddressLine2, a.AddressLine3, a.TownCity, a.County, a.Postcode, a.Latitude, a.Longitude, a.Id)

	} else {
		_, err = ds.postgres.Exec(sql, a.IdCountry, a.AddressLine1, a.AddressLine2, a.AddressLine3, a.TownCity, a.County, a.Postcode, a.Latitude, a.Longitude, a.Id)
	}
	return err
}

func (ds *Datastore) SaveAddress(a *Address) error {

	if a == nil {
		return errors.New("address should not be nil")
	}

	if a.Exists() {
		return ds.UpdateAddress(a)
	}

	return ds.InsertAddress(a)
}

func (ds *Datastore) UpsertAddress(a *Address) error {

	if a == nil {
		return errors.New("address should not be nil")
	}

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

func (ds *Datastore) DeleteAddressTx(tx *sql.Tx, a *Address) error {

	if a == nil {
		return errors.New("address should not be nil")
	}

	if !a.Exists() {
		return nil
	}

	if a.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.address WHERE id = $1`

	var err error
	if tx != nil {
		_, err = tx.Exec(sql, a.Id)

	} else {
		_, err = ds.postgres.Exec(sql, a.Id)

	}
	if err != nil {
		return err
	}

	a.SetDeleted()

	return err
}

func (ds *Datastore) DeleteAddress(a *Address) error {
	return ds.DeleteAddressTx(nil, a)
}

func (ds *Datastore) GetAddressCountry(a *Address) (*Country, error) {
	return ds.GetCountryById(a.IdCountry)
}

func (ds *Datastore) GetAddressByIdWithCountry(id int64) (*Address, error) {

	sql := ds.postgres.Rebind(`SELECT ` +
		`id, id_country, address_line1, address_line2, address_line3, town_city, county, postcode, latitude, longitude ` +
		`FROM places4all.address ` +
		`WHERE address.id = ?`)

	a := NewAddress(false)
	a.SetExists()

	err := ds.postgres.QueryRowx(sql, id).StructScan(a)
	if err != nil {
		return nil, err
	}

	a.Country, err = ds.GetCountryById(a.IdCountry)
	if err != nil {
		return nil, err
	}

	return a, err
}
