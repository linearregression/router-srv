package db

import (
	"errors"

	"github.com/micro/router-srv/proto/label"
)

type DB interface {
	Init() error
	Label
}

type Label interface {
	ReadLabel(id string) (*label.LabelSet, error)
	DeleteLabel(id string) error
	CreateLabel(l *label.LabelSet) error
	UpdateLabel(l *label.LabelSet) error
	SearchLabel(service, key string, limit, offset int64) ([]*label.LabelSet, error)
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
