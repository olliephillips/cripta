package internal

import (
	"github.com/google/uuid"
)

var (
	// config
	CONFIG_FILE = "config.env"

	// mqtt
	SERVER = "test.mosquitto.org"
	PORT   = "1883"

	// folders
	MAILBOX_FOLDER     = "mailbox"
	GROUPS_FOLDER      = "groups"
	FRIEND_KEYS_FOLDER = "friend_keys"

	// file names
	PUBLIC_KEY_FILE  = "my_public_key.txt"
	PRIVATE_KEY_FILE = ".my_private_key"
)

// ShortUID provides a 5 char ref which should be unique enough
func ShortUID() string {
	id := uuid.New()
	r := []rune(id.String())
	start := len(r) - 5
	return string(r[start:])
}
