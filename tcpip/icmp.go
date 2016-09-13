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

package tcpip

import (
	"encoding/binary"
)

type ICMPType byte

const (
	ICMPEcho ICMPType = 0x0
	ICMPRequest = 0x8
)

type ICMPPacket []byte

func (p ICMPPacket) Type() ICMPType {
	return ICMPType(p[0])
}

func (p ICMPPacket) SetType(v ICMPType) {
	p[0] = byte(v)
}

func (p ICMPPacket) Code() byte {
	return p[1]
}

func (p ICMPPacket) Checksum() uint16 {
	return binary.BigEndian.Uint16(p[2:])
}

func (p ICMPPacket) SetChecksum(sum [2]byte) {
	p[2] = sum[0]
	p[3] = sum[1]
}

func (p ICMPPacket) ResetChecksum() {
	p.SetChecksum(zeroChecksum)
	p.SetChecksum(Checksum(0, p))
}
