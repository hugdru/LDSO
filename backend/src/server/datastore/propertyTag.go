package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type PropertyTag struct {
	IdProperty int64 `json:"idProperty" db:"id_property"`
	IdTag      int64 `json:"idTag" db:"id_tag"`

	meta metadata.Metadata
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

func APropertyTag(allocateObjects bool) PropertyTag {
	propertyTag := PropertyTag{}
	//if allocateObjects {
	//}
	return propertyTag
}
func NewPropertyTag(allocateObjects bool) *PropertyTag {
	propertyTag := APropertyTag(allocateObjects)
	return &propertyTag
}

func (ds *Datastore) InsertPropertyTag(pt *PropertyTag) error {

	if pt.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.property_tag (` +
		`id_property, id_tag` +
		`) VALUES (` +
		`$1, $2` +
		`)`

	_, err := ds.postgres.Exec(sql, pt.IdProperty, pt.IdTag)
	if err != nil {
		return err
	}

	pt.SetExists()

	return err
}

func (ds *Datastore) UpdatePropertyTag(pt *PropertyTag) error {

	//if !pt.Exists() {
	//	return errors.New("update failed: does not exist")
	//}
	//
	//if pt.Deleted() {
	//	return errors.New("update failed: marked for deletion")
	//}
	//
	//const sql = `UPDATE places4all.property_tag SET (` +
	//	`` +
	//	`) = ( ` +
	//	`` +
	//	`) WHERE id_property $1 AND id_tag = $2`
	//
	//_, err := ds.postgres.Exec(sql, pt.IdProperty, pt.IdTag)
	//return err
	return errors.New("TO BE COMPLETED IF WE GET MORE DATABASE ROWS")
}

func (ds *Datastore) SavePropertyTag(pt *PropertyTag) error {
	if pt.Exists() {
		return ds.UpdatePropertyTag(pt)
	}

	return ds.InsertPropertyTag(pt)
}

func (ds *Datastore) UpsertPropertyTag(pt *PropertyTag) error {

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

	_, err := ds.postgres.Exec(sql, pt.IdProperty, pt.IdTag)
	if err != nil {
		return err
	}

	pt.SetExists()

	return err
}

func (ds *Datastore) DeletePropertyTag(pt *PropertyTag) error {

	if !pt.Exists() {
		return nil
	}

	if pt.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.property_tag WHERE id_property = $1 AND id_tag = $2`

	_, err := ds.postgres.Exec(sql, pt.IdProperty, pt.IdTag)
	if err != nil {
		return err
	}

	pt.SetDeleted()

	return err
}

func (ds *Datastore) GetPropertyTagProperty(pt *PropertyTag) (*Property, error) {
	return ds.GetPropertyById(pt.IdProperty)
}

func (ds *Datastore) GetPropertyTagTag(pt *PropertyTag) (*Tag, error) {
	return ds.GetTagById(pt.IdTag)
}

func (ds *Datastore) GetPropertyTagByIds(idProperty, idTag int64) (*PropertyTag, error) {

	const sql = `SELECT ` +
		`id_property, id_tag ` +
		`FROM places4all.property_tag ` +
		`WHERE id_property = $1 AND id_tag = $2`

	pt := APropertyTag(false)
	pt.SetExists()

	err := ds.postgres.QueryRow(sql, idProperty, idTag).Scan(&pt.IdProperty, &pt.IdTag)
	if err != nil {
		return nil, err
	}

	return &pt, err
}

func (ds *Datastore) GetPropertyTagByIdProperty(idProperty int64) ([]*Tag, error) {

	const sql = `SELECT ` +
		`tag.id, tag.name ` +
		`FROM places4all.tag ` +
		`JOIN places4all.property_tag ON property_tag.id_tag = tag.id ` +
		`WHERE property_tag.id_property = $1`

	rows, err := ds.postgres.Queryx(sql, idProperty)
	if err != nil {
		return nil, err
	}

	tags := make([]*Tag, 0)
	for rows.Next() {
		tag := NewTag(false)
		err := rows.StructScan(tag)
		if err != nil {
			return nil, err
		}
		tags = append(tags, tag)
	}

	return tags, err
}
