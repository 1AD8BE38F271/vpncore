// +build darwin
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
	"net"
	"fmt"
	"errors"
	"github.com/FTwOoO/vpncore/cmd"
	"github.com/FTwOoO/vpncore/routes"

)

func maskToString(m net.IPMask) string {
	//only support IP4
	if len(m) == 0 {
		return "<nil>"
	}
	return fmt.Sprintf("%d.%d.%d.%d", m[0], m[1], m[2], m[3])
}

func (ifce *Interface) SetupNetwork(ip net.IP, subnet net.IPNet, mtu int) (err error) {

	var c string
	var peer_ip net.IP

	err = ifce.changeMTU(mtu)
	if err != nil {
		return err
	}


	if ifce.IsTUN() {
		peer_ip = generatePeerIP(ip)
		c = fmt.Sprintf("ifconfig %s inet %s %s netmask %s",
			ifce.Name(), ip.String(), peer_ip.String(), maskToString(subnet.Mask))
	} else {
		c = fmt.Sprintf("ifconfig %s inet %s netmask %s",
			ifce.Name(), ip.String(), maskToString(subnet.Mask))
	}

	_, err = cmd.RunCommand(c)
	if err != nil {
		return
	} else {
		ifce.ip = ip
		ifce.subnet = subnet
		if ifce.IsTUN() {
			ifce.peer_ip = peer_ip
		}
	}

	err = ifce.setupRoutes()
	return
}

func (ifce *Interface) ServerSetupNatRules() (err error) {
	panic("Not implemented for this platform")
}

func (ifce *Interface) changeMTU(mtu int) (err error) {

	c := fmt.Sprintf("ifconfig %s mtu %d", ifce.Name(), mtu)
	_, err = cmd.RunCommand(c)
	if err != nil {
		return err
	}

	return nil
}

func (ifce *Interface) setupRoutes() (error) {

	if ifce.IP() == nil {
		return errors.New("Setup interface IP first!")
	}
	router, err := routes.NewRoutesManager()
	if err != nil {
		return err
	}

	err = router.AddRouteToNet(ifce.Name(), ifce.subnet, ifce.IP())
	return err
}


