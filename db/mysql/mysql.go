package mysql

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"

	_ "github.com/go-sql-driver/mysql"
	"github.com/micro/router-srv/db"
	label "github.com/micro/router-srv/proto/label"
	rule "github.com/micro/router-srv/proto/rule"
)

var (
	Url = "root@tcp(127.0.0.1:3306)/router"

	database string

	q = map[string]string{}

	// label queries
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

	// rule queries
	ruleQ = map[string]string{
		"deleteRule":                  "DELETE from %s.%s where id = ? limit 1",
		"createRule":                  `INSERT into %s.%s (id, service, version, weight, priority, label, value) values (?, ?, ?, ?, ?, ?, ?)`,
		"updateRule":                  "UPDATE %s.%s set service = ?, version = ?, weight = ?, priority = ?, label = ?, value = ? where id = ?",
		"readRule":                    "SELECT id, service, version, weight, priority, label, value from %s.%s where id = ? limit 1",
		"searchRule":                  "SELECT id, service, version, weight, priority, label, value from %s.%s limit ? offset ?",
		"searchRuleService":           "SELECT id, service, version, weight, priority, label, value from %s.%s where service = ? limit ? offset ?",
		"searchRuleServiceAndVersion": "SELECT id, service, version, weight, priority, label, value from %s.%s  where service = ? and version = ? limit ? offset ?",
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
	if _, err = d.Exec(ruleSchema); err != nil {
		return err
	}

	for query, statement := range labelQ {
		prepared, err := d.Prepare(fmt.Sprintf(statement, database, "labels"))
		if err != nil {
			return err
		}
		st[query] = prepared
	}

	for query, statement := range ruleQ {
		prepared, err := d.Prepare(fmt.Sprintf(statement, database, "rules"))
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

func (m *mysql) DeleteRule(id string) error {
	_, err := st["deleteRule"].Exec(id)
	return err
}

func (m *mysql) CreateRule(r *rule.RuleSet) error {
	_, err := st["createRule"].Exec(r.Id, r.Service, r.Version, r.Weight, r.Priority, r.Key, r.Value)
	return err
}

func (m *mysql) UpdateRule(r *rule.RuleSet) error {
	_, err := st["updateRule"].Exec(r.Service, r.Version, r.Weight, r.Priority, r.Key, r.Value, r.Id)
	return err
}

func (m *mysql) ReadRule(id string) (*rule.RuleSet, error) {
	r := &rule.RuleSet{}

	row := st["readRule"].QueryRow(id)
	// we dont return salt or secret
	if err := row.Scan(&r.Id, &r.Service, &r.Version, &r.Weight, &r.Priority, &r.Key, &r.Value); err != nil {
		if err == sql.ErrNoRows {
			return nil, db.ErrNotFound
		}
		return nil, err
	}

	return r, nil
}

func (m *mysql) SearchRule(service, version string, limit, offset int64) ([]*rule.RuleSet, error) {
	var rows *sql.Rows
	var err error

	if len(service) > 0 && len(version) > 0 {
		rows, err = st["searchRuleServiceAndVersion"].Query(service, version, limit, offset)
	} else if len(service) > 0 {
		rows, err = st["searchRuleService"].Query(service, limit, offset)
	} else {
		rows, err = st["searchRule"].Query(limit, offset)
	}

	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var rules []*rule.RuleSet

	for rows.Next() {
		r := &rule.RuleSet{}
		if err := rows.Scan(&r.Id, &r.Service, &r.Version, &r.Weight, &r.Priority, &r.Key, &r.Value); err != nil {
			if err == sql.ErrNoRows {
				return nil, db.ErrNotFound
			}
			return nil, err
		}

		rules = append(rules, r)

	}
	if rows.Err() != nil {
		return nil, err
	}

	return rules, nil
}
