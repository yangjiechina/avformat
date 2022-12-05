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
}

type transport struct {
	onPacketHandler       OnPacketHandler
	onDisConnectedHandler OnDisconnectedHandler
}

func (t *transport) SetOnPacketHandler(handler OnPacketHandler) {
	t.onPacketHandler = handler
}

func (t *transport) SetOnDisconnectedHandler(handler OnDisconnectedHandler) {
	t.onDisConnectedHandler = handler
}

type TCPClient struct {
	transport
	conn   net.Conn
	cancel context.CancelFunc
}

func NewTCPClient(localAddr *net.TCPAddr, serverIp string, serverPort int) (*TCPClient, error) {
	dialer := net.Dialer{LocalAddr: localAddr}
	if dial, err := dialer.Dial("tcp", fmt.Sprintf("%s:%d", serverIp, serverPort)); err != nil {
		return nil, err
	} else {
		return &TCPClient{conn: dial}, nil
	}
}

func (c *TCPClient) Write(data []byte) (int, error) {
	return c.conn.Write(data)
}

func (c *TCPClient) doRead() {
	var err error
	var n int
	var ctx context.Context
	ctx, c.cancel = context.WithCancel(context.Background())

	bytes := make([]byte, 16000)
	for ctx.Err() == nil {
		n, err = c.conn.Read(bytes)
		if err != nil {
			break
		}

		if c.onPacketHandler != nil {
			c.onPacketHandler(c.conn, bytes[:n])
		}
	}

	if c.onDisConnectedHandler != nil {
		c.onDisConnectedHandler(c.conn, err)
	}
}

func (c *TCPClient) Read() {
	go c.doRead()
}

func (c *TCPClient) Close() error {
	c.cancel()
	return c.conn.Close()
}
