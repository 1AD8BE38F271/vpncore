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

package crypt

import (
	"net"
	"github.com/1AD8BE38F271/vpncore/enc"
)


type cryptConn struct {
	net.Conn
	B enc.BlockCrypt
}

func NewCryptConn(conn net.Conn, block enc.BlockCrypt) (*cryptConn, error) {
	connection := new(cryptConn)
	connection.Conn = conn
	connection.B = block
	return connection, nil
}

func (c *cryptConn) Read(b []byte) (n int, err error) {

	buf := make([]byte, len(b))

	n, err = c.Conn.Read(buf)
	if err != nil {
		return
	}

	c.B.Decrypt(b[:n], buf[:n])
	return
}

func (c *cryptConn) Write(b []byte) (n int, err error) {
	buf := make([]byte, len(b))
	c.B.Encrypt(buf, b)
	return c.Conn.Write(buf)
}
