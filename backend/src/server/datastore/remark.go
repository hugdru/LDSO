package datastore

import (
	"errors"
	"gopkg.in/guregu/null.v3/zero"
	"server/datastore/metadata"
)

type Remark struct {
	Id            int64       `json:"id" db:"id"`
	IdAudit       int64       `json:"idAudit" db:"id_audit"`
	IdCriterion   int64       `json:"idCriterion" db:"id_criterion"`
	Observation   zero.String `json:"observation" db:"observation"`
	Image         []byte      `json:"-" db:"image"`
	ImageMimeType zero.String `json:"imageMimeType" db:"image_mimetype"`

	meta metadata.Metadata
}

func (r *Remark) SetExists() {
	r.meta.Exists = true
}

func (r *Remark) SetDeleted() {
	r.meta.Deleted = true
}

func (r *Remark) Exists() bool {
	return r.meta.Exists
}

func (r *Remark) Deleted() bool {
	return r.meta.Deleted
}

func ARemark(allocateObjects bool) Remark {
	remark := Remark{}
	//if allocateObjects {
	//}
	return remark
}
func NewRemark(allocateObjects bool) *Remark {
	remark := ARemark(allocateObjects)
	return &remark
}

func (ds *Datastore) InsertRemark(r *Remark) error {

	if r.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.remark (` +
		`id_audit, id_criterion, observation, image, image_mimetype ` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5 ` +
		`) RETURNING id`

	err := ds.postgres.QueryRow(sql, r.IdAudit, r.IdCriterion, r.Observation, r.Image, r.ImageMimeType).Scan(&r.Id)
	if err != nil {
		return err
	}

	r.SetExists()

	return err
}

func (ds *Datastore) UpdateRemark(r *Remark) error {

	if !r.Exists() {
		return errors.New("update failed: does not exist")
	}

	if r.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.remark SET (` +
		`observation, image, image_mimetype` +
		`) = (` +
		`$1, $2, $3` +
		`) WHERE id = $4 `

	_, err := ds.postgres.Exec(sql, r.Observation, r.Image, r.ImageMimeType, r.Id)
	return err
}

func (ds *Datastore) SaveRemark(r *Remark) error {
	if r.Exists() {
		return ds.UpdateRemark(r)
	}

	return ds.InsertRemark(r)
}

func (ds *Datastore) GetRemarkByIdsAuditCriterionRemark(idAudit, idCriterion, idRemark int64) (*Remark, error) {

	var err error

	const sql = `SELECT ` +
		`remark.id, remark.id_audit, remark.id_criterion, remark.observation, remark.image, remark.image_mimetype FROM places4all.remark ` +
		`WHERE remark.id_audit = $1 AND remark.id_criterion = $2 AND remark.id = $3`

	r := ARemark(true)
	r.SetExists()

	err = ds.postgres.QueryRow(sql, idRemark).Scan(
		&r.Id, &r.IdAudit, &r.IdCriterion, &r.Observation, &r.Image, &r.ImageMimeType,
	)

	if err != nil {
		return nil, err
	}

	return &r, err

}
func (ds *Datastore) GetRemarksByAuditCriterionIds(idAudit int64, idCriterion int64) ([]*Remark, error) {

	var err error

	const sql = `SELECT ` +
		`remark.id, remark.id_audit, remark.id_criterion, remark.observation, remark.image, remark.image_mimetype FROM places4all.remark ` +
		`WHERE remark.id_audit = $1 AND remark.id_criterion = $2`

	rows, err := ds.postgres.Queryx(sql, idAudit, idCriterion)

	if err != nil {
		return nil, err
	}

	remarks := make([]*Remark, 0)
	for rows.Next() {
		c := NewRemark(true)
		c.SetExists()
		err = rows.StructScan(c)
		if err != nil {
			return nil, err
		}
		remarks = append(remarks, c)
	}

	return remarks, err

}
