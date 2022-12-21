package utils

import (
	"context"
	"fmt"
	"net"
)

type OnPacketHandler func(conn net.Conn, data []byte)
type OnDisconnectedHandler func(conn net.Conn, err error)

type Transport interface {
	SetOnPacketHandler(OnPacketHandler)
	SetOnDisconnectedHandler(OnDisconnectedHandler)
	Conn() net.Conn
	Write([]byte) (int, error)
	Close() error
	Read()
	ListenPort() int
}

type transport struct {
	onPacketHandler       OnPacketHandler
	onDisConnectedHandler OnDisconnectedHandler
	conn                  net.Conn
	cancel                context.CancelFunc
	listPort              int
}

func (t *transport) SetOnPacketHandler(handler OnPacketHandler) {
	t.onPacketHandler = handler
}

func (t *transport) SetOnDisconnectedHandler(handler OnDisconnectedHandler) {
	t.onDisConnectedHandler = handler
}

func (t *transport) Conn() net.Conn {
	return t.conn
}

func (t *transport) Write(data []byte) (int, error) {
	return t.conn.Write(data)
}

func (t *transport) Close() error {
	t.cancel()
	return t.conn.Close()
}

func (t *transport) doRead() {
	var err error
	var n int
	var ctx context.Context
	ctx, t.cancel = context.WithCancel(context.Background())

	bytes := make([]byte, 16000)
	for ctx.Err() == nil {
		n, err = t.conn.Read(bytes)
		if err != nil {
			break
		}

		if t.onPacketHandler != nil {
			t.onPacketHandler(t.conn, bytes[:n])
		}
	}

	if t.onDisConnectedHandler != nil {
		t.onDisConnectedHandler(t.conn, err)
	}
}

func (t *transport) Read() {
	go t.doRead()
}

func (t *transport) ListenPort() int {
	return t.listPort
}

type TCPClient struct {
	transport
}

type UDPTransport struct {
	transport
}

func (u *UDPTransport) WriteTo(data []byte, ip string, port int) (int, error) {
	return u.conn.(*net.UDPConn).WriteTo(data, &net.UDPAddr{IP: net.ParseIP(ip), Port: port})
}

func NewTCPClient(localAddr *net.TCPAddr, serverIp string, serverPort int) (Transport, error) {
	dialer := net.Dialer{LocalAddr: localAddr}
	if dial, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort)); err != nil {
		return nil, err
	} else {
		return &TCPClient{transport: transport{conn: dial, listPort: dial.LocalAddr().(*net.TCPAddr).Port}}, nil
	}
}

func NewUDPTransport(port int) (Transport, error) {
	udp, err := net.ListenUDP("udp", &net.UDPAddr{IP: net.ParseIP("0.0.0.0"), Port: port})
	if err != nil {
		return nil, err
	}
	return &UDPTransport{transport{conn: udp, listPort: port}}, nil
}
