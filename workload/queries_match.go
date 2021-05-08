package workload

import (
	"fmt"
	"regexp"
	"strings"
)

// makeRegexpFromStatsQueries create regexp.Regexp from stat query pattern
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

// makeRegexpPattern creates regexp.Regexp from stat query pattern
// by replacing all $\d+ signs with .+
func makeRegexpPattern(q string) string {
	// escape all excape symbols and unescape $
	q = regexp.QuoteMeta(q)
	q = strings.Replace(q, "\\$", "$", -1)

	// replace all $\d+ signs with .+
	regexp := regexp.MustCompile("\\$\\d+")
	rp := regexp.ReplaceAllString(q, ".+")

	// pattern must be full matched
	return fmt.Sprintf("^%s$", rp)
}

// filterQueriesByRegexp filters proxied queries by stat query pattern
// returns proxied queries mathing patterns and all patterns with which
// proxied queries were matched (filtering 3rd party libs queries)
func filterQueriesByRegexp(rStatQueries map[*Query]*regexp.Regexp, qs *QuerySet) ([]string, []*Query) {
	filteredQueries := []string{}
	statQueries := map[*Query]bool{}

	for sq, regexq := range rStatQueries {
		for _, q := range qs.Queries {
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
