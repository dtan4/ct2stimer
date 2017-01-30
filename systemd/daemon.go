package systemd

import (
	"github.com/coreos/go-systemd/dbus"
	"github.com/pkg/errors"
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

// StartUnit starts the given systemd unit file
func (c *Client) StartUnit(name string) error {
	ch := make(chan string)

	if _, err := c.conn.StartUnit(name, "replace", ch); err != nil {
		return errors.Wrapf(err, "failed to start systemd unit %q", name)
	}

	if job := <-ch; job != "done" {
		return errors.Errorf("cannot start service %q, current status is %q", name, job)
	}

	return nil
}

// Reload reloads systemd unit files
func (c *Client) Reload() error {
	if err := c.conn.Reload(); err != nil {
		return errors.Wrap(err, "failed to reload systemd unit files")
	}

	return nil
}
