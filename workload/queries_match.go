package workload

import (
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
	regexp := regexp.MustCompile("\\$\\d+")
	rp := regexp.ReplaceAllString(q, "?")
	replaces := map[string]string{
		"?": ".+",
		"(": "\\(",
		")": "\\)",
	}
	for r, replace := range replaces {
		rp = strings.Replace(rp, r, replace, -1)
	}
	return rp
}

func filterQueriesByRegexp(regexpStatQueries map[*Query]*regexp.Regexp, qs *query.QuerySet) ([]string, []*Query) {
	filteredQueries := []string{}
	statQueries := map[*Query]bool{}

	for sq, regexq := range regexpStatQueries {
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
