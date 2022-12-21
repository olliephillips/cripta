package mailbox

import (
	"github.com/olliephillips/cripta/internal"
	"github.com/olliephillips/cripta/internal/crypt"
	"log"
	"strings"
	"time"
)

// Message is the unencrypted representation of message data
type Message struct {
	Id      string `json:"id"`
	Subject string `json:"subject"`
	Body    string `json:"body"`
	Sender  string `json:"sender"`
	Sent    string `json:"time"`
}

// OutboundMessage ideally stores the recipient
type OutboundMessage struct {
	Message
	To        string
	Signature []byte `json:"signature"`
}

type EncryptedMQTT struct {
	Topic   string
	Message []byte
}

// NewOutboundMessage parses the raw message and creates an OutboundMessage
func NewOutboundMessage(raw string, fromUser string) OutboundMessage {
	var sub, body string
	split := strings.SplitN(raw, " ", 2)
	user := strings.TrimLeft(split[0], "@")

	splitSubBody := strings.Split(split[1], "::")
	if len(splitSubBody) == 1 {
		sub = "No subject"
		body = splitSubBody[0]
	} else {
		sub = splitSubBody[0]
		body = splitSubBody[1]
	}

	obm := OutboundMessage{}
	obm.To = user
	obm.Id = internal.ShortUID()
	obm.Subject = sub
	obm.Body = body
	obm.Sender = fromUser
	obm.Sent = time.Now().Format("Mon Jan 2 15:04 MST 2006")

	// sign the message body and include sig
	sig, err := crypt.SignPayload([]byte(obm.Body), internal.PRIVATE_KEY_FILE)
	if err != nil {
		log.Printf("Failed to get signature: %v\n", err)
	}
	obm.Signature = sig

	return obm
}
