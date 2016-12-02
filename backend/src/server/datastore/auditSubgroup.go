package datastore

import (
	"errors"
	"server/datastore/generators"
	"server/datastore/metadata"
)

type AuditSubgroup struct {
	IdAudit    int64 `json:"id_audit" db:"id_audit"`
	IdSubgroup int64 `json:"id_subgroup" db:"id_subgroup"`

	meta metadata.Metadata
}

func (ac *AuditSubgroup) SetExists() {
	ac.meta.Exists = true
}

func (ac *AuditSubgroup) SetDeleted() {
	ac.meta.Deleted = true
}

func (ac *AuditSubgroup) Exists() bool {
	return ac.meta.Exists
}

func (ac *AuditSubgroup) Deleted() bool {
	return ac.meta.Deleted
}

func AAuditSubgroup(allocateObjects bool) AuditSubgroup {
	auditSubgroup := AuditSubgroup{}
	//if allocateObjects {
	//}
	return auditSubgroup
}
func NewAuditSubgroup(allocateObjects bool) *AuditSubgroup {
	auditSubgroup := AAuditSubgroup(allocateObjects)
	return &auditSubgroup
}

func (ds *Datastore) InsertAuditSubgroup(ac *AuditSubgroup) error {

	if ac.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.audit_subgroup (` +
		`id_audit, id_subgroup` +
		`) VALUES (` +
		`$1, $2` +
		`)`

	_, err := ds.postgres.Exec(sql, ac.IdAudit, ac.IdSubgroup)
	if err != nil {
		return err
	}

	ac.SetExists()

	return err
}

func (ds *Datastore) DeleteAuditSubgroup(ac *AuditSubgroup) error {

	if !ac.Exists() {
		return nil
	}

	if ac.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.audit_subgroup WHERE id_audit = $1 AND id_subgroup = $2`

	_, err := ds.postgres.Exec(sql, ac.IdAudit, ac.IdSubgroup)
	if err != nil {
		return err
	}

	ac.SetDeleted()

	return err
}

func (ds *Datastore) GetAuditSubgroupAudit(ac *AuditSubgroup) (*Audit, error) {
	return ds.GetAuditById(ac.IdAudit)
}

func (ds *Datastore) GetAuditSubgroupSubgroup(ac *AuditSubgroup) (*Subgroup, error) {
	return ds.GetSubgroupById(ac.IdSubgroup)
}

func (ds *Datastore) GetAuditSubgroupById(idAudit, idSubgroup int64) (*AuditSubgroup, error) {
	var err error

	const sql = `SELECT ` +
		`id_audit, id_subgroup ` +
		`FROM places4all.audit_subgroup ` +
		`WHERE id_audit = $1 AND id_subgroup = $2`

	ac := AAuditSubgroup(true)
	ac.SetExists()

	err = ds.postgres.QueryRowx(sql, idAudit, idSubgroup).StructScan(&ac)
	if err != nil {
		return nil, err
	}

	return &ac, err
}

func (ds *Datastore) GetAuditSubgroupsByIdAudit(idAudit int64, filter map[string]string) ([]int64, error) {

	where, values := generators.GenerateSearchClause(filter)

	sql := `SELECT ` +
		`id_subgroup ` +
		`FROM places4all.audit_subgroup ` +
		where +
		`WHERE id_audit = $1`
	sql = ds.postgres.Rebind(sql)
	values = append(values, idAudit)

	rows, err := ds.postgres.Queryx(sql, values...)
	if err != nil {
		return nil, err
	}

	subgroups := make([]int64, 0)
	for rows.Next() {
		var idSubgroup int64
		err := rows.Scan(&idSubgroup)
		if err != nil {
			return nil, err
		}
		subgroups = append(subgroups, idSubgroup)
	}

	return subgroups, err
}

// TODO: check if subgroup is from the template audit is using
func (ds *Datastore) SaveAuditSubgroup(idAudit int64, idTemplate int64, idsSubgroups []int64) error {

	const sql = `INSERT INTO places4all.audit_subgroup (` +
		`id_audit, id_subgroup` +
		`) VALUES (` +
		`$1, $2` +
		`)`

	tx, err := ds.postgres.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
			return
		}
		err = tx.Commit()
	}()

	for _, idSubgroup := range idsSubgroups {
		_, err := tx.Exec(sql, idAudit, idSubgroup)
		if err != nil {
			return err
		}
	}

	return err
}

func (ds *Datastore) DeleteAuditSubgroupsByIdAudit(idAudit int64) error {
	const sql = `DELETE FROM places4all.audit_subgroup WHERE id_audit = $1`

	_, err := ds.postgres.Exec(sql, idAudit)
	if err != nil {
		return err
	}

	return err
}
