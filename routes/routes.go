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
	"net"
	"errors"
	"sync"
)

type RoutesManager struct {
	default_gateway    net.IP
	default_nic        string

	is_gateway_changed bool
	host_routes        map[string]net.IP
	net_routes         map[string]net.IP
	mu                 sync.Mutex
}

func NewRoutesManager() (m *RoutesManager, err error) {
	m = new(RoutesManager)

	old_gateway, old_nic, err := m.GetCurrentNetGateway()
	if err != nil {
		return
	}

	m.host_routes = map[string]net.IP{}
	m.net_routes = map[string]net.IP{}

	m.default_gateway = old_gateway
	m.default_nic = old_nic
	return

}

func (self *RoutesManager)AddRouteToNet(iface string, dest net.IPNet, nextHop net.IP) (err error) {
	self.mu.Lock()
	defer self.mu.Unlock()

	err = addRouteToNet(iface, dest, nextHop)
	if err != nil {
		return
	}

	self.net_routes[dest.String()] = nextHop
	return

}

func (self *RoutesManager)AddRouteToHost(iface string, dest net.IP, nextHop net.IP) (err error) {
	self.mu.Lock()
	defer self.mu.Unlock()

	err = addRouteToHost(iface, dest, nextHop)
	if err != nil {
		return
	}

	self.host_routes[dest.String()] = nextHop
	return
}

func (self *RoutesManager) SetNewGateway(iface string, gw net.IP) (err error) {
	self.mu.Lock()
	defer self.mu.Unlock()

	if gw == nil {
		return errors.New("Argument gw is nil!")
	}

	err = redirectGateway(iface, gw)
	if err != nil {
		return err
	}

	self.is_gateway_changed = true

	return nil
}

func (self *RoutesManager) DeleteAllRoutes() (err error) {
	self.mu.Lock()
	defer self.mu.Unlock()

	for k, _ := range self.host_routes {
		delHostRoute(net.ParseIP(k))
	}

	for k, _ := range self.net_routes {
		_, subnet, _ := net.ParseCIDR(k)
		delNetRoute(*subnet)
	}

	return
}

func (self *RoutesManager) RestoreGateWay() (err error) {
	self.mu.Lock()
	defer self.mu.Unlock()

	if self.default_gateway == nil {
		return errors.New("Dont need to reset gateway!")
	}

	if self.is_gateway_changed == false {
		return
	}

	err = restoreGateWay(self.default_nic, self.default_gateway)
	if err != nil {
		return
	}

	self.is_gateway_changed = false
	return
}

