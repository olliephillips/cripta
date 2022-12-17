package display

import (
	"bufio"
	"fmt"
	"github.com/olliephillips/cripta/internal/mailbox"
	"strings"
)

// ShowMessages provides a mailbox type look and feel for the messages
func ShowMessages(msgs []mailbox.Message) string {

	// understand widths
	fromH := "FROM                            "
	subH := "SUBJECT                                              "
	idH := "ID     "

	fLen := len(fromH)
	sLen := len(subH)
	iLen := len(idH)

	totalLen := fLen + sLen + iLen

	// set underlines
	ulFrom := ""
	for i := 1; i < fLen; i++ {
		ulFrom += "-"
	}
	ulFrom += " "

	ulSub := ""
	for i := 1; i < sLen; i++ {
		ulSub += "-"
	}
	ulSub += " "

	ulId := ""
	for i := 1; i < iLen; i++ {
		ulId += "-"
	}
	ulId += " "

	// header
	header := fmt.Sprintf("\n%s%s%s\n", fromH, subH, idH)
	header += fmt.Sprintf("%s%s%s\n", ulFrom, ulSub, ulId)

	body := ""
	for _, msg := range msgs {
		row := ""
		// from
		fPadLen := fLen - len(msg.Sender)
		fPad := ""
		for i := 1; i < fPadLen; i++ {
			fPad += " "
		}
		row += fmt.Sprintf("@%s%s", msg.Sender, fPad)

		// sub
		sPadLen := sLen - len(msg.Subject) + 1
		sPad := ""
		for i := 1; i < sPadLen; i++ {
			sPad += " "
		}
		row += fmt.Sprintf("%s%s", msg.Subject, sPad)

		// id
		row += fmt.Sprintf("%s\n", msg.Id)
		body += row
	}

	// footer
	footer := ""
	for i := 1; i < totalLen; i++ {
		footer += "-"
	}
	footer += "\n"

	return fmt.Sprintf("%s%s%s", header, body, footer)
}

// ShowMessage provides a email type look and feel for the message
func ShowMessage(msg mailbox.Message) string {

	fromLine := fmt.Sprintf("\nFrom:    @%s\n", msg.Sender)
	timeLine := fmt.Sprintf("Sent:    %s\n", msg.Sent)
	subjectLine := fmt.Sprintf("Subject: %s\n", msg.Subject)
	sepLine := fmt.Sprintf("---------------------------------------------------\n")
	lineLength := len(sepLine)

	messageBody := wrapText(msg.Body, lineLength)

	message := fmt.Sprintf(`%s%s%s%s%s`, fromLine, timeLine, subjectLine, sepLine, messageBody)

	return message
}

// line wrapping helper
func wrapText(text string, wrapLength int) string {
	wrappedText := ""
	scanner := bufio.NewScanner(strings.NewReader(text))
	for scanner.Scan() {
		line := scanner.Text()
		for len(line) > wrapLength {
			// find the nearest space character before the wrap length
			spaceIndex := strings.LastIndex(line[:wrapLength], " ")
			if spaceIndex == -1 {
				// if no space is found, just break at the wrap length
				wrappedText += line[:wrapLength] + "\n"
				line = line[wrapLength:]
			} else {
				// otherwise, break at the space character
				wrappedText += line[:spaceIndex] + "\n"
				line = line[spaceIndex+1:]
			}
		}
		wrappedText += line + "\n"
	}
	return wrappedText
}
