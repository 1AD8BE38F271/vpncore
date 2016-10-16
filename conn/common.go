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
	"errors"
	"net"
	"strings"
	"fmt"
)


var (
	ErrInvalidArgs     = errors.New("Invalid arguments.")
	ErrInvalidCtx      = errors.New("Invalid context")


)

type ConnLayer int
const (
	STREAM_LAYER = ConnLayer(1)
	OBS_LAYER = ConnLayer(2)
	CRYPT_LAYER = ConnLayer(3)
	AUTH_LAYER = ConnLayer(4)
	APPCATIOIN_LAYER = ConnLayer(5)
)

type ConnDialer interface {
	Dial(net.Conn) (net.Conn, error)
}

type ConnListener interface {
	NewListener(net.Listener) (net.Listener, error)
}

type ConnLayerContext interface {
	Layer() ConnLayer
	ConnDialer
	ConnListener
}

type TransProtocol string
const (
	PROTO_TCP = TransProtocol("tcp")
	PROTO_KCP = TransProtocol("kcp")
	PROTO_OBFS4 = TransProtocol("obfs4")
)

func (self *TransProtocol) UnmarshalTOML(data []byte) (err error) {
	name := string(data)
	name = strings.TrimSpace(name)
	name = strings.Trim(name, "\"")

	switch name {
	case string(PROTO_TCP): *self = PROTO_TCP
	case string(PROTO_KCP): *self = PROTO_KCP
	case string(PROTO_OBFS4): *self = PROTO_OBFS4
	default:
		return fmt.Errorf("invalid protocal:%s", name)
	}
	return
}