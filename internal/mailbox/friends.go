package mailbox

import (
	"fmt"
	"github.com/olliephillips/cripta/internal"
	"io/ioutil"
	"strings"
)

// ListFriends returns a slice of strings representing public key filenames held for friends
func ListFriends() ([]string, error) {
	var friends []string
	files, err := ioutil.ReadDir(internal.FRIEND_KEYS_FOLDER)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		friends = append(friends, fmt.Sprintf("@%s", strings.TrimSuffix(file.Name(), ".txt")))
	}

	return friends, nil
}
