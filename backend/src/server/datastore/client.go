package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type Client struct {
	Id       int64 `json:"id" db:"id"`
	IdEntity int64 `json:"IdEntity" db:"id_entity"`

	// Objects
	Entity *Entity `json:"entity,omitempty"`

	meta metadata.Metadata
}

func (c *Client) SetExists() {
	c.meta.Exists = true
}

func (c *Client) SetDeleted() {
	c.meta.Deleted = true
}

func (c *Client) Exists() bool {
	return c.meta.Exists
}

func (c *Client) Deleted() bool {
	return c.meta.Deleted
}

func AClient(allocateObjects bool) Client {
	client := Client{}
	if allocateObjects {
		client.Entity = NewEntity(allocateObjects)
	}
	return client
}

func NewClient(allocateObjects bool) *Client {
	client := AClient(allocateObjects)
	return &client
}

func (ds *Datastore) InsertClient(c *Client) error {

	if c.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.client (` +
		`id_entity` +
		`) VALUES (` +
		`$1` +
		`) RETURNING id`

	err := ds.postgres.QueryRow(sql, c.IdEntity).Scan(&c.Id)
	if err != nil {
		return err
	}

	c.SetExists()

	return err
}

func (ds *Datastore) UpdateClient(c *Client) error {

	if !c.Exists() {
		return errors.New("update failed: does not exist")
	}

	if c.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.client SET (` +
		`id_entity` +
		`) = ( ` +
		`$1` +
		`) WHERE id = $2`

	_, err := ds.postgres.Exec(sql, c.IdEntity, c.Id)
	return err
}

func (ds *Datastore) SaveClient(c *Client) error {
	if c.Exists() {
		return ds.UpdateClient(c)
	}

	return ds.InsertClient(c)
}
func (ds *Datastore) UpsertClient(c *Client) error {

	if c.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.client (` +
		`id, id_entity` +
		`) VALUES (` +
		`$1, $2` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_entity` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_entity` +
		`)`

	_, err := ds.postgres.Exec(sql, c.Id, c.IdEntity)
	if err != nil {
		return err
	}

	c.SetExists()

	return err
}

func (ds *Datastore) DeleteClient(c *Client) error {

	if !c.Exists() {
		return nil
	}

	if c.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.client WHERE id = $1`

	_, err := ds.postgres.Exec(sql, c.Id)
	if err != nil {
		return err
	}

	c.SetDeleted()

	return err
}

func (ds *Datastore) GetClientEntity(c *Client) (*Entity, error) {
	return ds.GetEntityById(c.IdEntity)
}

func (ds *Datastore) GetClientById(id int64) (*Client, error) {

	const sql = `SELECT ` +
		`client.id, client.id_entity, ` +
		`entity.id, entity.id_country, entity.name, entity.email, ` +
		`entity.username, entity.password, entity.image, entity.banned, entity.banned_date, ` +
		`entity.reason, entity.mobilephone, entity.telephone, entity.created_date, ` +
		`country.id, country.name, country.iso2 ` +
		`FROM places4all.client ` +
		`JOIN places4all.entity ON entity.id = client.id_entity ` +
		`JOIN places4all.country ON country.id = entity.id_country ` +
		`WHERE client.id = $1`

	c := NewClient(true)
	c.SetExists()

	err := ds.postgres.QueryRow(sql, id).Scan(
		&c.Id, &c.IdEntity,
		&c.Entity.Id, &c.Entity.IdCountry, &c.Entity.Name, &c.Entity.Email,
		&c.Entity.Username, &c.Entity.Password, &c.Entity.Image, &c.Entity.Banned, &c.Entity.BannedDate,
		&c.Entity.Reason, &c.Entity.Mobilephone, &c.Entity.Telephone, &c.Entity.CreatedDate,
		&c.Entity.Country.Id, &c.Entity.Country.Name, &c.Entity.Country.Iso2,
	)
	if err != nil {
		return nil, err
	}

	return c, err
}
