package db

import (
	"errors"

	label "github.com/micro/router-srv/proto/label"
	rule "github.com/micro/router-srv/proto/rule"
)

type DB interface {
	Init() error
	Label
	Rule
}

type Label interface {
	ReadLabel(id string) (*label.LabelSet, error)
	DeleteLabel(id string) error
	CreateLabel(l *label.LabelSet) error
	UpdateLabel(l *label.LabelSet) error
	SearchLabel(service, key string, limit, offset int64) ([]*label.LabelSet, error)
}

type Rule interface {
	ReadRule(id string) (*rule.RuleSet, error)
	DeleteRule(id string) error
	CreateRule(l *rule.RuleSet) error
	UpdateRule(l *rule.RuleSet) error
	SearchRule(service, version string, limit, offset int64) ([]*rule.RuleSet, error)
}

var (
	db DB

	ErrNotFound = errors.New("not found")
)

func Register(backend DB) {
	db = backend
}

func Init() error {
	return db.Init()
}

func ReadLabel(id string) (*label.LabelSet, error) {
	return db.ReadLabel(id)
}

func CreateLabel(l *label.LabelSet) error {
	return db.CreateLabel(l)
}

func UpdateLabel(l *label.LabelSet) error {
	return db.UpdateLabel(l)
}

func DeleteLabel(id string) error {
	return db.DeleteLabel(id)
}

func SearchLabel(service, key string, limit, offset int64) ([]*label.LabelSet, error) {
	return db.SearchLabel(service, key, limit, offset)
}

func ReadRule(id string) (*rule.RuleSet, error) {
	return db.ReadRule(id)
}

func CreateRule(l *rule.RuleSet) error {
	return db.CreateRule(l)
}

func UpdateRule(l *rule.RuleSet) error {
	return db.UpdateRule(l)
}

func DeleteRule(id string) error {
	return db.DeleteRule(id)
}

func SearchRule(service, version string, limit, offset int64) ([]*rule.RuleSet, error) {
	return db.SearchRule(service, version, limit, offset)
}
