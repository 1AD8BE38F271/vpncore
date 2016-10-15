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

package stream

import (
	"net"
	"errors"
	"github.com/FTwOoO/vpncore/conn"
)

type StreamLayerContext struct {
	Protocol   conn.TransProtocol
	ListenAddr string
	RemoveAddr string
}

func (this *StreamLayerContext) Dial(_ net.Conn) (net.Conn, error) {

	switch this.Protocol {
	case conn.PROTO_TCP:
		c, err := net.Dial("tcp", this.RemoveAddr)
		if err != nil {
			return nil, err
		}

		return &streamConn{Conn:c, proto:this.Protocol}, nil


	case conn.PROTO_KCP:
		panic("not implemented!")
	case conn.PROTO_OBFS4:
		panic("not implemented!")
	}

	return nil, errors.New("Proto not supported!")
}

func (this *StreamLayerContext) NewListener(_ net.Listener) (net.Listener, error) {
	switch this.Protocol {
	case conn.PROTO_KCP:
		panic("not implemented yet!")
	case conn.PROTO_TCP:
		addr, err := net.ResolveTCPAddr("tcp4", this.ListenAddr)
		if err != nil {
			return nil, err
		}

		l, err := net.ListenTCP("tcp4", addr)
		if err != nil {
			return nil, err
		}
		return &streamListener{proto:this.Protocol, Listener:l}, nil
	default:
		return nil, errors.New("UNKOWN PROTOCOL!")
	}
}

func (this *StreamLayerContext) Layer() conn.ConnLayer {
	return conn.STREAM_LAYER
}

