package workload

import (
	"fmt"
	"pgworkload/query"
	"regexp"
	"strings"
)

func makeRegexpFromStatsQueries(queries []*Query) map[*Query]*regexp.Regexp {
	regMap := map[*Query]*regexp.Regexp{}
	for _, q := range queries {
		rp := makeRegexpPattern(q.Query)
		if _, inmap := regMap[q]; !inmap {
			queryRegex := regexp.MustCompile(rp)
			regMap[q] = queryRegex
		}
	}

	return regMap
}

func makeRegexpPattern(q string) string {
	q = regexp.QuoteMeta(q)
	q = strings.Replace(q, "\\$", "$", -1)

	regexp := regexp.MustCompile("\\$\\d+")
	rp := regexp.ReplaceAllString(q, "?")
	rp = strings.Replace(rp, "?", ".+", -1)

	return fmt.Sprintf("^%s$", rp)
}

func filterQueriesByRegexp(rStatQueries map[*Query]*regexp.Regexp, qs *query.QuerySet) ([]string, []*Query) {
	filteredQueries := []string{}
	statQueries := map[*Query]bool{}

	for sq, regexq := range rStatQueries {
		for q := range qs.Queries {
			if regexq.MatchString(q) {
				filteredQueries = append(filteredQueries, q)
				if _, inmap := statQueries[sq]; !inmap {
					statQueries[sq] = true
				}
			}
		}
	}

	matchedStatQueries := []*Query{}
	for sq := range statQueries {
		matchedStatQueries = append(matchedStatQueries, sq)
	}

	return filteredQueries, matchedStatQueries
}
