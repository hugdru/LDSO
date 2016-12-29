package datastore

import (
	"database/sql"
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

func (p *Property) AllSetIfNotEmptyOrNil(name, details string) error {
	if name != "" {
		p.Name = name
	}
	if details != "" {
		p.Details = details
	}

	return nil
}

func (p *Property) MustSet(name, details string) error {
	if name != "" {
		p.Name = name
	} else {
		return errors.New("name must be set")
	}
	if details != "" {
		p.Details = details
	} else {
		return errors.New("details must be set")
	}

	return nil
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

func (ds *Datastore) InsertPropertyTx(tx *sql.Tx, p *Property) error {

	if p == nil {
		return errors.New("property should not be nil")
	}

	if p.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.property (` +
		`id_address, name, details, created_date` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`) RETURNING id`

	var err error
	if tx != nil {
		err = tx.QueryRow(sql, p.IdAddress, p.Name, p.Details, p.CreatedDate).Scan(&p.Id)
	} else {
		err = ds.postgres.QueryRow(sql, p.IdAddress, p.Name, p.Details, p.CreatedDate).Scan(&p.Id)
	}

	if err != nil {
		return err
	}

	p.SetExists()

	return err
}

func (ds *Datastore) InsertProperty(p *Property) error {
	return ds.InsertPropertyTx(nil, p)
}

func (ds *Datastore) UpdateProperty(p *Property) error {
	return ds.UpdatePropertyTx(nil, p)
}

func (ds *Datastore) UpdatePropertyTx(tx *sql.Tx, p *Property) error {

	if p == nil {
		return errors.New("property should not be nil")
	}

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

	var err error
	if tx != nil {
		_, err = tx.Exec(sql, p.IdAddress, p.Name, p.Details, p.CreatedDate, p.Id)

	} else {
		_, err = ds.postgres.Exec(sql, p.IdAddress, p.Name, p.Details, p.CreatedDate, p.Id)
	}
	return err
}

func (ds *Datastore) SaveProperty(p *Property) error {

	if p == nil {
		return errors.New("property should not be nil")
	}

	if p.Exists() {
		return ds.UpdateProperty(p)
	}

	return ds.InsertProperty(p)
}

func (ds *Datastore) SavePropertyWithAddress(p *Property) error {

	if p == nil {
		return errors.New("property should not be nil")
	}

	if p.Exists() {
		return ds.UpdatePropertyWithAddress(p)
	}

	return ds.InsertPropertyWithAddress(p)
}

func (ds *Datastore) UpdatePropertyWithAddress(p *Property) error {

	if p == nil {
		return errors.New("property should not be nil")
	}

	tx, err := ds.postgres.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	err = ds.UpdateAddressTx(tx, p.Address)
	if err != nil {
		return err
	}

	err = ds.UpdatePropertyTx(tx, p)
	if err != nil {
		return err
	}

	return err
}

func (ds *Datastore) InsertPropertyWithAddress(p *Property) error {

	if p == nil {
		return errors.New("property should not be nil")
	}

	tx, err := ds.postgres.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	err = ds.InsertAddressTx(tx, p.Address)
	if err != nil {
		return err
	}

	err = ds.InsertPropertyTx(tx, p)
	if err != nil {
		return err
	}

	return err
}

func (ds *Datastore) UpsertProperty(p *Property) error {

	if p == nil {
		return errors.New("property should not be nil")
	}

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

func (ds *Datastore) DeletePropertyTx(tx *sql.Tx, p *Property) error {

	if p == nil {
		return errors.New("property should not be nil")
	}

	if !p.Exists() {
		return nil
	}

	if p.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.property WHERE id = $1`

	var err error
	if tx != nil {
		_, err = tx.Exec(sql, p.Id)

	} else {
		_, err = ds.postgres.Exec(sql, p.Id)

	}
	if err != nil {
		return err
	}

	p.SetDeleted()

	return err
}

func (ds *Datastore) DeleteProperty(p *Property) error {
	return ds.DeletePropertyTx(nil, p)
}

func (ds *Datastore) DeletePropertyWithAddress(p *Property) error {
	if p == nil {
		return errors.New("property should not be nil")
	}

	tx, err := ds.postgres.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	err = ds.DeletePropertyTx(tx, p)
	if err != nil {
		return err
	}
	err = ds.DeleteAddressTx(tx, p.Address)
	if err != nil {
		return err
	}

	return err
}

func (ds *Datastore) GetPropertyAddress(p *Property) (*Address, error) {
	return ds.GetAddressByIdWithCountry(p.IdAddress)
}

func (ds *Datastore) GetPropertyByIdWithAddressTagsOwners(id int64) (*Property, error) {

	p := NewProperty(false)
	err := ds.postgres.QueryRowx(`SELECT `+
		`id, id_address, name, details, created_date `+
		`FROM places4all.property `+
		`WHERE id = $1`,
		id).StructScan(p)
	if err != nil {
		return nil, err
	}

	p.Address, err = ds.GetAddressByIdWithCountry(p.IdAddress)
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

func (ds *Datastore) GetPropertiesWithForeignAddressTagsOwners(limit, offset int, filter map[string]interface{}) ([]*Property, error) {

	where, values := generators.GenerateAndSearchClause(filter)

	sql := ds.postgres.Rebind(`SELECT id, id_address, name, details, created_date ` +
		`FROM places4all.property ` +
		where + `ORDER BY property.id DESC LIMIT ` + strconv.Itoa(limit) +
		` OFFSET ` + strconv.Itoa(offset))

	rows, err := ds.postgres.Queryx(sql, values...)
	if err != nil {
		return nil, err
	}

	properties := make([]*Property, 0)
	for rows.Next() {
		p := NewProperty(false)
		p.SetExists()
		err := rows.StructScan(p)
		if err != nil {
			return nil, err
		}
		p.Address, err = ds.GetAddressByIdWithCountry(p.IdAddress)
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
