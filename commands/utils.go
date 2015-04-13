package commands

import (
	"fmt"

	"github.com/mattn/go-isatty"
	"github.com/wolfeidau/buildkite-cli/utils"
)

func toMap(columns []string, values []interface{}) map[string]interface{} {

	m := make(map[string]interface{})

	if len(columns) != len(values) {
		utils.Check(fmt.Errorf("Miss match in columns and values"))
	}

	for i, c := range columns {
		m[c] = values[i]
	}

	return m
}

func isTerminal(fd uintptr) bool {
	return isatty.IsTerminal(fd)
}
