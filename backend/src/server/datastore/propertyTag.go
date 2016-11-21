package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type PropertyTag struct {
	IdProperty int64 `json:"idProperty" db:"id_property"`
	IdTag      int64 `json:"idTag" db:"id_tag"`
	meta       metadata.Metadata
}

func (pt *PropertyTag) SetExists() {
	pt.meta.Exists = true
}

func (pt *PropertyTag) SetDeleted() {
	pt.meta.Deleted = true
}

func (pt *PropertyTag) Exists() bool {
	return pt.meta.Exists
}

func (pt *PropertyTag) Deleted() bool {
	return pt.meta.Deleted
}

func (ds *Datastore) InsertPropertyTag(pt *PropertyTag) error {
	var err error

	if pt.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.property_tag (` +
		`id_property` +
		`) VALUES (` +
		`$1` +
		`) RETURNING id_tag`

	err = ds.postgres.QueryRow(sql, pt.IdProperty).Scan(&pt.IdTag)
	if err != nil {
		return err
	}

	pt.SetExists()

	return nil
}

func (ds *Datastore) UpdatePropertyTag(pt *PropertyTag) error {
	var err error

	if !pt.Exists() {
		return errors.New("update failed: does not exist")
	}

	if pt.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.property_tag SET (` +
		`id_property` +
		`) = ( ` +
		`$1` +
		`) WHERE id_tag = $2`

	_, err = ds.postgres.Exec(sql, pt.IdProperty, pt.IdTag)
	return err
}

func (ds *Datastore) SavePropertyTag(pt *PropertyTag) error {
	if pt.Exists() {
		return ds.UpdatePropertyTag(pt)
	}

	return ds.InsertPropertyTag(pt)
}

func (ds *Datastore) UpsertPropertyTag(pt *PropertyTag) error {
	var err error

	if pt.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.property_tag (` +
		`id_property, id_tag` +
		`) VALUES (` +
		`$1, $2` +
		`) ON CONFLICT (id_tag) DO UPDATE SET (` +
		`id_property, id_tag` +
		`) = (` +
		`EXCLUDED.id_property, EXCLUDED.id_tag` +
		`)`

	_, err = ds.postgres.Exec(sql, pt.IdProperty, pt.IdTag)
	if err != nil {
		return err
	}

	pt.SetExists()

	return nil
}

func (ds *Datastore) DeletePropertyTag(pt *PropertyTag) error {
	var err error

	if !pt.Exists() {
		return nil
	}

	if pt.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.property_tag WHERE id_tag = $1`

	_, err = ds.postgres.Exec(sql, pt.IdTag)
	if err != nil {
		return err
	}

	pt.SetDeleted()

	return nil
}

func (ds *Datastore) GetPropertyTagProperty(pt *PropertyTag) (*Property, error) {
	return ds.GetPropertyById(pt.IdProperty)
}

func (ds *Datastore) GetPropertyTagTag(pt *PropertyTag) (*Tag, error) {
	return ds.GetTagById(pt.IdTag)
}

func (ds *Datastore) GetPropertyTagByIds(idProperty, idTag int64) (*PropertyTag, error) {
	var err error

	const sql = `SELECT ` +
		`id_property, id_tag ` +
		`FROM places4all.property_tag ` +
		`WHERE id_property = $1 AND id_tag = $2`

	pt := PropertyTag{}
	pt.SetExists()

	err = ds.postgres.QueryRow(sql, idProperty, idTag).Scan(&pt.IdProperty, &pt.IdTag)
	if err != nil {
		return nil, err
	}

	return &pt, nil
}
