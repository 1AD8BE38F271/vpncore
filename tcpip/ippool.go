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
	"errors"
	"net"
	"fmt"
	"encoding/binary"
	"sync"
)

type IP4Pool struct {
	Subnet   *net.IPNet
	ipMin    uint32
	ipMax    uint32

	pool     []bool
	poolLock sync.RWMutex
}

var poolFull = errors.New("IP Pool Full")

func NewIP4Pool(subnet *net.IPNet) (*IP4Pool, error) {
	p := new(IP4Pool)
	p.Subnet = subnet

	ip := subnet.IP.To4()
	if ip == nil {
		return nil, fmt.Errorf("Only support ipv4 :%v", subnet.IP)
	}

	var ipMin = net.IPv4Mask(0, 0, 0, 0)
	var ipMax = net.IPv4Mask(0, 0, 0, 0)

	for i := 0; i < 4; i++ {
		ipMin[i] = ip[i] & subnet.Mask[i]
		ipMax[i] = ip[i] | (^subnet.Mask[i])
	}

	ipMinValue := binary.BigEndian.Uint32(ipMin)
	ipMinValue += 1
	binary.BigEndian.PutUint32(ipMin, ipMinValue)
	ipMaxValue := binary.BigEndian.Uint32(ipMax)
	ipMaxValue -= 1

	//For A/B class, []bool will be large
	p.pool = make([]bool, ipMaxValue - ipMinValue + 1)
	p.ipMin = ipMinValue
	p.ipMax = ipMaxValue

	fmt.Printf("ipMin is %x\nipMax is %x\n", p.ipMin, p.ipMax)

	return p, nil

}

func (p *IP4Pool) Next() (net.IP, error) {
	p.poolLock.Lock()
	defer p.poolLock.Unlock()

	var i uint32
	for i = p.ipMin; i <= p.ipMax; i += 1 {
		if p.pool[i-p.ipMin] == false {
			fmt.Printf("current is %x\n", i)
			p.pool[i-p.ipMin] = true
			targetIp := net.IPv4(0, 0, 0, 0).To4()
			binary.BigEndian.PutUint32(targetIp, i)
			return targetIp, nil
		}
	}
	return nil, poolFull

}

func (p *IP4Pool) Release(ip net.IP) {
	val := binary.BigEndian.Uint32(ip)

	p.poolLock.Lock()
	defer p.poolLock.Unlock()
	p.pool[val - p.ipMin] = false
}
