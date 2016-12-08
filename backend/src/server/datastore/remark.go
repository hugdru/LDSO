package datastore

import (
	"gopkg.in/guregu/null.v3/zero"
	"server/datastore/metadata"
	"errors"

)

type Remark struct {
	Id          	   int64       `json:"id" db:"id"`
	IdAudit            int64       `json:"idAudit" db:"id_audit"`
	IdCriterion        int64       `json:"idCriterion" db:"id_criterion"`
	Observation        zero.String `json:"observation" db:"observation"`
	Image              []byte      `json:"image" db:"image"`

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

func RRemark(allocateObjects bool) Remark {
	remark := Remark{}
	//if allocateObjects {
	//}
	return remark
}
func NewRemark(allocateObjects bool) *Remark {
	remark := RRemark(allocateObjects)
	return &remark
}

func (ds *Datastore) InsertRemark(r *Remark) error {

	if r.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.remark (` +
		`id_audit, id_criterion, observation, image` +
		`) VALUES (` +
		`$1, $2, $3 , $4 ` +
		`)`

	_, err := ds.postgres.Exec(sql, r.IdAudit, r.IdCriterion, r.Observation, r.Image)
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
		`observation, image` +
		`) = (` +
		`$1, $2` +
		`) WHERE id = $3`

	_, err := ds.postgres.Exec(sql, r.Observation, r.Image, r.Id)
	return err
}

func (ds *Datastore) SaveRemark(r *Remark) error {
	if r.Exists() {
		return ds.UpdateRemark(r)
	}

	return ds.InsertRemark(r)
}


func (ds *Datastore) GetRemarkByAuditCriterionIds(idAudit int64, idCriterion int64, idRemark int64) (*Remark, error) {

	var err error

	const sql = `SELECT
		remark.id, remark.id_audit, remark.id_criterion, remark.observation, remark.image
		WHERE remark.id = $1 and remark.id_audit=$2 and  remark.id_criterion=$3 `

	r := RRemark(true)
	r.SetExists()

	err = ds.postgres.QueryRow(sql, idAudit, idCriterion, idRemark).Scan(
		&r.Id, &r.IdAudit, &r.IdCriterion, &r.Observation, &r.Image,
	)

	if err != nil {
		return nil, err
	}

	return &r, err

}
func (ds *Datastore) GetRemarksByAuditCriterionIds(idAudit int64, idCriterion int64) ([]*Remark, error) {

	var err error

	const sql = `SELECT
		remark.id, remark.id_audit, remark.id_criterion, remark.observation, remark.image
		WHERE remark.id_audit=$2 and  remark.id_criterion=$3 `


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
