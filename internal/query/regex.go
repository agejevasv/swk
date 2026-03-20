package query

import (
	"fmt"
	"regexp"
)

// Match represents a single regex match.
type Match struct {
	Value  string   `json:"value"`
	Start  int      `json:"start"`
	End    int      `json:"end"`
	Groups []string `json:"groups,omitempty"`
}

// RegexResult holds the result of a regex test.
type RegexResult struct {
	Matched bool    `json:"matched"`
	Matches []Match `json:"matches"`
}

func RegexTest(input, pattern string, global bool) (*RegexResult, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return nil, fmt.Errorf("invalid regex pattern: %w", err)
	}

	result := &RegexResult{}

	var allLocs [][]int
	if global {
		allLocs = re.FindAllStringSubmatchIndex(input, -1)
	} else {
		if loc := re.FindStringSubmatchIndex(input); loc != nil {
			allLocs = [][]int{loc}
		}
	}

	if len(allLocs) == 0 {
		return result, nil
	}

	result.Matched = true
	for _, loc := range allLocs {
		result.Matches = append(result.Matches, buildMatch(input, loc))
	}

	return result, nil
}

func buildMatch(input string, loc []int) Match {
	m := Match{
		Value: input[loc[0]:loc[1]],
		Start: loc[0],
		End:   loc[1],
	}
	for i := 2; i < len(loc); i += 2 {
		if loc[i] >= 0 {
			m.Groups = append(m.Groups, input[loc[i]:loc[i+1]])
		} else {
			m.Groups = append(m.Groups, "")
		}
	}
	return m
}

func RegexReplace(input, pattern, replacement string) (string, error) {
	re, err := regexp.Compile(pattern)
	if err != nil {
		return "", fmt.Errorf("invalid regex pattern: %w", err)
	}
	return re.ReplaceAllString(input, replacement), nil
}
