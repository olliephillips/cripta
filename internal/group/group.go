package group

import (
	"bufio"
	"fmt"
	"github.com/olliephillips/cripta/internal"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
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

// ListMembers extracts a the list of memnbers for a group
func ListMembers(group string) ([]string, error) {
	var members []string
	groupFile, err := os.Open(filepath.Join(internal.GROUPS_FOLDER, fmt.Sprintf("%s.txt", group)))
	if err != nil {
		return nil, err
	}
	defer groupFile.Close()

	scanner := bufio.NewScanner(groupFile)
	for scanner.Scan() {
		ln := scanner.Text()
		if ln != "" {
			// empty check
			if !strings.Contains(ln, " ") {
				//spaces check
				if strings.HasPrefix(ln, "@") {
					// good
					members = append(members, ln)
				} else {
					// missing @ just add it
					members = append(members, fmt.Sprintf("@%s", ln))
				}
			} else {
				log.Println("group memmbers; skipping", ln)
			}
		}
	}

	return members, scanner.Err()
}
