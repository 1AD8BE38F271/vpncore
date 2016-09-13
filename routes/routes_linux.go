// +build linux

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
	"bytes"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"github.com/FTwOoO/vpncore/cmd"

)

func addRouteToHost(iface string, dest net.IP, nextHop net.IP) (err error) {

	c := fmt.Sprintf("ip -4 r a %s via %s dev %s", dest.String(), nextHop.String(), iface)
	_, err = cmd.Runcommand(c)

	if err != nil {
		return
	}

	return

}

func addRouteToNet(iface string, subnet net.IPNet, nextHop net.IP) (err error) {

	c := fmt.Sprintf("ip -4 route add %s via %s dev %s", subnet.String(), nextHop.String(), iface)
	_, err = cmd.Runcommand(c)

	if err != nil {
	}
	return
}

func delNetRoute(dest net.IPNet) (err error) {
	c := fmt.Sprintf("ip -4 route del %s", dest.String())
	_, err = cmd.Runcommand(c)

	if err != nil {
		return
	}
	return
}

func delHostRoute(dest net.IP) (err error) {
	c := fmt.Sprintf("ip -4 route del %s", dest.String())
	_, err = cmd.Runcommand(c)

	if err != nil {
		return
	}
	return
}

//TODO:
func redirectGateway(iface string, gw net.IP) (err error) {

	subnets := []string{"0.0.0.0/1", "128.0.0.0/1"}

	for _, subnet := range subnets {
		c := fmt.Sprintf("ip -4 route add %s via %s dev %s", subnet, gw.String(), iface)
		_, err = cmd.Runcommand(c)
		if err != nil {
			return err
		}
	}
	return nil
}

// TODO:
func restoreGateWay(ifce string, gw net.IP) (err error) {
	return nil
}

func (self *RoutesManager) GetCurrentNetGateway() (gw net.IP, dev string, err error) {

	file, err := os.Open("/proc/net/route")
	if err != nil {
		return nil, "", err
	}

	defer file.Close()
	rd := bufio.NewReader(file)

	s2byte := func(s string) byte {
		b, _ := strconv.ParseUint(s, 16, 8)
		return byte(b)
	}

	for {
		line, isPrefix, err := rd.ReadLine()

		if err != nil {
			//logger.Error(err.Error())
			return nil, "", err
		}
		if isPrefix {
			return nil, "", errors.New("Line Too Long!")
		}
		buf := bytes.NewBuffer(line)
		scanner := bufio.NewScanner(buf)
		scanner.Split(bufio.ScanWords)
		tokens := make([]string, 0, 8)

		for scanner.Scan() {
			tokens = append(tokens, scanner.Text())
		}

		iface := tokens[0]
		dest := tokens[1]
		gw := tokens[2]
		mask := tokens[7]

		if bytes.Equal([]byte(dest), []byte("00000000")) &&
			bytes.Equal([]byte(mask), []byte("00000000")) {
			a := s2byte(gw[6:8])
			b := s2byte(gw[4:6])
			c := s2byte(gw[2:4])
			d := s2byte(gw[0:2])

			ip := net.IPv4(a, b, c, d)

			return ip, iface, nil
		}

	}
	return nil, "", errors.New("No default gateway found")
}

