package datastore

import (
	"database/sql"
	"errors"
	"server/datastore/metadata"
	"time"
)

type Gallery struct {
	Id          int64          `json:"id" db:"id"`
	IdProperty  int64          `json:"idProperty" db:"id_property"`
	Name        string         `json:"name" db:"name"`
	Description sql.NullString `json:"description" db:"description"`
	Created     *time.Time     `json:"created" db:"created"`
	meta        metadata.Metadata
}

func (g *Gallery) SetExists() {
	g.meta.Exists = true
}

func (g *Gallery) SetDeleted() {
	g.meta.Deleted = true
}

func (g *Gallery) Exists() bool {
	return g.meta.Exists
}

func (g *Gallery) Deleted() bool {
	return g.meta.Deleted
}

func (ds *Datastore) InsertGallery(g *Gallery) error {
	var err error

	if g.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.gallery (` +
		`id_property, name, description, created` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`) RETURNING id`

	err = ds.postgres.QueryRow(sql, g.IdProperty, g.Name, g.Description, g.Created).Scan(&g.Id)
	if err != nil {
		return err
	}

	g.SetExists()

	return nil
}

func (ds *Datastore) UpdateGallery(g *Gallery) error {
	var err error

	if !g.Exists() {
		return errors.New("update failed: does not exist")
	}

	if g.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.gallery SET (` +
		`id_property, name, description, created` +
		`) = ( ` +
		`$1, $2, $3, $4` +
		`) WHERE id = $5`

	_, err = ds.postgres.Exec(sql, g.IdProperty, g.Name, g.Description, g.Created, g.Id)
	return err
}

func (ds *Datastore) SaveGallery(g *Gallery) error {
	if g.Exists() {
		return ds.UpdateGallery(g)
	}

	return ds.InsertGallery(g)
}

func (ds *Datastore) UpsertGallery(g *Gallery) error {
	var err error

	if g.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.gallery (` +
		`id, id_property, name, description, created` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_property, name, description, created` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_property, EXCLUDED.name, EXCLUDED.description, EXCLUDED.created` +
		`)`

	_, err = ds.postgres.Exec(sql, g.Id, g.IdProperty, g.Name, g.Description, g.Created)
	if err != nil {
		return err
	}

	g.SetExists()

	return nil
}

func (ds *Datastore) DeleteGallery(g *Gallery) error {
	var err error

	if !g.Exists() {
		return nil
	}

	if g.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.gallery WHERE id = $1`

	_, err = ds.postgres.Exec(sql, g.Id)
	if err != nil {
		return err
	}

	g.SetDeleted()

	return nil
}

func (ds *Datastore) GetGalleryProperty(g *Gallery) (*Property, error) {
	return ds.GetPropertyById(g.IdProperty)
}

func (ds *Datastore) GetGalleryById(id int64) (*Gallery, error) {
	var err error

	const sql = `SELECT ` +
		`id, id_property, name, description, created ` +
		`FROM places4all.gallery ` +
		`WHERE id = $1`

	g := Gallery{}
	g.SetExists()

	err = ds.postgres.QueryRow(sql, id).Scan(&g.Id, &g.IdProperty, &g.Name, &g.Description, &g.Created)
	if err != nil {
		return nil, err
	}

	return &g, nil
}
