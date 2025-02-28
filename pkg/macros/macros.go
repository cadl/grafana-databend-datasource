package macros

import (
	"errors"
	"fmt"
	"math"
	"strings"

	"github.com/grafana/sqlds/v2"
)

var (
	ErrorNoArgumentsToMacro           = errors.New("expected minimum of 1 argument. But no argument found")
	ErrorInsufficientArgumentsToMacro = errors.New("expected number of arguments not matching")
)

type timeQueryType string

const (
	timeQueryTypeFrom timeQueryType = "from"
	timeQueryTypeTo   timeQueryType = "to"
)

func newTimeFilter(queryType timeQueryType, query *sqlds.Query) (string, error) {
	date := query.TimeRange.From
	if queryType == timeQueryTypeTo {
		date = query.TimeRange.To
	}
	return fmt.Sprintf("TO_TIMESTAMP(%d)", date.Unix()), nil
}

// FromTimeFilter return time filter query based on grafana's timepicker's from time
func FromTimeFilter(query *sqlds.Query, args []string) (string, error) {
	return newTimeFilter(timeQueryTypeFrom, query)
}

// ToTimeFilter return time filter query based on grafana's timepicker's to time
func ToTimeFilter(query *sqlds.Query, args []string) (string, error) {
	return newTimeFilter(timeQueryTypeTo, query)
}

func TimeFilter(query *sqlds.Query, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("%w: expected 1 argument, received %d", sqlds.ErrorBadArgumentCount, len(args))
	}

	var (
		column = args[0]
		from   = query.TimeRange.From.UTC().Unix()
		to     = query.TimeRange.To.UTC().Unix()
	)

	return fmt.Sprintf("%s >= '%d' AND %s <= '%d'", column, from, column, to), nil
}

func DateFilter(query *sqlds.Query, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("%w: expected 1 argument, received %d", sqlds.ErrorBadArgumentCount, len(args))
	}
	var (
		column = args[0]
		from   = query.TimeRange.From.Format("2006-01-02")
		to     = query.TimeRange.To.Format("2006-01-02")
	)

	return fmt.Sprintf("%s >= '%s' AND %s <= '%s'", column, from, column, to), nil
}

func TimeFilterMs(query *sqlds.Query, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("%w: expected 1 argument, received %d", sqlds.ErrorBadArgumentCount, len(args))
	}

	var (
		column = args[0]
		from   = query.TimeRange.From.UTC().UnixMilli()
		to     = query.TimeRange.To.UTC().UnixMilli()
	)

	return fmt.Sprintf("%s >= '%d' AND %s <= '%d'", column, from, column, to), nil
}

func TimeInterval(query *sqlds.Query, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("%w: expected 1 argument, received %d", sqlds.ErrorBadArgumentCount, len(args))
	}

	seconds := math.Max(query.Interval.Seconds(), 1)
	return fmt.Sprintf("TO_TIMESTAMP( TO_UNIX_TIMESTAMP(TO_TIMESTAMP(%s)) // %d * %d)", args[0], int(seconds), int(seconds)), nil
}

func TimeIntervalMs(query *sqlds.Query, args []string) (string, error) {
	if len(args) != 1 {
		return "", fmt.Errorf("%w: expected 1 argument, received %d", sqlds.ErrorBadArgumentCount, len(args))
	}
	return TimeInterval(query, args)
}

func IntervalSeconds(query *sqlds.Query, args []string) (string, error) {
	seconds := math.Max(query.Interval.Seconds(), 1)
	return fmt.Sprintf("%d", int(seconds)), nil
}

// RemoveQuotesInArgs remove all quotes from macro arguments and return
func RemoveQuotesInArgs(args []string) []string {
	updatedArgs := []string{}
	for _, arg := range args {
		replacer := strings.NewReplacer(
			"\"", "",
			"'", "",
		)
		updatedArgs = append(updatedArgs, replacer.Replace(arg))
	}
	return updatedArgs
}

// IsValidComparisonPredicates checks for a string and return true if it is a valid SQL comparison predicate
func IsValidComparisonPredicates(comparison_predicates string) bool {
	switch comparison_predicates {
	case "=", "!=", "<>", "<", "<=", ">", ">=":
		return true
	}
	return false
}
