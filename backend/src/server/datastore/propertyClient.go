package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type PropertyClient struct {
	IdProperty int64 `json:"idProperty" db:"id_property"`
	IdClient   int64 `json:"idClient" db:"id_client"`
	meta       metadata.Metadata
}

func (p *PropertyClient) SetExists() {
	p.meta.Exists = true
}

func (p *PropertyClient) SetDeleted() {
	p.meta.Deleted = true
}

func (p *PropertyClient) Exists() bool {
	return p.meta.Exists
}

func (p *PropertyClient) Deleted() bool {
	return p.meta.Deleted
}

func (ds *Datastore) InsertPropertyClient(pc *PropertyClient) error {
	var err error

	if pc.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.property_client (` +
		`id_property` +
		`) VALUES (` +
		`$1` +
		`) RETURNING id_client`

	err = ds.postgres.QueryRow(sql, pc.IdProperty).Scan(&pc.IdClient)
	if err != nil {
		return err
	}

	pc.SetExists()

	return nil
}

func (ds *Datastore) UpdatePropertyClient(pc *PropertyClient) error {
	var err error

	if !pc.Exists() {
		return errors.New("update failed: does not exist")
	}

	if pc.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.property_client SET (` +
		`id_property` +
		`) = ( ` +
		`$1` +
		`) WHERE id_client = $2`

	_, err = ds.postgres.Exec(sql, pc.IdProperty, pc.IdClient)
	return err
}

func (ds *Datastore) SavePropertyClient(pc *PropertyClient) error {
	if pc.Exists() {
		return ds.UpdatePropertyClient(pc)
	}

	return ds.InsertPropertyClient(pc)
}

func (ds *Datastore) UpsertPropertyClient(pc *PropertyClient) error {
	var err error

	if pc.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.property_client (` +
		`id_property, id_client` +
		`) VALUES (` +
		`$1, $2` +
		`) ON CONFLICT (id_client) DO UPDATE SET (` +
		`id_property, id_client` +
		`) = (` +
		`EXCLUDED.id_property, EXCLUDED.id_client` +
		`)`

	_, err = ds.postgres.Exec(sql, pc.IdProperty, pc.IdClient)
	if err != nil {
		return err
	}

	pc.SetExists()

	return nil
}

func (ds *Datastore) DeletePropertyClient(pc *PropertyClient) error {
	var err error

	if !pc.Exists() {
		return nil
	}

	if pc.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.property_client WHERE id_client = $1`

	_, err = ds.postgres.Exec(sql, pc.IdClient)
	if err != nil {
		return err
	}

	pc.SetDeleted()

	return nil
}

func (ds *Datastore) GetPropertyClientClient(pc *PropertyClient) (*Client, error) {
	return ds.GetClientById(pc.IdClient)
}

func (ds *Datastore) GetPropertyClientProperty(pc *PropertyClient) (*Property, error) {
	return ds.GetPropertyById(pc.IdProperty)
}

func (ds *Datastore) GetPropertyClientByIds(idProperty, idClient int64) (*PropertyClient, error) {
	var err error

	const sql = `SELECT ` +
		`id_property, id_client ` +
		`FROM places4all.property_client ` +
		`WHERE id_property = $1 AND id_client = $2`

	pc := PropertyClient{}
	pc.SetExists()

	err = ds.postgres.QueryRow(sql, idProperty, idClient).Scan(&pc.IdProperty, &pc.IdClient)
	if err != nil {
		return nil, err
	}

	return &pc, nil
}
