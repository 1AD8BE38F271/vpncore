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

var (
	zeroChecksum = [2]byte{0x00, 0x00}
)

func Sum(b []byte) uint32 {
	var sum uint32

	n := len(b)
	for i := 0; i < n; i = i + 2 {
		sum += (uint32(b[i]) << 8)
		if i + 1 < n {
			sum += uint32(b[i + 1])
		}
	}
	return sum
}

// checksum for Internet Protocol family headers
func Checksum(sum uint32, b []byte) (answer [2]byte) {
	sum += Sum(b)
	sum = (sum >> 16) + (sum & 0xffff)
	sum += (sum >> 16)
	sum = ^sum
	answer[0] = byte(sum >> 8)
	answer[1] = byte(sum)
	return
}
