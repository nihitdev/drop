package network

import (
	"fmt"
	"net"
	"net/http"
)

func LocalIP() (string, error) {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		return "", err
	}
	defer conn.Close()

	addr, ok := conn.LocalAddr().(*net.UDPAddr)
	if !ok {
		return "", fmt.Errorf("unexpected network address %T", conn.LocalAddr())
	}
	return addr.IP.String(), nil
}

func Listen(port int) (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf(":%d", port))
}

func Port(listener net.Listener) int {
	if addr, ok := listener.Addr().(*net.TCPAddr); ok {
		return addr.Port
	}
	return 0
}

func URL(ip string, listener net.Listener, token string) string {
	return fmt.Sprintf("http://%s:%d/d/%s", ip, Port(listener), token)
}

func ClientIP(r *http.Request) string {
	host, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return host
}
