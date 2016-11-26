package datastore

import (
	"encoding/json"
	"errors"
	"server/datastore/metadata"
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

	res, err := ds.postgres.Exec(sql, p.IdAddress, p.Name, p.Details, p.CreatedDate)
	if err != nil {
		return err
	}
	p.Id, err = res.LastInsertId()
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
	return ds.GetAddressById(p.IdAddress)
}

func (ds *Datastore) GetPropertyById(id int64) (*Property, error) {

	p := NewProperty(false)
	p.Address = NewAddress(true)
	var tagsJson, ownersJson string
	err := ds.postgres.QueryRow(`
	SELECT p.id, p.id_address, p.name, p.details, p.created_date, address.id,
	address.address_line1, address.address_line2, address.address_line3,
		address.town_city, address.county, address.postcode,
		address.latitude, address.longitude, country.id,
		country.name, country.iso2,
		(SELECT json_agg(ts)
		FROM (
			SELECT tag.id, tag.name
			FROM places4all.property_tag
			JOIN places4all.tag on tag.id = property_tag.id_tag
			WHERE property_tag.id_property = p.id
		) ts) AS tags,
		(SELECT json_agg(owrs)
		FROM (
			SELECT entity.id, entity.name, entity.image
			FROM places4all.entity
			JOIN places4all.client ON client.id_entity = entity.id
			JOIN places4all.property_client ON property_client.id_client = client.id
			WHERE id_property = p.id
		) owrs) AS owners
	FROM places4all.property AS p
	JOIN places4all.address ON address.id = p.id_address
	JOIN places4all.country ON country.id = address.id_country
	WHERE p.id = $1;
	`, id).Scan(&p.Id, &p.IdAddress, &p.Name, &p.Details, &p.CreatedDate, &p.Address.Id,
		&p.Address.AddressLine1, &p.Address.AddressLine2, &p.Address.AddressLine3,
		&p.Address.TownCity, &p.Address.County, &p.Address.Postcode,
		&p.Address.Latitude, &p.Address.Longitude, &p.Address.Country.Id,
		&p.Address.Country.Name, &p.Address.Country.Iso2,
		&tagsJson, &ownersJson,
	)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal([]byte(tagsJson), &p.Tags)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal([]byte(ownersJson), &p.Owners)
	if err != nil {
		return nil, err
	}

	return p, err
}
