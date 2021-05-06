package query

type QuerySet struct {
	Queries map[string]bool
}

var QSet *QuerySet

func (qs *QuerySet) Append(q string) {
	if q == "" {
		return
	}

	_, inset := qs.Queries[q]
	if !inset {
		qs.Queries[q] = true
	}
}

func (qs *QuerySet) Contains(q string) bool {
	return qs.Queries[q]
}

func ResetQuerySet() *QuerySet {
	QSet = &QuerySet{Queries: map[string]bool{}}
	return QSet
}
