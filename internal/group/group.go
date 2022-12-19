package group

import (
	"github.com/olliephillips/cripta/internal"
	"io/ioutil"
	"strings"
)

// ListGroups returns a slice of strings
func ListGroups() ([]string, error) {
	var groups []string
	files, err := ioutil.ReadDir(internal.GROUPS_FOLDER)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		groups = append(groups, strings.TrimSuffix(file.Name(), ".txt"))
	}

	return groups, nil
}
