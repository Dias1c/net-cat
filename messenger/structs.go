package messenger

import (
	"net"
	"sync"
)

// Server -
type Server struct {
	Server         net.Listener        // Server connection
	Connections    map[net.Conn]string // map[connection]Name
	MaxConnections int                 // 0 = no limit
	AllMessages    []string            // History of messages
	mutex          sync.Mutex          // Mutex for sync messages
}
