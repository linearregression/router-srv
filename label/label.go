package label

import (
	"math/rand"
	"sort"
	"time"

	"github.com/micro/go-micro/registry"
	"github.com/micro/router-srv/db"
	"github.com/micro/router-srv/proto/label"
)

// Using label package as labels may be cached later on

type Labels []*label.LabelSet

func (l Labels) Len() int {
	return len(l)
}

func (l Labels) Less(i, j int) bool {
	return l[i].Priority > l[j].Priority
}

func (l Labels) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

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

func Apply(labels []*label.LabelSet, s *registry.Service) {
	// sort low to high
	sortedLabels := Labels(labels)

	sort.Sort(sortedLabels)

	for _, label := range sortedLabels {
		// does the label have a version and does it match?
		if len(label.Version) > 0 && label.Version != s.Version {
			continue
		}

		// ok so either the version matches or its a generic apply

		// delete they label
		if label.Weight == 0 {
			v := s.Metadata[label.Key]
			// if there's a key value and its the same
			// or if the key value is blank
			if i := len(label.Value); i > 0 && v == label.Value || i == 0 {
				delete(s.Metadata, label.Key)
			}

			// delete from the nodes
			for _, node := range s.Nodes {
				v := node.Metadata[label.Key]
				if i := len(label.Value); i > 0 && v == label.Value || i == 0 {
					delete(node.Metadata, label.Key)
				}
			}
			// other weight is greater than 0, create label
		} else {
			// Apply at the top level
			// Should we actually apply as a sampling here?
			// Or at the individual node basis?
			if label.Weight == 100 {
				if s.Metadata == nil {
					s.Metadata = make(map[string]string)
				}
				s.Metadata[label.Key] = label.Value
			}

			// Apply to nodes
			for _, node := range s.Nodes {
				// Apply based on the weight
				if rand.Int63n(100) <= label.Weight {
					if node.Metadata == nil {
						node.Metadata = make(map[string]string)
					}
					node.Metadata[label.Key] = label.Value
				}
			}
		}
	}
}
