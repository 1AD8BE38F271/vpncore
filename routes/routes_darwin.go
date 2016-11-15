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

// Handle virtual interfaces

package routes

import (
	"bufio"
	"strings"
	"errors"
	"bytes"
	"fmt"
	"net"
	"github.com/1AD8BE38F271/vpncore/cmd"
)

func addRouteToNet(iface string, subnet net.IPNet, nextHop net.IP) (err error) {

	c := fmt.Sprintf("route add -net %s %s", subnet.String(), nextHop.String())
	_, err = cmd.RunCommand(c)

	if err != nil {
		return
	}
	return
}

func addRouteToHost(iface string, dest net.IP, nextHop net.IP) (err error) {

	c := fmt.Sprintf("route add -host %s %s", dest.String(), nextHop.String())
	_, err = cmd.RunCommand(c)

	if err != nil {
		return
	}
	return
}

func delNetRoute(dest net.IPNet) error {
	c := fmt.Sprintf("route delete -net %s", dest.String())
	_, err := cmd.RunCommand(c)

	if err != nil {
		return err
	}
	return nil
}

func delHostRoute(dest net.IP) error {
	c := fmt.Sprintf("route delete -host %s", dest.String())
	_, err := cmd.RunCommand(c)

	if err != nil {
		return err
	}
	return nil
}

func redirectGateway(iface string, gw net.IP) error {
	c := "route delete default"
	_, err := cmd.RunCommand(c)
	if err != nil {
		return err
	}

	c = fmt.Sprintf("route add -net 0.0.0.0 %s", gw.String())
	_, err = cmd.RunCommand(c)

	if err != nil {
		return err
	}

	return nil
}

func restoreGateWay(ifce string, gw net.IP) (err error) {
	c := "route delete default"
	_, err = cmd.RunCommand(c)
	if err != nil {
		return err
	}

	c = fmt.Sprintf("route add default %s", gw.String())
	_, err = cmd.RunCommand(c)

	if err != nil {
		return err
	}
	return
}

func (self *RoutesManager) GetCurrentNetGateway() (gw net.IP, dev string, err error) {

	out, err := cmd.RunCommand("netstat -rn")
	if err != nil {
		return
	}

	reader := bufio.NewReader(strings.NewReader(out))

	for {
		line, isPrefix, err := reader.ReadLine()

		if err != nil {
			break
		}
		if isPrefix {
			break
		}
		buf := bytes.NewBuffer(line)
		scanner := bufio.NewScanner(buf)
		scanner.Split(bufio.ScanWords)
		tokens := make([]string, 0, 6)

		for scanner.Scan() {
			tokens = append(tokens, scanner.Text())
		}

		if len(tokens) < 6 {
			continue
		}

		iface := tokens[5]
		dest := tokens[0]
		gateway := tokens[1]

		if bytes.Equal([]byte(dest), []byte("default")) {
			return net.ParseIP(gateway), iface, nil
		}

	}

	return nil, "", errors.New("No gateway found!")
}




