package input

import (
	"bufio"
	"fmt"
	"github.com/olliephillips/cripta/internal/display"
	"github.com/olliephillips/cripta/internal/group"
	"github.com/olliephillips/cripta/internal/mailbox"
	"log"
	"os"
	"strings"
)

const (
	DELETE_ALL = iota + 1
	DELETE_SINGLE
)

// ReadStdin monitors Stdin for input and responds accordingly
func ReadStdin(outbox chan<- string, disconnect chan<- struct{}, mb mailbox.MailStorer) {
	reader := bufio.NewReader(os.Stdin)

	// we use this as a flag to confirm deletion
	deleteRequest := false
	deleteOperation := 0
	deleteId := ""

	// loop to capture input
	for {
		fmt.Print("> ")
		text, _ := reader.ReadString('\n')
		text = strings.Replace(text, "\n", "", -1)

		// multi-word commands
		split := strings.SplitN(text, " ", 2)
		splitLen := len(split)

		switch {
		case text == "help":
			fmt.Println("Showing available commands..")
			fmt.Println(printCommands())

		case text == "quit":
			disconnect <- struct{}{}
			fmt.Println("Exiting")
			return
		case text == "list":
			messages, err := mb.List()
			if err != nil {
				log.Printf("There was a problem, could not list messages: %v\n", err)
			}

			// process messages
			output := display.ShowMessages(messages)
			fmt.Println(output)
		case strings.Contains(text, "read") && splitLen > 1:
			message, err := mb.Get(split[1])
			if err != nil {
				log.Printf("There was a problem, could not get message: %v\n", err)
				panic(err.(any))
			}
			// process message
			output := display.ShowMessage(message)
			fmt.Println(output)
		case strings.HasPrefix(text, "@") && splitLen > 1:
			fmt.Println("Sending message..")
			outbox <- text
		case strings.Contains(text, "delete") && splitLen > 1:
			fmt.Println("Are you sure you wish to delete this message..?")
			fmt.Println("Type 'confirm' to proceed..")
			deleteRequest = true
			deleteOperation = DELETE_SINGLE
			deleteId = split[1]
		case text == "empty":
			fmt.Println("Are you sure you wish to delete all messages..?")
			fmt.Println("Type 'confirm' to proceed..")
			deleteRequest = true
			deleteOperation = DELETE_ALL
		case text == "confirm":
			if deleteRequest == false && deleteOperation == 0 {
				fmt.Println("Use 'delete' or 'empty' command first..")
				continue
			}
			switch deleteOperation {
			case DELETE_ALL:
				messages, err := mb.List()
				if err != nil {
					log.Println("There was a problem, could not get messages for deletion")
					panic(err.(any))
				}

				deleteList := []string{}
				for _, msg := range messages {
					deleteList = append(deleteList, msg.Id)
				}

				err = mb.Delete(deleteList...)
				if err != nil {
					log.Println("There was a problem, could not delete messages")
					panic(err.(any))
				}

				fmt.Println("Deleted all messages")

				// reset
				deleteRequest = false
				deleteOperation = 0
			case DELETE_SINGLE:
				err := mb.Delete(deleteId)
				if err != nil {
					log.Println("There was a problem, could not delete message")
					panic(err.(any))
				}

				fmt.Printf("Deleted message (id: %v)\n", deleteId)

				// reset
				deleteRequest = false
				deleteOperation = 0
				deleteId = ""
			}
		case strings.HasPrefix(text, "->") && splitLen > 1:
			fmt.Println("group send - not implemented..")
		case text == "friends":
			var out string
			fr, err := mailbox.ListFriends()
			if err != nil {
				log.Println("There was a problem, could not list friends")
				panic(err.(any))
			}
			for _, v := range fr {
				out += fmt.Sprintf("%s\n", v)
			}
			fmt.Println(out)
		case text == "groups":
			var out string
			fr, err := group.ListGroups()
			if err != nil {
				log.Println("There was a problem, could not list groups")
				panic(err.(any))
			}
			for _, v := range fr {
				out += fmt.Sprintf("%s\n", v)
			}
			fmt.Println(out)
		case text == "":
		default:
			fmt.Println("Not understood, showing available commands..")
			fmt.Println(printCommands())
		}
	}
}

func printCommands() string {
	cmd := `
Command List
------------
help			Show available commands
list			List all messages in mailbox
friends			List available friends (with public keys)
groups			List available groups
read <msgId>		Print the message with id <msgId> to the console
delete <msgId>		Delete the messsage with id <msgId> from mailbox
empty			Delete all messages in the mailbox
quit			Quit the program

To send a message to a single recipient:-
> @<username> <subject>::<message>

To send to a group of recipients (group must exist):-
> -><groupname> <subject>::<message>

The message(s) will be queue and published.
`
	return cmd
}
