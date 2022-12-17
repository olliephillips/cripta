package main

import (
	"encoding/json"
	"errors"
	"fmt"
	MQTT "github.com/eclipse/paho.mqtt.golang"
	"github.com/joho/godotenv"
	"github.com/olliephillips/cripta/internal"
	"github.com/olliephillips/cripta/internal/crypt"
	"github.com/olliephillips/cripta/internal/input"
	"github.com/olliephillips/cripta/internal/mailbox"
	"log"
	"os"
	"path/filepath"
	"strings"
)

var (
	version   = "Development build"
	goversion = "Unknown"
)

func init() {
	fmt.Println("\nCripta Messenger")

	fmt.Println("---------------------")
	fmt.Printf("Version: %s\nGo version: %s\n", version, goversion)
	fmt.Println("---------------------")
	fmt.Println("\nEnter 'help' for command list.\n")

	// read config
	err := godotenv.Load(internal.CONFIG_FILE)
	if err != nil {
		log.Fatalf("Failed to read database config file: %v\n", err)
	}

	// create mailbox folder (we only do fs at moment)
	mbFolder := filepath.Join(".", internal.MAILBOX_FOLDER)
	if err := os.MkdirAll(mbFolder, os.ModePerm); err != nil {
		log.Fatalf("Problem creating mailbox folder: %v\n", err)
	}

	// create groups folder
	groupFolder := filepath.Join(".", internal.GROUPS_FOLDER)
	if err := os.MkdirAll(groupFolder, os.ModePerm); err != nil {
		log.Fatalf("Problem creating groups folder: %v\n", err)
	}

	// create friend public keys folder
	keysFolder := filepath.Join(".", internal.FRIEND_KEYS_FOLDER)
	if err := os.MkdirAll(keysFolder, os.ModePerm); err != nil {
		log.Fatalf("Problem creating friends keys folder: %v\n", err)
	}

	// check if user has pub/private key
	pubKeyPath := filepath.Join(".", internal.PUBLIC_KEY_FILE)
	privateKeyPath := filepath.Join(".", internal.PRIVATE_KEY_FILE)
	makeNewKeys := false
	if _, err := os.Stat(pubKeyPath); errors.Is(err, os.ErrNotExist) {
		makeNewKeys = true
	}
	if _, err := os.Stat(privateKeyPath); errors.Is(err, os.ErrNotExist) {
		makeNewKeys = true
	}

	// if not make new
	if makeNewKeys {
		log.Println("Creating your public/private key pair")
		if err := crypt.MakeKeys(pubKeyPath, privateKeyPath); err != nil {
			log.Fatalf("Problem creating the public/private key pair: %v\n", err)
		}
	}
}

func main() {
	// store
	mb := mailbox.NewFSMailbox(internal.MAILBOX_FOLDER)

	// we use channels for the outbox messages and send queue
	outbox := make(chan string, 1)
	sendQueue := make(chan mailbox.EncryptedMQTT, 1)
	disconnect := make(chan struct{})
	quit := make(chan struct{})

	// mqtt setup
	go func() {
		// subscription handler
		var f MQTT.MessageHandler = func(client MQTT.Client, mqttMsg MQTT.Message) {
			// decrypt payload
			jsonPayload, err := crypt.DecryptMessageWithPrivateKey(mqttMsg.Payload(), internal.PRIVATE_KEY_FILE)
			if err != nil {
				log.Printf("Problem decrypting message: %v\n", err)
				return
			}

			// make message
			var msg mailbox.Message
			if err := json.Unmarshal(jsonPayload, &msg); err != nil {
				log.Printf("Problem unmarshalling message: %v\n", err)
				return
			}

			if err := mb.Set(msg); err != nil {
				log.Printf("Problem writing message: %v\n", err)
			}
		}

		mqttServer := fmt.Sprintf("tcp://%s:%s", internal.SERVER, internal.PORT)
		opts := MQTT.NewClientOptions().AddBroker(mqttServer)
		opts.SetClientID(os.Getenv("TWITTER_USERNAME"))
		opts.SetDefaultPublishHandler(f)

		// create and start a client using the above ClientOptions
		c := MQTT.NewClient(opts)
		if token := c.Connect(); token.Wait() && token.Error() != nil {
			panic(token.Error().(any))
		}

		// subscribe to our own topic to collect messages to us
		subscribeTopic := fmt.Sprintf("uk/misc/testing/%s", strings.ToLower(os.Getenv("TWITTER_USERNAME")))
		if token := c.Subscribe(subscribeTopic, 2, f); token.Wait() && token.Error() != nil {
			fmt.Println(token.Error())
			os.Exit(1)
		}

		for {
			select {
			case msg := <-sendQueue:
				token := c.Publish(msg.Topic, 2, false, msg.Message)
				token.Wait()
			case <-disconnect:
				c.Unsubscribe(subscribeTopic)
				c.Disconnect(0)
				// we've gracefully shutdown the MQTT subscribe and connection
				// we can quit here
				quit <- struct{}{}
			}
		}
	}()

	// start reading Stdin
	go func() {
		input.ReadStdin(outbox, disconnect, mb)
	}()

	// message routing
	for {
		select {
		case msg := <-outbox:
			// process for sending
			obm := mailbox.NewOutboundMessage(msg, os.Getenv("TWITTER_USERNAME"))

			// marshal for json
			json, err := json.Marshal(obm)

			// encrypt
			encrypted, err := crypt.EncryptMessageWithPublicKey(json, obm.To, internal.FRIEND_KEYS_FOLDER)
			if err != nil {
				log.Printf("Could not encrypt message: %v", err)
				continue
			}

			topic := fmt.Sprintf("uk/misc/testing/%s", strings.ToLower(obm.Sender))
			payload := encrypted

			send := mailbox.EncryptedMQTT{topic, payload}
			sendQueue <- send
		case <-quit:
			return
		}
	}
}
