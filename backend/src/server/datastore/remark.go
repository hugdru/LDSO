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
	ImageMimetype zero.String `json:"-" db:"image_mimetype"`
	ImageHash     zero.String `json:"imageLocation" db:"image_hash"`

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

	if r == nil {
		return errors.New("remark should not be nil")
	}

	if r.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.remark (` +
		`id_audit, id_criterion, observation, image, image_mimetype, image_hash ` +
		`) VALUES (` +
		`$1, $2, $3, $4, $5, $6 ` +
		`) RETURNING id`

	err := ds.postgres.QueryRow(sql, r.IdAudit, r.IdCriterion, r.Observation, r.Image, r.ImageMimetype, r.ImageHash).Scan(&r.Id)
	if err != nil {
		return err
	}

	r.SetExists()

	return err
}

func (ds *Datastore) UpdateRemark(r *Remark) error {

	if r == nil {
		return errors.New("remark should not be nil")
	}

	if !r.Exists() {
		return errors.New("update failed: does not exist")
	}

	if r.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.remark SET (` +
		`observation, image, image_mimetype, image_hash` +
		`) = (` +
		`$1, $2, $3, $4` +
		`) WHERE id = $5 `

	_, err := ds.postgres.Exec(sql, r.Observation, r.Image, r.ImageMimetype, r.ImageHash, r.Id)
	return err
}

func (ds *Datastore) SaveRemark(r *Remark) error {

	if r == nil {
		return errors.New("remark should not be nil")
	}

	if r.Exists() {
		return ds.UpdateRemark(r)
	}

	return ds.InsertRemark(r)
}

func (ds *Datastore) GetRemarkByIdsAuditCriterionRemark(idAudit, idCriterion, idRemark int64) (*Remark, error) {

	const sql = `SELECT ` +
		`remark.id, remark.id_audit, remark.id_criterion, remark.observation, remark.image, remark.image_mimetype, remark.image_hash FROM places4all.remark ` +
		`WHERE remark.id_audit = $1 AND remark.id_criterion = $2 AND remark.id = $3`

	r := ARemark(true)
	r.SetExists()

	err := ds.postgres.QueryRow(sql, idAudit, idCriterion, idRemark).Scan(
		&r.Id, &r.IdAudit, &r.IdCriterion, &r.Observation, &r.Image, &r.ImageMimetype, &r.ImageHash,
	)

	if err != nil {
		return nil, err
	}

	return &r, err

}

func (ds *Datastore) GetRemarksByIdsAuditCriterion(idAudit int64, idCriterion int64) ([]*Remark, error) {

	const sql = `SELECT ` +
		`remark.id, remark.id_audit, remark.id_criterion, remark.observation, remark.image, remark.image_mimetype, image_hash FROM places4all.remark ` +
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

func (ds *Datastore) DeleteRemarkByIdsAuditCriterionRemark(idAudit, idCriterion, idRemark int64) error {

	const sql = `DELETE FROM places4all.remark ` +
		`WHERE remark.id_audit = $1 AND remark.id_criterion = $2 AND remark.id = $3`

	_, err := ds.postgres.Exec(sql, idAudit, idCriterion, idRemark)

	if err != nil {
		return err
	}

	return err
}

func (ds *Datastore) DeleteRemarkByIdsAuditCriterion(idAudit, idCriterion int64) error {

	const sql = `DELETE FROM places4all.remark ` +
		`WHERE remark.id_audit = $1 AND remark.id_criterion = $2`

	_, err := ds.postgres.Exec(sql, idAudit, idCriterion)

	if err != nil {
		return err
	}

	return err
}
