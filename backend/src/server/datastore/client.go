package datastore

import (
	"database/sql"
	"errors"
	"server/datastore/generators"
	"server/datastore/metadata"
	"strconv"
)

func clientVisibility(restricted bool) string {
	const clientRestricted = "client.id_entity "
	const clientAll = "client.id_entity "
	if restricted {
		return clientRestricted
	}
	return clientAll
}

type Client struct {
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

func (c *Client) MustSet(idEntity int64) error {
	if idEntity != 0 {
		c.IdEntity = idEntity
	} else {
		return errors.New("idEntity must be set")
	}

	return nil
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
	return ds.InsertClientTx(nil, c)
}

func (ds *Datastore) InsertClientTx(tx *sql.Tx, c *Client) error {

	if c == nil {
		return errors.New("client should not be nil")
	}

	if c.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.client (` +
		`id_entity` +
		`) VALUES (` +
		`$1` +
		`)`

	var err error
	if tx != nil {
		_, err = tx.Exec(sql, c.IdEntity)

	} else {
		_, err = ds.postgres.Exec(sql, c.IdEntity)
	}
	if err != nil {
		return err
	}

	c.SetExists()

	return err
}

func (ds *Datastore) UpdateClient(c *Client) error {
	return ds.UpdateClientTx(nil, c)
}

func (ds *Datastore) UpdateClientTx(tx *sql.Tx, c *Client) error {

	//if c == nil {
	//	return errors.New("client should not be nil")
	//}
	//
	//if !c.Exists() {
	//	return errors.New("update failed: does not exist")
	//}
	//
	//if c.Deleted() {
	//	return errors.New("update failed: marked for deletion")
	//}
	//
	//const sql = `UPDATE places4all.client SET (` +
	//	`id_entity` +
	//	`) = ( ` +
	//	`$1` +
	//	`) WHERE id = $2`
	//
	//_, err := ds.postgres.Exec(sql, c.IdEntity, c.Id)
	//return err
	return errors.New("NOT IMPLEMENTED")
}

func (ds *Datastore) SaveClient(c *Client) error {
	return ds.SaveClientTx(nil, c)
}

func (ds *Datastore) SaveClientTx(tx *sql.Tx, c *Client) error {

	if c == nil {
		return errors.New("client should not be nil")
	}

	if c.Exists() {
		return ds.UpdateClientTx(tx, c)
	}

	return ds.InsertClientTx(tx, c)
}

func (ds *Datastore) UpsertClient(c *Client) error {

	//if c == nil {
	//	return errors.New("client should not be nil")
	//}
	//
	//if c.Exists() {
	//	return errors.New("insert failed: already exists")
	//}
	//
	//const sql = `INSERT INTO places4all.client (` +
	//	`id_entity` +
	//	`) VALUES (` +
	//	`$1` +
	//	`) ON CONFLICT (id_entity) DO UPDATE SET (` +
	//	`id_entity` +
	//	`) = (` +
	//	`EXCLUDED.id_entity` +
	//	`)`
	//
	//_, err := ds.postgres.Exec(sql, c.IdEntity)
	//if err != nil {
	//	return err
	//}
	//
	//c.SetExists()
	//
	//return err
	return errors.New("NOT IMPLEMENTED")
}

func (ds *Datastore) DeleteClient(c *Client) error {

	if c == nil {
		return errors.New("client should not be nil")
	}

	if !c.Exists() {
		return nil
	}

	if c.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.client WHERE id_entity = $1`

	_, err := ds.postgres.Exec(sql, c.IdEntity)
	if err != nil {
		return err
	}

	c.SetDeleted()

	return err
}

func (ds *Datastore) GetClientEntity(c *Client, withCountry, restricted bool) (*Entity, error) {
	return ds.GetEntityById(c.IdEntity, withCountry, restricted)
}

func (ds *Datastore) GetClients(limit, offset int, withEntity, restricted bool) ([]*Client, error) {
	rows, err := ds.postgres.Queryx(ds.postgres.Rebind(`SELECT ` +
		clientVisibility(restricted) +
		`FROM places4all.client ` +
		`ORDER BY client.id_entity DESC LIMIT ` + strconv.Itoa(limit) +
		`OFFSET ` + strconv.Itoa(offset)))

	if err != nil {
		return nil, err
	}

	client := make([]*Client, 0)
	for rows.Next() {
		c := NewClient(false)
		c.SetExists()
		err = rows.StructScan(c)
		if err != nil {
			return nil, err
		}
		if withEntity {
			c.Entity, err = ds.GetEntityById(c.IdEntity, true, restricted)
			if err != nil {
				return nil, err
			}
		}
		client = append(client, c)
	}

	return client, err
}

func (ds *Datastore) GetClientById(idEntity int64, withEntity, restricted bool) (*Client, error) {
	filter := make(map[string]interface{})
	filter["id_entity"] = idEntity
	return ds.GetClient(filter, withEntity, restricted)
}

func (ds *Datastore) GetClient(filter map[string]interface{}, withEntity, restricted bool) (*Client, error) {
	where, values := generators.GenerateAndSearchClause(filter)
	sql := ds.postgres.Rebind(`SELECT ` +
		clientVisibility(restricted) +
		`FROM places4all.client ` +
		where)

	a := AClient(false)
	a.SetExists()

	err := ds.postgres.QueryRowx(sql, values...).StructScan(&a)
	if err != nil {
		return nil, err
	}

	if withEntity {
		a.Entity, err = ds.GetEntityById(a.IdEntity, true, restricted)
		if err != nil {
			return nil, err
		}
	}

	return &a, err
}
