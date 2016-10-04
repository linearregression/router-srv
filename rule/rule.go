package rule

import (
	"math/rand"
	"sort"
	"time"

	proto "github.com/micro/go-os/router/proto"
	"github.com/micro/router-srv/db"
	rule "github.com/micro/router-srv/proto/rule"
)

// Using rule package as rules may be cached later on

type Rules []*rule.RuleSet

func (l Rules) Len() int {
	return len(l)
}

func (l Rules) Less(i, j int) bool {
	return l[i].Priority > l[j].Priority
}

func (l Rules) Swap(i, j int) {
	l[i], l[j] = l[j], l[i]
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

func Create(l *rule.RuleSet) error {
	return db.CreateRule(l)
}

func Update(l *rule.RuleSet) error {
	return db.UpdateRule(l)
}

func Delete(id string) error {
	return db.DeleteRule(id)
}

func Read(id string) (*rule.RuleSet, error) {
	return db.ReadRule(id)
}

func Search(service, version string, limit, offset int64) ([]*rule.RuleSet, error) {
	return db.SearchRule(service, version, limit, offset)
}

// Apply rule to a service.
func Apply(rules []*rule.RuleSet, s *proto.Service, labels map[string]string) *proto.Service {
	// no rules
	if len(rules) == 0 {
		return s
	}

	// sort low to high
	sortedRules := Rules(rules)
	sort.Sort(sortedRules)

	// remove flag we need at the end
	var remove bool
	var matched bool

	nodes := map[string]*proto.Node{}

RULE:
	for _, rule := range sortedRules {
		// same version? if not, skip
		if rule.Version != s.Version {
			continue
		}

		var isLabelRule bool

		if len(rule.Key) > 0 {
			isLabelRule = true
		}

		// is it a label rule? do we have labels to match?
		if isLabelRule {
			// no labels but a rule key, skip
			if len(labels) == 0 {
				continue
			}

			var seen bool
			for k, v := range labels {
				// TODO: if rule.Value is nil then glob match
				if k == rule.Key && v == rule.Value {
					seen = true
					break
				}
			}
			// no matching labels to the rule, skip
			if !seen {
				continue
			}
		}

		// it's either a label rule that matched or a regular rule
		matched = true

		// if weight is 0 and there is no label rule
		// delete the service entirely
		if rule.Weight == 0 && !isLabelRule {
			remove = true
			nodes = nil
			continue RULE
		}

		if rule.Weight == 100 {
			switch isLabelRule {
			case true:
				lv, ok := s.Metadata[rule.Key]
				if ok && rule.Value == lv {
					remove = false
				}
			default:
				remove = false
			}
		}

		if nodes == nil {
			nodes = map[string]*proto.Node{}
		}

		// Apply to nodes
	NODE:
		for _, node := range s.Nodes {
			// create the sampling right here
			r := rand.Int63n(100) <= rule.Weight

			// there is no label rule
			// generic service level rule always
			// overrides the label rule
			if !isLabelRule {
				if r {
					nodes[node.Id] = node
				} else {
					delete(nodes, node.Id)
				}
				continue NODE
			}

			// hasLabelRule

			lv, ok := node.Metadata[rule.Key]
			// if matches then sample
			if ok && lv == rule.Value {
				if r {
					nodes[node.Id] = node
				}
				continue NODE
			}
		}

		// flip remove if there's no nodes left
		if len(nodes) == 0 {
			remove = true
		} else {
			remove = false
		}
	}

	// explicitly told to remove service
	if remove {
		return nil
	}

	// we didn't actually match any rules
	if !matched {
		return s
	}

	// matched but told to keep service
	// rebuild the node list based on the rules applied
	var snodes []*proto.Node
	for _, node := range nodes {
		snodes = append(snodes, node)
	}
	s.Nodes = snodes

	return s
}
