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
package vpncore

import (
	"strings"
	"net"
)

type DNSList []net.IP

func (l *DNSList) String() string {
	r := make([]string, len(*l))

	for i, value := range *l {
		r[i] = value.String()
	}

	return strings.Join(r, " ")
}

func (ls *DNSList)UnmarshalText(text []byte) error {
	s := string(text)
	var l []string

	if strings.Contains(s, ",") {
		l = strings.Split(s, ",")
	} else {
		l = strings.Split(s, " ")
	}

	for _, v := range l {
		v = strings.TrimSpace(v)
		*ls = append(*ls, net.ParseIP(v))
	}
	return nil

}

type DNSManager struct {
	old_dns DNSList
}