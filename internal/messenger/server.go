package messenger

import (
	"bufio"
	"errors"
	"fmt"
	"log"
	"net"
	"strings"
	"time"
)

// ====================
// ====== Public ======
// ====================

// Constructor - for create Server with input data
func (s *Server) Constructor(port string, maxConn int) error {
	serv, err := net.Listen("tcp", port)
	if err != nil {
		return err
	}
	// Set settings
	s.Server = serv
	s.MaxConnections = maxConn
	s.Connections = make(map[net.Conn]string, maxConn)
	s.UsedNames = make(map[string]bool, maxConn)
	return nil
}

// CanConnect - check connection for connect
func (s *Server) CanConnect(conn net.Conn) bool {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return !(s.MaxConnections != 0 && len(s.Connections) >= s.MaxConnections)
}

// ConnectMessenger - Join conn to group after writing name and starting comminicate with him.
// Closes connection on ERROR or finish Chatting
func (s *Server) ConnectMessenger(conn net.Conn) {
	if !s.CanConnect(conn) {
		fmt.Fprint(conn, "The room is full, please try again later...")
		conn.Close()
		return
	}
	fmt.Fprint(conn, WelcomMessage)
	name, err := bufio.NewReader(conn).ReadString('\n')
	if err != nil {
		log.Printf("ConnectMessenger: %v\n", err.Error())
		return
	}
	name = strings.Replace(name, "\n", "", 1)
	err = s.addConnection(conn, name)
	if err != nil {
		fmt.Fprint(conn, err.Error())
		conn.Close()
		return
	}
	s.startChatting(conn)
	s.removeConnection(conn)
}

// CloseServer - closing all connections and close connections
func (s *Server) CloseServer() {
	log.Println("Closing Server")
	s.mutex.Lock()
	for conn := range s.Connections {
		fmt.Fprintf(conn, "\n%sServer Was Closed!%s", BgColorRed, ColorReset)
		conn.Close()
	}
	s.mutex.Unlock()
	s.Server.Close()
	log.Println("Server Closed")
}

// ====================
// ===== PRIVATE ======
// ====================

// StartChatting - Load old messages and this connection starts communicate with any connections
func (s *Server) startChatting(conn net.Conn) {
	s.loadMessages(conn)
	message := getFomattedMessage(s, conn, "", ModeJoinChat)
	s.sendMessage(conn, message)
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			break
		}
		message = getFomattedMessage(s, conn, message, ModeSendMessage)
		s.sendMessage(conn, message)
		s.saveMessage(message)
	}
	message = getFomattedMessage(s, conn, "", ModeLeftChat)
	s.sendMessage(conn, message)
}

// SendMessage - Sends message for all users by mode
func (s *Server) sendMessage(conn net.Conn, message string) {
	if message == "" {
		fmt.Fprintf(conn, PatternSending, time.Now().Format(TimeDefault), s.Connections[conn])
		return
	}
	// SendingMessage
	time := time.Now().Format(TimeDefault)
	sendMessage := fmt.Sprintf("%s!%s\n%s", ColorYellow, ColorReset, message)
	s.mutex.Lock()
	for con := range s.Connections {
		if con != conn {
			fmt.Fprint(con, sendMessage)
		}
		fmt.Fprintf(con, PatternSending, time, s.Connections[con])
	}
	s.mutex.Unlock()
}

// LoadMessages - Sends History of message for conn
func (s *Server) loadMessages(conn net.Conn) {
	for _, message := range s.AllMessages {
		fmt.Fprintf(conn, message)
	}
}

// SaveMessage - Saving message to AllMessages (safe)
func (s *Server) saveMessage(message string) {
	s.mutex.Lock()
	s.AllMessages = append(s.AllMessages, message)
	s.mutex.Unlock()
}

// AddConnection - Adding conncetion to s.Connections (safe)
func (s *Server) addConnection(conn net.Conn, name string) error {
	if name == "" {
		return errors.New("Name cant be empty")
	}
	s.mutex.Lock()
	defer s.mutex.Unlock()
	if s.MaxConnections != 0 && len(s.Connections) >= s.MaxConnections {
		return fmt.Errorf("The room is full [%v]", conn.RemoteAddr())
	} else if s.UsedNames[name] {
		return fmt.Errorf("Name '%s' is Exist [%v]", name, conn.RemoteAddr())
	}
	s.UsedNames[name] = true
	s.Connections[conn] = name
	log.Printf("Connected %v", conn.RemoteAddr())
	return nil
}

// RemoveConnection - Removing Connection from s.Connections (safe)
func (s *Server) removeConnection(conn net.Conn) {
	s.mutex.Lock()
	delete(s.UsedNames, s.Connections[conn])
	delete(s.Connections, conn)
	log.Printf("Connect %v was left", conn.RemoteAddr())
	s.mutex.Unlock()
}
