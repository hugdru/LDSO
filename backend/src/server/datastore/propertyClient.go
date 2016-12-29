package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type PropertyClient struct {
	IdProperty int64 `json:"idProperty" db:"id_property"`
	IdClient   int64 `json:"idClient" db:"id_client"`

	meta metadata.Metadata
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

func APropertyClient(allocateObjects bool) PropertyClient {
	propertyClient := PropertyClient{}
	//if allocateObjects {
	//}
	return propertyClient
}
func NewPropertyClient(allocateObjects bool) *PropertyClient {
	propertyClient := APropertyClient(allocateObjects)
	return &propertyClient
}

func (ds *Datastore) InsertPropertyClient(pc *PropertyClient) error {

	if pc == nil {
		return errors.New("propertyClient should not be nil")
	}

	if pc.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.property_client ` +
		`(id_property, id_client) VALUES ($1, $2)`

	_, err := ds.postgres.Exec(sql, pc.IdProperty, pc.IdClient)
	if err != nil {
		return err
	}

	pc.SetExists()

	return err
}

func (ds *Datastore) UpdatePropertyClient(pc *PropertyClient) error {

	//if pc == nil {
	//	return errors.New("propertyClient should not be nil")
	//}
	//
	//if !pc.Exists() {
	//	return errors.New("update failed: does not exist")
	//}
	//
	//if pc.Deleted() {
	//	return errors.New("update failed: marked for deletion")
	//}
	//
	//const sql = `UPDATE places4all.property_client SET (` +
	//	`` +
	//	`) = ( ` +
	//	`` +
	//	`) WHERE id_property = $1 AND id_client = $2`
	//
	//_, err := ds.postgres.Exec(sql, pc.IdProperty, pc.IdClient)
	//return err
	return errors.New("TO BE COMPLETED IF WE GET MORE DATABASE ROWS")
}

func (ds *Datastore) SavePropertyClient(pc *PropertyClient) error {
	if pc.Exists() {
		return ds.UpdatePropertyClient(pc)
	}

	return ds.InsertPropertyClient(pc)
}

func (ds *Datastore) UpsertPropertyClient(pc *PropertyClient) error {

	if pc == nil {
		return errors.New("propertyClient should not be nil")
	}

	if pc.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.property_client (` +
		`id_property, id_client` +
		`) VALUES (` +
		`$1, $2` +
		`) ON CONFLICT (id_property, id_client) DO UPDATE SET (` +
		`id_property, id_client` +
		`) = (` +
		`EXCLUDED.id_property, EXCLUDED.id_client` +
		`)`

	_, err := ds.postgres.Exec(sql, pc.IdProperty, pc.IdClient)
	if err != nil {
		return err
	}

	pc.SetExists()

	return err
}

func (ds *Datastore) DeletePropertyClient(pc *PropertyClient) error {

	if pc == nil {
		return errors.New("propertyClient should not be nil")
	}

	if !pc.Exists() {
		return nil
	}

	if pc.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.property_client WHERE id_property = $1 AND id_client = $2`

	_, err := ds.postgres.Exec(sql, pc.IdProperty, pc.IdClient)
	if err != nil {
		return err
	}

	pc.SetDeleted()

	return err
}

func (ds *Datastore) GetPropertyClientClient(pc *PropertyClient) (*Client, error) {
	return ds.GetClientByIdWithEntity(pc.IdClient)
}

func (ds *Datastore) GetPropertyClientProperty(pc *PropertyClient) (*Property, error) {
	return ds.GetPropertyByIdWithAddressTagsOwners(pc.IdProperty)
}

func (ds *Datastore) GetPropertyClientByIds(idProperty, idClient int64) (*PropertyClient, error) {

	const sql = `SELECT ` +
		`id_property, id_client ` +
		`FROM places4all.property_client ` +
		`WHERE id_property = $1 AND id_client = $2`

	pc := APropertyClient(false)
	pc.SetExists()

	err := ds.postgres.QueryRow(sql, idProperty, idClient).Scan(&pc.IdProperty, &pc.IdClient)
	if err != nil {
		return nil, err
	}

	return &pc, err
}

func (ds *Datastore) GetPropertyClientsByIdProperty(idProperty int64) ([]*Client, error) {

	const sql = `SELECT client.id_entity, entity.id, entity.name, entity.username, entity.image ` +
		`FROM places4all.entity ` +
		`JOIN places4all.client ON client.id_entity = entity.id ` +
		`JOIN places4all.property_client ON property_client.id_client = client.id_entity ` +
		`WHERE property_client.id_property = $1`

	rows, err := ds.postgres.Queryx(sql, idProperty)
	if err != nil {
		return nil, err
	}

	clients := make([]*Client, 0)
	for rows.Next() {
		client := NewClient(false)
		client.Entity = NewEntity(false)
		client.SetExists()
		err := rows.Scan(&client.IdEntity, &client.Entity.Id, &client.Entity.Name, &client.Entity.Username, &client.Entity.Image)
		if err != nil {
			return nil, err
		}
		clients = append(clients, client)
	}

	return clients, err
}
