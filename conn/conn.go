/*
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU General Public License as published by
 * the Free Software Foundation, either version 3 of the License, or
 * (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU General Public License for more details.
 *
 * You should have received a copy of the GNU General Public License
 * along with this program.  If not, see <http://www.gnu.org/licenses/>.
 *
 * Author: FTwOoO <booobooob@gmail.com>
 */

package conn

import (
	"net"
	"time"
	"github.com/FTwOoO/go-enc"
)

type Connection struct {
	C net.Conn
	B enc.BlockCrypt
}

func NewConnection(conn net.Conn, config *enc.BlockConfig) (*Connection, error) {
	connection := new(Connection)
	connection.C = conn
	block, err := enc.NewBlock(config)
	if err != nil {
		return nil, err
	} else {
		connection.B = block
		return connection, nil
	}

}

func (c *Connection) Read(b []byte) (n int, err error) {
	n, err = c.C.Read(b)
	if err != nil {
		return
	}

	c.B.Decrypt(b[:n], b[:n])
	return

}

func (c *Connection) Write(b []byte) (n int, err error) {
	n, err = c.C.Write(b)
	if err != nil {
		return
	}

	c.B.Encrypt(b[:n], b[:n])
	return
}

func (c *Connection) Close() error {
	return c.C.Close()
}

func (c *Connection) LocalAddr() net.Addr {
	return c.C.LocalAddr()
}

func (c *Connection) RemoteAddr() net.Addr {
	return c.C.RemoteAddr()
}

func (c *Connection)SetDeadline(t time.Time) error {
	return c.C.SetDeadline(t)
}

func (c *Connection)SetReadDeadline(t time.Time) error {
	return c.C.SetReadDeadline(t)
}

func (c *Connection)SetWriteDeadline(t time.Time) error {
	return c.C.SetWriteDeadline(t)
}