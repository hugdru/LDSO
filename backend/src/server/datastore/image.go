package datastore

import (
	"errors"
	"gopkg.in/guregu/null.v3/zero"
	"server/datastore/metadata"
	"time"
)

type Image struct {
	Id          int64       `json:"id" db:"id"`
	IdGallery   int64       `json:"idGallery" db:"id_gallery"`
	Name        string      `json:"name" db:"name"`
	Description zero.String `json:"description" db:"description"`
	Image       []byte      `json:"image" db:"image"`
	CreatedDate time.Time   `json:"createdDate" db:"createdDate"`

	meta metadata.Metadata
}

func (i *Image) SetExists() {
	i.meta.Exists = true
}

func (i *Image) SetDeleted() {
	i.meta.Deleted = true
}

func (i *Image) Exists() bool {
	return i.meta.Exists
}

func (i *Image) Deleted() bool {
	return i.meta.Deleted
}

func AImage(allocateObjects bool) Image {
	image := Image{}
	//if allocateObjects {
	//}
	return image
}

func NewImage(allocateObjects bool) *Image {
	image := AImage(allocateObjects)
	return &image
}

func (ds *Datastore) InsertImage(i *Image) error {

	if i.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.image (` +
		`id_gallery, name, description, image, created_date` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5` +
		`) RETURNING id`

	res, err := ds.postgres.Exec(sql, i.IdGallery, i.Name, i.Description, i.Image, i.CreatedDate)
	if err != nil {
		return err
	}
	i.Id, err = res.LastInsertId()
	if err != nil {
		return err
	}

	i.SetExists()

	return err
}

func (ds *Datastore) UpdateImage(i *Image) error {

	if !i.Exists() {
		return errors.New("update failed: does not exist")
	}

	if i.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.image SET (` +
		`id_gallery, name, description, image, created` +
		`) = ( ` +
		`$1, $2, $3, $4, $5` +
		`) WHERE id = $6`

	_, err := ds.postgres.Exec(sql, i.IdGallery, i.Name, i.Description, i.Image, i.CreatedDate, i.Id)
	return err
}

func (ds *Datastore) SaveImage(i *Image) error {
	if i.Exists() {
		return ds.UpdateImage(i)
	}

	return ds.InsertImage(i)
}

func (ds *Datastore) UpsertImage(i *Image) error {

	if i.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.image (` +
		`id, id_gallery, name, description, image, created_date` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, id_gallery, name, description, image, created_date` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.id_gallery, EXCLUDED.name, EXCLUDED.description, EXCLUDED.image, EXCLUDED.created_date` +
		`)`

	_, err := ds.postgres.Exec(sql, i.Id, i.IdGallery, i.Name, i.Description, i.Image, i.CreatedDate)
	if err != nil {
		return err
	}

	i.SetExists()

	return err
}

func (ds *Datastore) DeleteImage(i *Image) error {

	if !i.Exists() {
		return nil
	}

	if i.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.image WHERE id = $1`

	_, err := ds.postgres.Exec(sql, i.Id)
	if err != nil {
		return err
	}

	i.SetDeleted()

	return err
}

func (ds *Datastore) GetImageGallery(i *Image) (*Gallery, error) {
	return ds.GetGalleryById(i.IdGallery)
}

func (ds *Datastore) GetImageById(id int64) (*Image, error) {

	const sql = `SELECT ` +
		`id, id_gallery, name, description, image, created_date ` +
		`FROM places4all.image ` +
		`WHERE id = $1`

	i := Image{}
	i.SetExists()

	err := ds.postgres.QueryRowx(sql, id).StructScan(&i)
	if err != nil {
		return nil, err
	}

	return &i, err
}
