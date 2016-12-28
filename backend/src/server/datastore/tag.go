package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type Tag struct {
	Id   int64  `json:"id" db:"id"`
	Name string `json:"name" db:"name"`

	meta metadata.Metadata
}

func (t *Tag) SetExists() {
	t.meta.Exists = true
}

func (t *Tag) SetDeleted() {
	t.meta.Deleted = true
}

func (t *Tag) Exists() bool {
	return t.meta.Exists
}

func (t *Tag) Deleted() bool {
	return t.meta.Deleted
}

func ATag(allocateObjects bool) Tag {
	tag := Tag{}
	//if allocateObjects {
	//}
	return tag
}

func NewTag(allocateObjects bool) *Tag {
	tag := ATag(allocateObjects)
	return &tag
}

func (ds *Datastore) InsertTag(t *Tag) error {

	if t == nil {
		return errors.New("tag should not be nil")
	}

	if t.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.tag (` +
		`name` +
		`) VALUES (` +
		`$1` +
		`) RETURNING id`

	err := ds.postgres.QueryRow(sql, t.Name).Scan(&t.Id)
	if err != nil {
		return err
	}

	t.SetExists()

	return err
}

func (ds *Datastore) UpdateTag(t *Tag) error {

	if t == nil {
		return errors.New("tag should not be nil")
	}

	if !t.Exists() {
		return errors.New("update failed: does not exist")
	}

	if t.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.tag SET (` +
		`name` +
		`) = ( ` +
		`$1` +
		`) WHERE id = $2`

	_, err := ds.postgres.Exec(sql, t.Name, t.Id)
	return err
}

func (ds *Datastore) SaveTag(t *Tag) error {
	if t.Exists() {
		return ds.UpdateTag(t)
	}

	return ds.InsertTag(t)
}

func (ds *Datastore) UpsertTag(t *Tag) error {

	if t == nil {
		return errors.New("tag should not be nil")
	}

	if t.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.tag (` +
		`id, name` +
		`) VALUES (` +
		`$1, $2` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, name` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.name` +
		`)`

	_, err := ds.postgres.Exec(sql, t.Id, t.Name)
	if err != nil {
		return err
	}

	t.SetExists()

	return err
}

func (ds *Datastore) DeleteTag(t *Tag) error {

	if t == nil {
		return errors.New("tag should not be nil")
	}

	if !t.Exists() {
		return nil
	}

	if t.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.tag WHERE id = $1`

	_, err := ds.postgres.Exec(sql, t.Id)
	if err != nil {
		return err
	}

	t.SetDeleted()

	return err
}

func (ds *Datastore) GetTagByName(name string) (*Tag, error) {

	const sql = `SELECT ` +
		`id, name ` +
		`FROM places4all.tag ` +
		`WHERE name = $1`

	t := Tag{}
	t.SetExists()

	err := ds.postgres.QueryRow(sql, name).Scan(&t.Id, &t.Name)
	if err != nil {
		return nil, err
	}

	return &t, err
}

func (ds *Datastore) GetTagById(id int64) (*Tag, error) {

	const sql = `SELECT ` +
		`id, name ` +
		`FROM places4all.tag ` +
		`WHERE id = $1`

	t := Tag{}
	t.SetExists()

	err := ds.postgres.QueryRow(sql, id).Scan(&t.Id, &t.Name)
	if err != nil {
		return nil, err
	}

	return &t, err
}
