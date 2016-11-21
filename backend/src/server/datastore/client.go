package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type Client struct {
	Id       int64 `json:"id" db:"id"`
	IdEntity int64 `json:"idPerson" db:"id_entity"`
	meta     metadata.Metadata
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

func (ds *Datastore) InsertClient(c *Client) error {
	var err error

	if c.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.client (` +
		`id_person` +
		`) VALUES (` +
		`$1` +
		`) RETURNING id`

	err = ds.postgres.QueryRow(sql, c.IdPerson).Scan(&c.Id)
	if err != nil {
		return err
	}

	c.SetExists()

	return nil
}

func (ds *Datastore) UpdateClient(c *Client) error {
	var err error

	if !c.Exists() {
		return errors.New("update failed: does not exist")
	}

	if c.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.client SET (` +
		`id_person` +
		`) = ( ` +
		`$1` +
		`) WHERE id = $2`

	_, err = ds.postgres.Exec(sql, c.IdPerson, c.Id)
	return err
}

func (ds *Datastore) SaveClient(c *Client) error {
	if c.Exists() {
		return ds.UpdateClient(c)
	}

	return ds.InsertClient(c)
}
func (ds *Datastore) UpsertClient(c *Client) error {
	var err error

	if c.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.client (` +
		`id, id_person` +
		`) VALUES (` +
		`$1, $2` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_person` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_person` +
		`)`

	_, err = ds.postgres.Exec(sql, c.Id, c.IdPerson)
	if err != nil {
		return err
	}

	c.SetExists()

	return nil
}

func (ds *Datastore) DeleteClient(c *Client) error {
	var err error

	if !c.Exists() {
		return nil
	}

	if c.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.client WHERE id = $1`

	_, err = ds.postgres.Exec(sql, c.Id)
	if err != nil {
		return err
	}

	c.SetDeleted()

	return nil
}

func (ds *Datastore) GetClientPerson(c *Client) (*Person, error) {
	return ds.GetPersonById(c.IdPerson)
}

func (ds *Datastore) GetClientById(id int64) (*Client, error) {
	var err error

	const sql = `SELECT ` +
		`id, id_person ` +
		`FROM places4all.client ` +
		`WHERE id = $1`

	c := Client{}
	c.SetExists()

	err = ds.postgres.QueryRow(sql, id).Scan(&c.Id, &c.IdPerson)
	if err != nil {
		return nil, err
	}

	return &c, nil
}
