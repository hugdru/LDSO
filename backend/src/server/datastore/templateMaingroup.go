package datastore

import (
	"errors"
	"server/datastore/metadata"
)

type TemplateMaingroup struct {
	IdTemplate  int64 `json:"idTemplate" db:"id_template"`
	IdMaingroup int64 `json:"idMaingroup" db:"id_maingroup"`
	meta        metadata.Metadata
}

func (tm *TemplateMaingroup) SetExists() {
	tm.meta.Exists = true
}

func (tm *TemplateMaingroup) SetDeleted() {
	tm.meta.Deleted = true
}

func (tm *TemplateMaingroup) Exists() bool {
	return tm.meta.Exists
}

func (tm *TemplateMaingroup) Deleted() bool {
	return tm.meta.Deleted
}

func (ds *Datastore) Insert(tm *TemplateMaingroup) error {
	var err error

	if tm.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.template_maingroup (` +
		`id_template` +
		`) VALUES (` +
		`$1` +
		`) RETURNING id_maingroup`

	err = ds.postgres.QueryRow(sql, tm.IdTemplate).Scan(&tm.IdMaingroup)
	if err != nil {
		return err
	}

	tm.SetExists()

	return nil
}

func (ds *Datastore) Update(tm *TemplateMaingroup) error {
	var err error

	if !tm.Exists() {
		return errors.New("update failed: does not exist")
	}

	if tm.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.template_maingroup SET (` +
		`id_template` +
		`) = ( ` +
		`$1` +
		`) WHERE id_maingroup = $2`

	_, err = ds.postgres.Exec(sql, tm.IdTemplate, tm.IdMaingroup)
	return err
}

func (ds *Datastore) Save(tm *TemplateMaingroup) error {
	if tm.Exists() {
		return ds.Update(tm)
	}

	return ds.Insert(tm)
}

func (ds *Datastore) Upsert(tm *TemplateMaingroup) error {
	var err error

	if tm.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.template_maingroup (` +
		`id_template, id_maingroup` +
		`) VALUES (` +
		`$1, $2` +
		`) ON CONFLICT (id_maingroup) DO UPDATE SET (` +
		`id_template, id_maingroup` +
		`) = (` +
		`EXCLUDED.id_template, EXCLUDED.id_maingroup` +
		`)`

	_, err = ds.postgres.Exec(sql, tm.IdTemplate, tm.IdMaingroup)
	if err != nil {
		return err
	}

	tm.SetExists()

	return nil
}

func (ds *Datastore) Delete(tm *TemplateMaingroup) error {
	var err error

	if !tm.Exists() {
		return nil
	}

	if tm.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.template_maingroup WHERE id_maingroup = $1`

	_, err = ds.postgres.Exec(sql, tm.IdMaingroup)
	if err != nil {
		return err
	}

	tm.SetDeleted()

	return nil
}

func (ds *Datastore) Maingroup(tm *TemplateMaingroup) (*Maingroup, error) {
	return ds.GetMaingroupById(tm.IdMaingroup)
}

func (ds *Datastore) Template(tm *TemplateMaingroup) (*Template, error) {
	return ds.GetTemplateById(tm.IdTemplate)
}

func (ds *Datastore) TemplateMaingroupByIDTemplateIDMaingroup(idTemplate, idMaingroup int64) (*TemplateMaingroup, error) {
	var err error

	const sql = `SELECT ` +
		`id_template, id_maingroup ` +
		`FROM places4all.template_maingroup ` +
		`WHERE id_template = $1 AND id_maingroup = $2`

	tm := TemplateMaingroup{}
	tm.SetExists()

	err = ds.postgres.QueryRow(sql, idTemplate, idMaingroup).Scan(&tm.IdTemplate, &tm.IdMaingroup)
	if err != nil {
		return nil, err
	}

	return &tm, nil
}
