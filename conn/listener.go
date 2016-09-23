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

type Listener struct {
	listener net.Listener
	proto    TransProtocol
	blockConfig *enc.BlockConfg

}

func NewListener(proto TransProtocol, listenAddr string, blockConfig enc.BlockConfg) (l net.Listener, err error) {
	switch proto {
	case PROTO_KCP:
		return &Listener{proto:proto, listener:&KCPListener{}, blockConfig:blockConfig}, nil
	case PROTO_TCP:
		addr, err := net.ResolveTCPAddr("tcp4", listenAddr)
		if err != nil {
			return nil, err
		}
		l, err := net.ListenTCP("tcp4", addr)
		return l, nil
	default:
		return nil, errors.New("UNKOWN PROTOCOL!")
	}
}

func (l *Listener) Accept() (net.Conn, error) {
	for {
		conn, err := l.listener.Accept()
		if err != nil {
			return nil, err
		} else {
			return NewConnection(conn, l.blockConfig)
		}
	}
	return
}

// Close closes the listener.
// Any blocked Accept operations will be unblocked and return errors.
func (l *Listener) Close() error {
	return l.listener.Close()
}

// Addr returns the listener's network address.
func (l *Listener) Addr() net.Addr {
	return l.listener.Addr()
}

