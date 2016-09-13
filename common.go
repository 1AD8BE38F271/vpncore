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
	"regexp"
	"strconv"
	"fmt"
	"strings"
	"os/exec"
)

func generatePeerIP(ip net.IP) (peerIp net.IP) {
	peerIp = net.IP(make([]byte, 4))
	copy([]byte(peerIp), []byte(ip.To4()))
	peerIp[3]++
	return
}

func checkTunName(ifName string) bool {
	rex, _ := regexp.Compile("tun\\d+")
	if ! rex.MatchString(ifName) {
		return false
	}
	return true
}

func checkTapName(ifName string) bool {
	rex, _ := regexp.Compile("tap\\d+")
	if ! rex.MatchString(ifName) {
		return false
	}
	return true
}

func getTunTapIndex(ifName string) (index int, err error) {
	rex, _ := regexp.Compile("tap|tun(\\d+)")
	if ! rex.MatchString(ifName) {
		return 0, fmt.Errorf("Error tun/tap name:%s", ifName)
	}

	ifNum := rex.FindStringSubmatch(ifName)[1]
	index, err = strconv.Atoi(ifNum)
	return

}

func runCommand(cmd string) (out string, err error) {

	args := strings.Split(cmd, " ")
	command := exec.Command(args[0], args[1:]...)

	out_bytes, err := command.CombinedOutput()
	out = string(out_bytes)
	if err != nil {
		fmt.Printf("RUN[%s]==>\n %s==================\n", cmd, out)
		return
	}
	fmt.Printf("RUN[%s]==>\n %s==================\n", cmd, out)
	return
}