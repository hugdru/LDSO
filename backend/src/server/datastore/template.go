package datastore

import (
	"errors"
	"gopkg.in/guregu/null.v3/zero"
	"server/datastore/generators"
	"server/datastore/metadata"
	"strconv"
	"time"
)

type Template struct {
	Id          int64       `json:"id" db:"id"`
	Name        string      `json:"name" db:"name"`
	Description zero.String `json:"description" db:"description"`
	CreatedDate time.Time   `json:"createdDate" db:"created_date"`

	// Objects
	Maingroups []*Maingroup `json:"maingroups,omitempty"`

	meta metadata.Metadata
}

func (t *Template) SetExists() {
	t.meta.Exists = true
}

func (t *Template) SetDeleted() {
	t.meta.Deleted = true
}

func (t *Template) Exists() bool {
	return t.meta.Exists
}

func (t *Template) Deleted() bool {
	return t.meta.Deleted
}

func (t *Template) MustSet(name string) error {

	if name != "" {
		t.Name = name
	} else {
		return errors.New("name must be set")
	}

	return nil
}

func (t *Template) AllSetIfNotEmptyOrNil(name, description string) error {
	if name != "" {
		t.Name = name
	}

	return t.OptionalSetIfNotEmptyOrNil(description)
}

func (t *Template) OptionalSetIfNotEmptyOrNil(description string) error {
	if description != "" {
		t.Description = zero.StringFrom(description)
	}

	return nil
}

func ATemplate(allocateObjects bool) Template {
	template := Template{}
	if allocateObjects {
		template.Maingroups = make([]*Maingroup, 0)
	}
	return template
}

func NewTemplate(allocateObjects bool) *Template {
	template := ATemplate(allocateObjects)
	return &template
}

func (ds *Datastore) InsertTemplate(t *Template) error {

	if t == nil {
		return errors.New("template should not be nil")
	}

	if t.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.template (` +
		`name, description, created_date` +
		`) VALUES (` +
		`$1, $2, $3` +
		`) RETURNING id`
	err := ds.postgres.QueryRow(sql, t.Name, t.Description, t.CreatedDate).Scan(&t.Id)
	if err != nil {
		return err
	}

	t.SetExists()

	return err
}

func (ds *Datastore) UpdateTemplate(t *Template) error {

	if t == nil {
		return errors.New("template should not be nil")
	}

	if !t.Exists() {
		return errors.New("update failed: does not exist")
	}

	if t.Deleted() {
		return errors.New("update failed: marked for deletion")
	}

	const sql = `UPDATE places4all.template SET (` +
		`name, description, created_date` +
		`) = ( ` +
		`$1, $2, $3` +
		`) WHERE id = $4`

	_, err := ds.postgres.Exec(sql, t.Name, t.Description, t.CreatedDate, t.Id)
	return err
}

func (ds *Datastore) SaveTemplate(t *Template) error {
	if t.Exists() {
		return ds.UpdateTemplate(t)
	}

	return ds.InsertTemplate(t)
}

func (ds *Datastore) UpsertTemplate(t *Template) error {

	if t == nil {
		return errors.New("template should not be nil")
	}

	if t.Exists() {
		return errors.New("insert failed: already exists")
	}

	const sql = `INSERT INTO places4all.template (` +
		`id, name, description, created_date` +
		`) VALUES (` +
		`$1, $2, $3, $4` +
		`) ON CONFLICT (id) DO UPDATE SET (` +
		`id, name, description, created_date` +
		`) = (` +
		`EXCLUDED.id, EXCLUDED.name, EXCLUDED.description, EXCLUDED.created_date` +
		`)`

	_, err := ds.postgres.Exec(sql, t.Id, t.Name, t.Description, t.CreatedDate)
	if err != nil {
		return err
	}

	t.SetExists()

	return err
}

func (ds *Datastore) DeleteTemplate(t *Template) error {

	if t == nil {
		return errors.New("template should not be nil")
	}

	if !t.Exists() {
		return nil
	}

	if t.Deleted() {
		return nil
	}

	const sql = `DELETE FROM places4all.template WHERE id = $1`

	_, err := ds.postgres.Exec(sql, t.Id)
	if err != nil {
		return err
	}

	t.SetDeleted()

	return err
}

func (ds *Datastore) DeleteTemplateById(id int64) error {

	const sql = `DELETE FROM places4all.template WHERE id = $1`

	_, err := ds.postgres.Exec(sql, id)
	if err != nil {
		return err
	}

	return err
}

func (ds *Datastore) GetTemplateById(id int64) (*Template, error) {

	const sql = `SELECT ` +
		`id, name, description, created_date ` +
		`FROM places4all.template ` +
		`WHERE id = $1`

	t := ATemplate(false)
	t.SetExists()

	err := ds.postgres.QueryRow(sql, id).Scan(&t.Id, &t.Name, &t.Description, &t.CreatedDate)
	if err != nil {
		return nil, err
	}

	return &t, err
}

func (ds *Datastore) GetTemplatesWithMaingroups(limit, offset int, filter map[string]interface{}) ([]*Template, error) {

	where, values := generators.GenerateAndSearchClause(filter)

	sql := `SELECT template.id, template.name, template.description, template.created_date ` +
		`FROM places4all.template ` +
		where +
		`ORDER BY template.id DESC LIMIT ` + strconv.Itoa(limit) +
		` OFFSET ` + strconv.Itoa(offset)
	sql = ds.postgres.Rebind(sql)

	rows, err := ds.postgres.Queryx(sql, values...)
	if err != nil {
		return nil, err
	}

	templates := make([]*Template, 0)
	for rows.Next() {
		template := NewTemplate(false)
		err := rows.StructScan(template)
		if err != nil {
			return nil, err
		}
		templates = append(templates, template)
		template.Maingroups, err = ds.GetMaingroupsByTemplateIdWithSubgroups(template.Id)
		if err != nil {
			return nil, err
		}
	}

	return templates, err
}

func (ds *Datastore) GetTemplateWithMaingroups(id int64) (*Template, error) {
	rows := ds.postgres.QueryRowx(
		`SELECT template.id, template.name, template.description, template.created_date `+
			`FROM places4all.template `+
			`WHERE id = $1`, id)

	template := NewTemplate(false)
	err := rows.StructScan(template)
	if err != nil {
		return nil, err
	}
	template.Maingroups, err = ds.GetMaingroupsByTemplateIdWithSubgroups(template.Id)
	if err != nil {
		return nil, err
	}

	return template, err
}
