package datastore

import (
	"errors"
	"server/datastore/generators"
	"server/datastore/metadata"
	"strconv"
	"time"
)

type Property struct {
	Id          int64     `json:"id" db:"id"`
	IdAddress   int64     `json:"idAddress" db:"id_address"`
	Name        string    `json:"name" db:"name"`
	Details     string    `json:"details" db:"details"`
	CreatedDate time.Time `json:"createdDate" db:"created_date"`

	// Objects
	Address *Address  `json:"address,omitempty"`
	Owners  []*Client `json:"owners,omitempty"`
	Tags    []*Tag    `json:"tags,omitempty"`

	meta metadata.Metadata
}

func (p *Property) SetExists() {
	p.meta.Exists = true
}

func (p *Property) SetDeleted() {
	p.meta.Deleted = true
}

func (p *Property) Exists() bool {
	return p.meta.Exists
}

func (p *Property) Deleted() bool {
	return p.meta.Deleted
}

func AProperty(allocateObjects bool) Property {
	property := Property{}
	if allocateObjects {
		property.Address = NewAddress(allocateObjects)
		property.Owners = make([]*Client, 0)
		property.Tags = make([]*Tag, 0)
	}
	return property
}

func NewProperty(allocateObjects bool) *Property {
	property := AProperty(allocateObjects)
	return &property
}

func (ds *Datastore) InsertProperty(p *Property) error {

	if p.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.property (` +
		`id_address, name, details, created_date` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`) RETURNING id`

	err := ds.postgres.QueryRow(sql, p.IdAddress, p.Name, p.Details, p.CreatedDate).Scan(&p.Id)
	if err != nil {
		return err
	}

	p.SetExists()

	return err
}

func (ds *Datastore) UpdateProperty(p *Property) error {

	if !p.Exists() {
		return errors.New("update failed: does not exist")
	}

	if p.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.property SET (` +
		`id_address, name, details, created_date` +
		`) = ( ` +
		`$1, $2, $3, $4` +
		`) WHERE id = $5`

	_, err := ds.postgres.Exec(sql, p.IdAddress, p.Name, p.Details, p.CreatedDate, p.Id)
	return err
}

func (ds *Datastore) SaveProperty(p *Property) error {
	if p.Exists() {
		return ds.UpdateProperty(p)
	}

	return ds.InsertProperty(p)
}

func (ds *Datastore) UpsertProperty(p *Property) error {

	if p.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.property (` +
		`id, id_address, name, details, created` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_address, name, details, created` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_address, EXCLUDED.name, EXCLUDED.details, EXCLUDED.created` +
		`)`

	_, err := ds.postgres.Exec(sql, p.Id, p.IdAddress, p.Name, p.Details, p.CreatedDate)
	if err != nil {
		return err
	}

	p.SetExists()

	return err
}

func (ds *Datastore) DeleteProperty(p *Property) error {

	if !p.Exists() {
		return nil
	}

	if p.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.property WHERE id = $1`

	_, err := ds.postgres.Exec(sql, p.Id)
	if err != nil {
		return err
	}

	p.SetDeleted()

	return err
}

func (ds *Datastore) GetPropertyAddress(p *Property) (*Address, error) {
	return ds.GetAddressByIdWithForeign(p.IdAddress)
}

func (ds *Datastore) GetPropertyById(id int64) (*Property, error) {

	p := NewProperty(false)
	err := ds.postgres.QueryRowx(`SELECT `+
		`p.id, p.id_address, p.name, p.details, p.created_date ` +
		`FROM places4all.property as p `+
		`WHERE p.id = $1`,
		id).StructScan(p)
	if err != nil {
		return nil, err
	}

	p.Address, err = ds.GetAddressByIdWithForeign(p.IdAddress)
	if err != nil {
		return nil, err
	}
	p.Tags, err = ds.GetPropertyTagByIdProperty(p.Id)
	if err != nil {
		return nil, err
	}
	p.Owners, err = ds.GetPropertyClientsByIdProperty(p.Id)
	if err != nil {
		return nil, err
	}

	return p, err
}

func (ds *Datastore) GetProperties(limit, offset int, filter map[string]interface{}) ([]*Property, error) {

	where, values := generators.GenerateAndSearchClause(filter)

	sql := `SELECT p.id, p.id_address, p.name, p.details, p.created_date, ` +
		`address.id, address.id_country, address.address_line1, address.address_line2, ` +
		`address.address_line3, address.town_city, address.county, ` +
		`address.postcode, address.latitude, address.longitude, ` +
		`country.id, country.name, country.iso2 ` +
		`FROM places4all.property AS p ` +
		`JOIN places4all.address ON address.id = p.id_address ` +
		`JOIN places4all.country ON country.id = address.id_country ` +
		where + `ORDER BY p.id DESC LIMIT ` + strconv.Itoa(limit) +
		` OFFSET ` + strconv.Itoa(offset)
	sql = ds.postgres.Rebind(sql)

	rows, err := ds.postgres.Queryx(sql, values...)
	if err != nil {
		return nil, err
	}

	properties := make([]*Property, 0)
	for rows.Next() {
		p := NewProperty(false)
		p.SetExists()
		p.Address = NewAddress(true)
		p.Address.SetExists()
		err := rows.Scan(&p.Id, &p.IdAddress, &p.Name, &p.Details, &p.CreatedDate, &p.Address.Id,
			&p.Address.IdCountry, &p.Address.AddressLine1, &p.Address.AddressLine2, &p.Address.AddressLine3,
			&p.Address.TownCity, &p.Address.County, &p.Address.Postcode,
			&p.Address.Latitude, &p.Address.Longitude, &p.Address.Country.Id,
			&p.Address.Country.Name, &p.Address.Country.Iso2)
		if err != nil {
			return nil, err
		}
		p.Tags, err = ds.GetPropertyTagByIdProperty(p.Id)
		if err != nil {
			return nil, err
		}
		p.Owners, err = ds.GetPropertyClientsByIdProperty(p.Id)
		if err != nil {
			return nil, err
		}
		properties = append(properties, p)
		if err != nil {
			return nil, err
		}
	}

	return properties, err
}
