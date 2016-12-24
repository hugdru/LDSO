package datastore

import (
	"errors"
	"gopkg.in/guregu/null.v3/zero"
	"server/datastore/generators"
	"server/datastore/metadata"
)

type AuditCriterion struct {
	IdAudit     int64       `json:"id_audit" db:"id_audit"`
	IdCriterion int64       `json:"id_criterion" db:"id_criterion"`
	Value       zero.Int    `json:"value" db:"value"`
	Observation zero.String `json:"observation" db:"observation"`

	meta metadata.Metadata
}

func (ac *AuditCriterion) SetExists() {
	ac.meta.Exists = true
}

func (ac *AuditCriterion) SetDeleted() {
	ac.meta.Deleted = true
}

func (ac *AuditCriterion) Exists() bool {
	return ac.meta.Exists
}

func (ac *AuditCriterion) Deleted() bool {
	return ac.meta.Deleted
}

func AAuditCriterion(allocateObjects bool) AuditCriterion {
	auditCriterion := AuditCriterion{}
	//if allocateObjects {
	//}
	return auditCriterion
}
func NewAuditCriterion(allocateObjects bool) *AuditCriterion {
	auditCriterion := AAuditCriterion(allocateObjects)
	return &auditCriterion
}

func (ds *Datastore) InsertAuditCriterion(ac *AuditCriterion) error {

	if ac.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.audit_criterion (` +
		`id_audit, id_criterion, value, observation` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`)`

	_, err := ds.postgres.Exec(sql, ac.IdAudit, ac.IdCriterion, ac.Value, ac.Observation)
	if err != nil {
		return err
	}

	ac.SetExists()

	return err
}

func (ds *Datastore) UpdateAuditCriterion(ac *AuditCriterion) error {

	if !ac.Exists() {
		return errors.New("update failed: does not exist")
	}

	if ac.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.audit_criterion SET (` +
		`value, observation` +
		`) = (` +
		`$1, $2` +
		`) WHERE id_audit = $3 AND id_criterion = $4`

	_, err := ds.postgres.Exec(sql, ac.Value, ac.Observation, ac.IdAudit, ac.IdCriterion)
	return err
}

func (ds *Datastore) SaveAuditCriterion(ac *AuditCriterion) error {
	if ac.Exists() {
		return ds.UpdateAuditCriterion(ac)
	}

	return ds.InsertAuditCriterion(ac)
}

func (ds *Datastore) DeleteAuditCriterion(ac *AuditCriterion) error {

	if !ac.Exists() {
		return nil
	}

	if ac.Deleted() {
		return nil
	}

	sql := ds.postgres.Rebind(`DELETE FROM places4all.audit_criterion WHERE id_audit = ? AND id_criterion = ?`)

	_, err := ds.postgres.Exec(sql, ac.IdAudit, ac.IdCriterion)
	if err != nil {
		return err
	}

	ac.SetDeleted()

	return err
}

func (ds *Datastore) UpsertAuditCriterion(ac *AuditCriterion) error {

	if ac.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.audit_criterion (` +
		`id_audit, id_criterion, value, observation` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`) ON CONFLICT (id_audit, id_criterion) DO UPDATE SET (` +
		`id_audit, id_criterion, value, observation` +
		`) = (` +
		`EXCLUDED.id_audit, EXCLUDED.id_criterion, EXCLUDED.value, EXCLUDED.observation` +
		`)`

	_, err := ds.postgres.Exec(sql, ac.IdAudit, ac.IdCriterion, ac.Value, ac.Observation)
	if err != nil {
		return err
	}

	ac.SetExists()

	return err
}

func (ds *Datastore) DeleteAuditCriterionByIds(idAudit, idCriterion int64) error {
	const sql = `DELETE FROM places4all.audit_criterion WHERE id_audit = $1 AND id_criterion = $2`

	_, err := ds.postgres.Exec(sql, idAudit, idCriterion)
	if err != nil {
		return err
	}

	return err
}

func (ds *Datastore) DeleteAuditCriterionByIdAudit(idAudit int64) error {
	const sql = `DELETE FROM places4all.audit_criterion WHERE id_audit = $1`

	_, err := ds.postgres.Exec(sql, idAudit)
	if err != nil {
		return err
	}

	return err
}

func (ds *Datastore) GetAuditCriterionAudit(ac *AuditCriterion) (*Audit, error) {
	return ds.GetAuditById(ac.IdAudit)
}

func (ds *Datastore) GetAuditCriterionCriterion(ac *AuditCriterion) (*Criterion, error) {
	return ds.GetCriterionById(ac.IdCriterion)
}

func (ds *Datastore) GetAuditCriterionByIds(idAudit, idCriterion int64) (*AuditCriterion, error) {

	const sql = `SELECT ` +
		`id_audit, id_criterion, value, observation ` +
		`FROM places4all.audit_criterion ` +
		`WHERE id_audit = $1 AND id_criterion = $2`

	ac := AAuditCriterion(true)
	ac.SetExists()

	err := ds.postgres.QueryRowx(sql, idAudit, idCriterion).StructScan(&ac)
	if err != nil {
		return nil, err
	}

	return &ac, err
}

func (ds *Datastore) GetAuditCriteria(filter map[string]interface{}) ([]*AuditCriterion, error) {

	where, values := generators.GenerateAndSearchClause(filter)

	sql := ds.postgres.Rebind(`SELECT id_audit, id_criterion, value, observation ` +
		`FROM places4all.audit_criterion ` +
		where)

	rows, err := ds.postgres.Queryx(sql, values...)
	if err != nil {
		return nil, err
	}

	auditCriteria := make([]*AuditCriterion, 0)
	for rows.Next() {
		auditCriterion := NewAuditCriterion(false)
		auditCriterion.SetExists()
		err := rows.StructScan(auditCriterion)
		if err != nil {
			return nil, err
		}
		auditCriteria = append(auditCriteria, auditCriterion)
		if err != nil {
			return nil, err
		}
	}

	return auditCriteria, err
}
