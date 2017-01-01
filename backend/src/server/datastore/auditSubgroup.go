package datastore

import (
	"errors"
	"server/datastore/generators"
	"server/datastore/metadata"
)

type AuditSubgroup struct {
	IdAudit    int64 `json:"idAudit" db:"id_audit"`
	IdSubgroup int64 `json:"idSubgroup" db:"id_subgroup"`

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

func (ds *Datastore) InsertAuditSubgroup(as *AuditSubgroup) error {

	if as == nil {
		return errors.New("auditorSubgroup should not be nil")
	}

	if as.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.audit_subgroup (` +
		`id_audit, id_subgroup` +
		`) VALUES (` +
		`$1, $2` +
		`)`

	_, err := ds.postgres.Exec(sql, as.IdAudit, as.IdSubgroup)
	if err != nil {
		return err
	}

	as.SetExists()

	return err
}

func (ds *Datastore) DeleteAuditSubgroup(as *AuditSubgroup) error {

	if as == nil {
		return errors.New("auditorSubgroup should not be nil")
	}

	if !as.Exists() {
		return nil
	}

	if as.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.audit_subgroup WHERE id_audit = $1 AND id_subgroup = $2`

	_, err := ds.postgres.Exec(sql, as.IdAudit, as.IdSubgroup)
	if err != nil {
		return err
	}

	as.SetDeleted()

	return err
}

func (ds *Datastore) GetAuditSubgroupAudit(as *AuditSubgroup) (*Audit, error) {
	return ds.GetAuditById(as.IdAudit)
}

func (ds *Datastore) GetAuditSubgroupSubgroup(as *AuditSubgroup) (*Subgroup, error) {
	return ds.GetSubgroupById(as.IdSubgroup)
}

func (ds *Datastore) GetAuditSubgroupById(idAudit, idSubgroup int64) (*AuditSubgroup, error) {
	var err error

	const sql = `SELECT ` +
		`id_audit, id_subgroup ` +
		`FROM places4all.audit_subgroup ` +
		`WHERE id_audit = $1 AND id_subgroup = $2`

	as := AAuditSubgroup(true)
	as.SetExists()

	err = ds.postgres.QueryRowx(sql, idAudit, idSubgroup).StructScan(&as)
	if err != nil {
		return nil, err
	}

	return &as, err
}

func (ds *Datastore) GetAuditSubgroupsByIdAudit(idAudit int64, filter map[string]interface{}) ([]int64, error) {

	where, values := generators.GenerateAndSearchClause(filter)

	sql := ds.postgres.Rebind(`SELECT ` +
		`id_subgroup ` +
		`FROM places4all.audit_subgroup ` +
		where +
		`WHERE id_audit = $1`)

	rows, err := ds.postgres.Queryx(sql, append(values, idAudit)...)
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

func (ds *Datastore) SaveAuditSubgroup(idAudit int64, idTemplate int64, idsSubgroups []int64) (def error) {

	if idsSubgroups == nil {
		return errors.New("idsSubgroups should not be nil")
	}

	in, values := generators.GenerateIn(idsSubgroups)
	if in == "" || values == nil {
		return errors.New("GenerateIn failed")
	}

	sqlCheck := ds.postgres.Rebind(`SELECT count(*) as count ` +
		`FROM places4all.subgroup ` +
		`JOIN places4all.maingroup ON maingroup.id = subgroup.id_maingroup ` +
		`JOIN places4all.template ON template.id = maingroup.id_template ` +
		`WHERE subgroup.id IN (` + in + `)`)
	var count int
	err := ds.postgres.QueryRow(sqlCheck, values...).Scan(&count)
	if err != nil {
		return err
	}
	if count != len(idsSubgroups) {
		return errors.New("Some subgroups don't belong to template")
	}

	const sqlInsert = `INSERT INTO places4all.audit_subgroup (` +
		`id_audit, id_subgroup` +
		`) VALUES (` +
		`$1, $2` +
		`)`

	tx, err := ds.postgres.Begin()
	if err != nil {
		return err
	}
	defer func() {
		if def != nil {
			tx.Rollback()
			return
		}
		def = tx.Commit()
	}()

	for _, idSubgroup := range idsSubgroups {
		_, err = tx.Exec(sqlInsert, idAudit, idSubgroup)
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
