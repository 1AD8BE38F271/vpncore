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

type UDPPacket []byte

func (p UDPPacket) SourcePort() uint16 {
	return binary.BigEndian.Uint16(p)
}

func (p UDPPacket) SetSourcePort(port uint16) {
	binary.BigEndian.PutUint16(p, port)
}

func (p UDPPacket) DestinationPort() uint16 {
	return binary.BigEndian.Uint16(p[2:])
}

func (p UDPPacket) SetDestinationPort(port uint16) {
	binary.BigEndian.PutUint16(p[2:], port)
}

func (p UDPPacket) SetChecksum(sum [2]byte) {
	p[6] = sum[0]
	p[7] = sum[1]
}

func (p UDPPacket) Checksum() uint16 {
	return binary.BigEndian.Uint16(p[6:])
}

func (p UDPPacket) ResetChecksum(psum uint32) {
	// psum is calc by ip4packet.PseudoSum()
	p.SetChecksum(zeroChecksum)
	p.SetChecksum(Checksum(psum, p))
}
