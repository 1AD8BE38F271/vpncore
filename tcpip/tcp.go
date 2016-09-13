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

type TCPPacket []byte

func (p TCPPacket) SourcePort() uint16 {
	return binary.BigEndian.Uint16(p)
}

func (p TCPPacket) SetSourcePort(port uint16) {
	binary.BigEndian.PutUint16(p, port)
}

func (p TCPPacket) DestinationPort() uint16 {
	return binary.BigEndian.Uint16(p[2:])
}

func (p TCPPacket) SetDestinationPort(port uint16) {
	binary.BigEndian.PutUint16(p[2:], port)
}

func (p TCPPacket) SetChecksum(sum [2]byte) {
	p[16] = sum[0]
	p[17] = sum[1]
}

func (p TCPPacket) Checksum() uint16 {
	return binary.BigEndian.Uint16(p[16:])
}

func (p TCPPacket) ResetChecksum(psum uint32) {
	// psum is calc by ip4packet.PseudoSum()
	p.SetChecksum(zeroChecksum)
	p.SetChecksum(Checksum(psum, p))
}
