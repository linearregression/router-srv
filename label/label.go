package label

import (
	"github.com/micro/router-srv/db"
	"github.com/micro/router-srv/proto/label"
)

// Using label package as labels may be cached later on

func Create(l *label.LabelSet) error {
	return db.CreateLabel(l)
}

func Update(l *label.LabelSet) error {
	return db.UpdateLabel(l)
}

func Delete(id string) error {
	return db.DeleteLabel(id)
}

func Read(id string) (*label.LabelSet, error) {
	return db.ReadLabel(id)
}

func Search(service, key string, limit, offset int64) ([]*label.LabelSet, error) {
	return db.SearchLabel(service, key, limit, offset)
}
