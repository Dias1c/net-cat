package messenger

import (
	"fmt"
	"net"
	"time"
)

// GetFomattedMessage -
func getFomattedMessage(serv *Server, conn net.Conn, message string, mode int) string {
	name := serv.Connections[conn]
	// Change Message
	switch mode {
	case ModeSendMessage:
		if message == "\n" {
			return ""
		}
		time := time.Now().Format(TimeDefault)
		message = fmt.Sprintf(PatternMessage, time, name, message)
	case ModeJoinChat:
		message = fmt.Sprintf(ColorYellow+PatternJoinChat+ColorReset, name)
	case ModeLeftChat:
		message = fmt.Sprintf(ColorYellow+PatternLeftChat+ColorReset, name)
	}
	return message
}
