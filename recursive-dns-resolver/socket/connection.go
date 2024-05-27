package socket

import (
	"fmt"
	"net"
	"time"
)

func SetupUDPConnection(serverAddr string, timeout time.Duration) (net.Conn, error) {
	socket := ":53"
	conn, err := net.Dial("udp", serverAddr+socket)
	if err != nil {
		return nil, fmt.Errorf("failed to create UDP connection: %w", err)
	}

	return conn, nil
}
func OpenUDPConnection(rootserver string) (net.Conn, error) {
	connection, err := SetupUDPConnection(rootserver, 5*time.Second)
	if err != nil {
		return nil, fmt.Errorf("error setting up udp connection: %w", err)
	}
	return connection, nil
}

func CloseUDPConnection(connection net.Conn) {
	connection.Close()
}
