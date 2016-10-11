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
package dns

import (
	"net"
	"fmt"
	"bufio"
	"strings"
	"regexp"
	"net/textproto"
	"errors"
	"github.com/FTwOoO/vpncore/cmd"
)

func (self *DNSManager) SetupNewDNS(new_dns []net.IP) (err error) {

	old_dns, err := self.GetCurrentDNS()
	if err != nil {
		old_dns = nil
	}

	l := (DNSList)(new_dns)
	_, err = cmd.RunCommand(fmt.Sprintf("networksetup -setdnsservers WI-Fi %s", l.String()))
	if err != nil {
		return
	}

	self.old_dns = old_dns
	return
}

func (self *DNSManager) RestoreDNS() (err error) {

	if self.old_dns == nil {
		self.old_dns = "empty"
	}

	_, err = cmd.RunCommand(fmt.Sprintf("networksetup -setdnsservers WI-Fi %s", self.old_dns.String()))
	if err != nil {
		return
	}

	self.old_dns = nil
	return
}

func (self *DNSManager) GetCurrentDNS() (l DNSList, err error) {

	current_device, err := self.getActiveDevice()
	if err != nil {
		return
	}

	out, err := cmd.RunCommand("networksetup -listnetworkserviceorder")
	if err != nil {
		return
	}

	reader := bufio.NewReader(strings.NewReader(out))

	hardware := ""

	for {
		line, _, err := reader.ReadLine()

		if err != nil {
			return nil, err
		}

		rex, _ := regexp.Compile("^\\(Hardware\\s+Port\\:\\s*(.+)\\s*,\\s*Device\\:\\s*(.+)\\)\\s*")
		if ! rex.MatchString(string(line)) {
			continue
		}

		x := rex.FindStringSubmatch(string(line))
		if x[2] != current_device {
			continue
		} else {
			hardware = x[1]
			break
		}

	}

	if hardware == "" {
		return nil, fmt.Errorf("Cant get current DNS for device %s!", current_device)
	}

	out, err = cmd.RunCommand(fmt.Sprintf("networksetup -getdnsservers %s", hardware))

	ips := strings.Split(out, "\n")

	var current_dns []net.IP = []net.IP{}

	for _, ip := range ips {
		old_ip := net.ParseIP(textproto.TrimString(ip))
		if old_ip != nil {
			current_dns = append(current_dns, old_ip)

		}
	}

	if len(current_dns) > 0 {
		return current_dns, nil
	} else {
		return nil, fmt.Errorf("Cant get current DNS for device %s!", current_device)

	}
}

func (self *DNSManager) getActiveDevice() (device string, err error) {
	out, err := cmd.RunCommand("ifconfig")
	reader := bufio.NewReader(strings.NewReader(out))
	var current_device string

	for {
		line, isPrefix, err := reader.ReadLine()
		if isPrefix {
			break
		}

		if err != nil {
			return "", err
		}

		rex, _ := regexp.Compile("^\\s*status\\s*\\:\\s*active\\s*")
		if rex.MatchString(string(line)) {
			return current_device, nil
		}

		rex, _ = regexp.Compile("^\\s*(\\w+)\\s*\\:\\s*flags\\=.*")
		if rex.MatchString(string(line)) {
			current_device = rex.FindStringSubmatch(string(line))[1]
		}
	}

	return "", errors.New("Cant find active network device!")
}

