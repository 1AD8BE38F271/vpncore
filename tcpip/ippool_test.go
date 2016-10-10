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
	"testing"
	"net"
)

func TestNewIP4Pool(t *testing.T) {
	subnet := &net.IPNet{
		IP:net.IP{192, 168, 0, 1},
		Mask:net.IPv4Mask(0xff, 0xff, 0xff, 0),
	}

	pool, err := NewIP4Pool(subnet)
	if err != nil {
		t.Fatal(err)
	}

	for i := 1; i < 255; i++ {
		_, err := pool.Next()
		if err != nil {
			t.Fatal(err)
		}
	}

	_, err = pool.Next()
	if err != poolFull {
		t.Fail()
	}

}
