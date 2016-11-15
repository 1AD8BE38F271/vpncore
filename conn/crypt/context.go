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
	"github.com/1AD8BE38F271/vpncore/conn"
)

type CryptLayerContext struct {
	*enc.BlockConfig
}

func (this *CryptLayerContext) Dial(c net.Conn) (net.Conn, error) {

	block, err := enc.NewBlock(this.BlockConfig)
	if err != nil {
		return nil, err
	}

	return NewCryptConn(c, block)
}

func (this *CryptLayerContext) NewListener(l net.Listener) (net.Listener, error) {
	block, err := enc.NewBlock(this.BlockConfig)
	if err != nil {
		return nil, err
	}

	return &cryptListener{Listener:l, block:block}, nil
}


func (this *CryptLayerContext) 	Layer() conn.ConnLayer {
	return conn.CRYPT_LAYER
}
