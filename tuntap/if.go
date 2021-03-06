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
package tuntap

import (
	"io"
	"net"
)

const (
	DEFAULT_HWADDR_PREFIX = "01:02:03:04:05:"
	DEFUALT_HWADDR_BRD = "ff:ff:ff:ff:ff:ff"
)

// Interface is a TUN/TAP interface.
type Interface struct {
	ip            net.IP
	subnet        net.IPNet
	peer_ip       net.IP
	isTAP         bool
	io.ReadWriteCloser
	name          string
}

// Create a new TAP interface whose name is ifName.
// If ifName is empty, a default name (tap0, tap1, ... ) will be assigned.
// ifName should not exceed 16 bytes.
func NewTAP(ifName string) (ifce *Interface, err error) {
	return newTAP(ifName)
}

// Create a new TUN interface whose name is ifName.
// If ifName is empty, a default name (tap0, tap1, ... ) will be assigned.
// ifName should not exceed 16 bytes.
func NewTUN(ifName string) (ifce *Interface, err error) {
	return newTUN(ifName)
}

// Returns true if ifce is a TUN interface, otherwise returns false;
func (ifce *Interface) IsTUN() bool {
	return !ifce.isTAP
}

// Returns true if ifce is a TAP interface, otherwise returns false;
func (ifce *Interface) IsTAP() bool {
	return ifce.isTAP
}

func (ifce *Interface) PeerIP() net.IP {
	if ifce.IsTUN() {
		return ifce.peer_ip
	} else {
		return nil
	}
}

func (ifce *Interface) IP() net.IP {
	return ifce.ip
}

func (ifce *Interface) Net() net.IPNet {
	return ifce.subnet
}

// Returns the interface name of ifce, e.g. tun0, tap1, etc..
func (ifce *Interface) Name() string {
	return ifce.name
}

