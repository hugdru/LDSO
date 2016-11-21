package datastore

import (
	"errors"
	"server/datastore/metadata"
	"time"
)

type Property struct {
	Id        int64      `json:"id" db:"id"`
	IdAddress int64      `json:"idAddress, omitempty" db:"id_address"`
	Name      string     `json:"name" db:"name"`
	Details   string     `json:"details" db:"details"`
	Created   *time.Time `json:"created" db:"created"`

	// TODO: Add more data to property.sql, do the get and the rest
	Address   *Address   `json:"address, omitempty"`
	Owners    *[]Client  `json:"owners, omitempty"`
	Tags      *[]Tag     `json:"tags, omitempty"`
	Galleries *[]Gallery `json:"galleries, omitempty"`

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

func (ds *Datastore) InsertProperty(p *Property) error {
	var err error

	if p.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.property (` +
		`id_address, name, details, created` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`) RETURNING id`

	err = ds.postgres.QueryRow(sql, p.IdAddress, p.Name, p.Details, p.Created).Scan(&p.Id)
	if err != nil {
		return err
	}

	p.SetExists()

	return nil
}

func (ds *Datastore) UpdateProperty(p *Property) error {
	var err error

	if !p.Exists() {
		return errors.New("update failed: does not exist")
	}

	if p.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.property SET (` +
		`id_address, name, details, created` +
		`) = ( ` +
		`$1, $2, $3, $4` +
		`) WHERE id = $5`

	_, err = ds.postgres.Exec(sql, p.IdAddress, p.Name, p.Details, p.Created, p.Id)
	return err
}

func (ds *Datastore) SaveProperty(p *Property) error {
	if p.Exists() {
		return ds.UpdateProperty(p)
	}

	return ds.InsertProperty(p)
}

func (ds *Datastore) UpsertProperty(p *Property) error {
	var err error

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

	_, err = ds.postgres.Exec(sql, p.Id, p.IdAddress, p.Name, p.Details, p.Created)
	if err != nil {
		return err
	}

	p.SetExists()

	return nil
}

func (ds *Datastore) DeleteProperty(p *Property) error {
	var err error

	if !p.Exists() {
		return nil
	}

	if p.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.property WHERE id = $1`

	_, err = ds.postgres.Exec(sql, p.Id)
	if err != nil {
		return err
	}

	p.SetDeleted()

	return nil
}

func (ds *Datastore) GetPropertyAddress(p *Property) (*Address, error) {
	return ds.GetAddressById(p.IdAddress)
}

func (ds *Datastore) GetPropertyById(id int64) (*Property, error) {
	var err error

	const sql = `SELECT ` +
		`id, id_address, name, details, created ` +
		`FROM places4all.property ` +
		`WHERE id = $1`

	p := Property{}
	p.SetExists()

	err = ds.postgres.QueryRow(sql, id).Scan(&p.Id, &p.IdAddress, &p.Name, &p.Details, &p.Created)
	if err != nil {
		return nil, err
	}

	return &p, nil
}
