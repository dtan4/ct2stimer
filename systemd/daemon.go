package systemd

import (
	"github.com/coreos/go-systemd/dbus"
)

// Client represents systemd D-Bus API client.
type Client struct {
	conn *dbus.Conn
}

// NewClient creates new Client object
func NewClient(conn *dbus.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

// NewConn establishes a new connection to D-Bus
func NewConn() (*dbus.Conn, error) {
	return dbus.New()
}

// Reload reloads systemd unit files
func (c *Client) Reload() error {
	if err := c.conn.Reload(); err != nil {
		return err
	}

	return nil
}
