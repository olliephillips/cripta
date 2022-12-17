package mailbox

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"path/filepath"
)

// MailStorer so we can enable additional storage targets
type MailStorer interface {
	List() ([]Message, error)
	Get(string) (Message, error)
	Set(Message) error
	Delete(...string) error
}

// FSMailbox implements MailStorer providing a file system store
type FSMailbox struct {
	location string
}

// NewFSMailbox creates a new file system store
func NewFSMailbox(location string) *FSMailbox {
	return &FSMailbox{location}
}

// List returns all messages in the mailbox
func (mb *FSMailbox) List() ([]Message, error) {
	var msgs []Message
	files, err := ioutil.ReadDir(mb.location)
	if err != nil {
		return nil, err
	}

	for _, file := range files {
		var msg Message
		data, err := os.ReadFile(filepath.Join(mb.location, file.Name()))
		if err != nil {
			return nil, err
		}

		if err := json.Unmarshal(data, &msg); err != nil {
			return nil, err
		}

		msgs = append(msgs, msg)
	}

	return msgs, nil
}

// Get retrieves a single message from the
func (mb *FSMailbox) Get(msgId string) (Message, error) {

	messages, err := mb.List()
	if err != nil {
		return Message{}, err
	}

	found := false
	for _, msg := range messages {
		if msg.Id == msgId {
			found = true
			return msg, nil
		}
	}

	if !found {
		return Message{}, errors.New("message was not found")
	}

	return Message{}, nil
}

// Set adds a message to the store
func (mb *FSMailbox) Set(msg Message) error {
	fileName := msg.Id
	msgPath := filepath.Join(mb.location, fileName)
	msgBytes, err := json.Marshal(msg)
	if err != nil {
		return err
	}

	// write file
	if err := os.WriteFile(msgPath, msgBytes, 0644); err != nil {
		return err
	}

	return nil
}

// Delete removes one or more messages from the store
func (mb *FSMailbox) Delete(msgID ...string) error {

	for _, v := range msgID {
		if err := os.Remove(filepath.Join(mb.location, v)); err != nil {
			return err
		}
	}

	return nil
}
