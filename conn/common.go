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
type TransProtocol uint32

func (self *TransProtocol) UnmarshalTOML(data []byte) (err error) {
	name := string(data)
	name = strings.TrimSpace(name)
	name = strings.Trim(name, "\"")

	switch name {
	case "tcp": *self = 0
	case "kcp": *self = 1
	case "obfs4": *self = 2
	default:
		return fmt.Errorf("invalid protocal:%s", name)
	}
	return
}

const (
	PROTO_TCP = TransProtocol(0)
	PROTO_KCP = TransProtocol(1)
	PROTO_OBFS4 = TransProtocol(2)
)

