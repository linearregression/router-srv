package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/micro/router-srv/db"
	"github.com/micro/router-srv/proto/label"
)

var (
	Url = "root@tcp(127.0.0.1:3306)/router"

	database string

	q = map[string]string{}

	// account queries
	labelQ = map[string]string{
		"deleteLabel":              "DELETE from %s.%s where id = ? limit 1",
		"createLabel":              `INSERT into %s.%s (id, service, version, weight, priority, label, value) values (?, ?, ?, ?, ?, ?, ?)`,
		"updateLabel":              "UPDATE %s.%s set service = ?, version = ?, weight = ?, priority = ?, label = ?, value = ? where id = ?",
		"readLabel":                "SELECT id, service, version, weight, priority, label, value from %s.%s where id = ? limit 1",
		"searchLabel":              "SELECT id, service, version, weight, priority, label, value from %s.%s limit ? offset ?",
		"searchLabelService":       "SELECT id, service, version, weight, priority, label, value from %s.%s where service = ? limit ? offset ?",
		"searchLabelKey":           "SELECT id, service, version, weight, priority, label, value from %s.%s where label = ? limit ? offset ?",
		"searchLabelServiceAndKey": "SELECT id, service, version, weight, priority, label, value from %s.%s  where service = ? and label = ? limit ? offset ?",
	}

	st = map[string]*sql.Stmt{}
)

type mysql struct {
	db *sql.DB
}

func init() {
	db.Register(new(mysql))
}

func (m *mysql) Init() error {
	var d *sql.DB
	var err error

	parts := strings.Split(Url, "/")
	if len(parts) != 2 {
		return errors.New("Invalid database url")
	}

	if len(parts[1]) == 0 {
		return errors.New("Invalid database name")
	}

	url := parts[0]
	database := parts[1]

	if d, err = sql.Open("mysql", url+"/"); err != nil {
		return err
	}
	if _, err := d.Exec("CREATE DATABASE IF NOT EXISTS " + database); err != nil {
		return err
	}
	d.Close()
	if d, err = sql.Open("mysql", Url); err != nil {
		return err
	}
	if _, err = d.Exec(labelSchema); err != nil {
		return err
	}

	for query, statement := range labelQ {
		prepared, err := d.Prepare(fmt.Sprintf(statement, database, "labels"))
		if err != nil {
			return err
		}
		st[query] = prepared
	}

	m.db = d

	return nil
}

func (m *mysql) DeleteLabel(id string) error {
	_, err := st["deleteLabel"].Exec(id)
	return err
}

func (m *mysql) CreateLabel(l *label.LabelSet) error {
	_, err := st["createLabel"].Exec(l.Id, l.Service, l.Version, l.Weight, l.Priority, l.Key, l.Value)
	return err
}

func (m *mysql) UpdateLabel(l *label.LabelSet) error {
	_, err := st["updateLabel"].Exec(l.Service, l.Version, l.Weight, l.Priority, l.Key, l.Value, l.Id)
	return err
}

func (m *mysql) ReadLabel(id string) (*label.LabelSet, error) {
	l := &label.LabelSet{}

	r := st["readLabel"].QueryRow(id)
	// we dont return salt or secret
	if err := r.Scan(&l.Id, &l.Service, &l.Version, &l.Weight, &l.Priority, &l.Key, &l.Value); err != nil {
		if err == sql.ErrNoRows {
			return nil, db.ErrNotFound
		}
		return nil, err
	}

	return l, nil
}

func (m *mysql) SearchLabel(service, key string, limit, offset int64) ([]*label.LabelSet, error) {
	var r *sql.Rows
	var err error

	if len(service) > 0 && len(key) > 0 {
		r, err = st["searchLabelServiceAndKey"].Query(service, key, limit, offset)
	} else if len(service) > 0 {
		r, err = st["searchLabelService"].Query(service, limit, offset)
	} else if len(key) > 0 {
		r, err = st["searchLabelKey"].Query(key, limit, offset)
	} else {
		r, err = st["searchLabel"].Query(limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer r.Close()

	var labels []*label.LabelSet

	for r.Next() {
		l := &label.LabelSet{}
		if err := r.Scan(&l.Id, &l.Service, &l.Version, &l.Weight, &l.Priority, &l.Key, &l.Value); err != nil {
			if err == sql.ErrNoRows {
				return nil, db.ErrNotFound
			}
			return nil, err
		}

		labels = append(labels, l)

	}
	if r.Err() != nil {
		return nil, err
	}

	return labels, nil
}
