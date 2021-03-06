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
 * Modification Author: 1AD8BE38F271 <1AD8BE38F271@protonmail.com>
 */
package tuntap

import (
	"testing"
	"net"
	"time"
	"fmt"
	"encoding/hex"
	"sync"
	"github.com/1AD8BE38F271/vpncore/tcpip"
	"github.com/1AD8BE38F271/vpncore/routes"

	"github.com/1AD8BE38F271/vpncore/cmd"
    . "github.com/1AD8BE38F271/vpncore"
)

const BUFFERSIZE = 1522

func startRead(wg *sync.WaitGroup, ch chan <- []byte, ifce *Interface) {
	wg.Add(1)
	defer wg.Done()

	for {
		buffer := make([]byte, BUFFERSIZE)
		n, err := ifce.Read(buffer)
		if err == nil {
			Logger.Infof("Received a packet(%d bytes from %s)\n", n, ifce.Name())
			buffer = buffer[:n:n]
			ch <- buffer
		} else {
			Logger.Error(err)
			return
		}
	}
}

func startPing(dst net.IP) {
	c := time.After(1 * time.Second)
	select {
	case <-c:
		c := fmt.Sprintf("ping -c 5 %s", dst.String())
		cmd.RunCommand(c)
		return
	}
}

func ip4BroadcastAddr(subnet net.IPNet) (brdIp net.IP) {

	brdIp = net.IP{0, 0, 0, 0}
	for i := 0; i < 4; i++ {
		brdIp[i] = subnet.IP[i] | (0xFF ^ subnet.Mask[i])
	}
	return

}

func testInterface(ifce *Interface, ip net.IP, subnet net.IPNet) {

	err := ifce.SetupNetwork(ip, nil, subnet, 1400)
	if err != nil {
		panic(err)
	}

	router, err := routes.NewRoutesManager()
	if err != nil {
		panic(err)
	}
	defer router.Destroy()

	err = router.SetNewGateway(ifce.Name(), ifce.IP())
	if err != nil {
		panic(err)
	}

	dataCh := make(chan []byte, 8)
	wg := &sync.WaitGroup{}
	go startRead(wg, dataCh, ifce)

	if ifce.IsTUN() {
		startPing(ifce.PeerIP())
	} else {
		startPing(ip4BroadcastAddr(ifce.Net()))
	}

	timeout := time.NewTimer(5 * time.Second).C

	readFrame:
	for {
		select {
		case buffer := <-dataCh:
			var ipPacket tcpip.IPv4Packet

			if ifce.IsTAP() {
				ethertype := tcpip.MACPacket(buffer).MACEthertype()
				if ethertype != tcpip.IPv4 {
					continue readFrame
				}
				if !tcpip.IsBroadcast(tcpip.MACPacket(buffer).MACDestination()) {
					continue readFrame
				}

				ipPacket = tcpip.IPv4Packet(tcpip.MACPacket(buffer).MACPayload())
			} else {
				ipPacket = tcpip.IPv4Packet(buffer)
			}

			if !tcpip.IsIPv4(ipPacket) {
				continue readFrame
			}

			if !ipPacket.SourceIP().Equal(ifce.IP()) {
				continue readFrame
			}
			if ipPacket.Protocol() != tcpip.ICMP {
				continue readFrame
			}
			Logger.Infof("Received ICMP frame: %#v\n", hex.EncodeToString(ipPacket))
			break readFrame

		case <-timeout:
			panic("Waiting for broadcast packet timeout")

		}
	}

	Logger.Infof("Close the iterface %s\n", ifce.Name())
	ifce.Close()
	wg.Wait()

}

func TestAll(t *testing.T) {
	subnet := net.IPNet{IP:[]byte{192, 168, 99, 0}, Mask:net.IPv4Mask(255, 255, 255, 0)}
	ip := net.IP{192, 168, 99, 1}

	ifce, err := NewTUN("tun2")
	if err != nil {
		t.Fatal(err)
	}

	Logger.Infof("create %s\n", ifce.Name())
	testInterface(ifce, ip, subnet)
	ifce.Close()

	ifce2, err := NewTAP("tap1")
	if err != nil {
		t.Fatal(err)
	}
	Logger.Infof("create %s\n", ifce2.Name())
	testInterface(ifce2, ip, subnet)
	ifce2.Close()

}
