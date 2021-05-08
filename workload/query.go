package workload

import "strings"

// QuerySet struct contains map of proxied queries
type QuerySet struct {
	Queries []string
}

var QSet *QuerySet

// QuerySet.Append appends proxied query to queries array
func (qs *QuerySet) Append(q string) {
	// checks if query is not empty
	if q == "" {
		return
	}
	// remove redundant spaces characters
	q = strings.TrimSpace(q)
	qs.Queries = append(qs.Queries, q)
}

// ResetQuerySet resets query set to empty map
func ResetQuerySet() *QuerySet {
	QSet = &QuerySet{Queries: []string{}}
	return QSet
}
