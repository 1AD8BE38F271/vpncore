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

type Ethertype [2]byte

// Ethertype values. From: http://en.wikipedia.org/wiki/Ethertype
var (
	IPv4 = Ethertype{0x08, 0x00}
	ARP = Ethertype{0x08, 0x06}
	WakeOnLAN = Ethertype{0x08, 0x42}
	TRILL = Ethertype{0x22, 0xF3}
	DECnetPhase4 = Ethertype{0x60, 0x03}
	RARP = Ethertype{0x80, 0x35}
	AppleTalk = Ethertype{0x80, 0x9B}
	AARP = Ethertype{0x80, 0xF3}
	IPX1 = Ethertype{0x81, 0x37}
	IPX2 = Ethertype{0x81, 0x38}
	QNXQnet = Ethertype{0x82, 0x04}
	IPv6 = Ethertype{0x86, 0xDD}
	EthernetFlowControl = Ethertype{0x88, 0x08}
	IEEE802_3 = Ethertype{0x88, 0x09}
	CobraNet = Ethertype{0x88, 0x19}
	MPLSUnicast = Ethertype{0x88, 0x47}
	MPLSMulticast = Ethertype{0x88, 0x48}
	PPPoEDiscovery = Ethertype{0x88, 0x63}
	PPPoESession = Ethertype{0x88, 0x64}
	JumboFrames = Ethertype{0x88, 0x70}
	HomePlug1_0MME = Ethertype{0x88, 0x7B}
	IEEE802_1X = Ethertype{0x88, 0x8E}
	PROFINET = Ethertype{0x88, 0x92}
	HyperSCSI = Ethertype{0x88, 0x9A}
	AoE = Ethertype{0x88, 0xA2}
	EtherCAT = Ethertype{0x88, 0xA4}
	EthernetPowerlink = Ethertype{0x88, 0xAB}
	LLDP = Ethertype{0x88, 0xCC}
	SERCOS3 = Ethertype{0x88, 0xCD}
	HomePlugAVMME = Ethertype{0x88, 0xE1}
	MRP = Ethertype{0x88, 0xE3}
	IEEE802_1AE = Ethertype{0x88, 0xE5}
	IEEE1588 = Ethertype{0x88, 0xF7}
	IEEE802_1ag = Ethertype{0x89, 0x02}
	FCoE = Ethertype{0x89, 0x06}
	FCoEInit = Ethertype{0x89, 0x14}
	RoCE = Ethertype{0x89, 0x15}
	CTP = Ethertype{0x90, 0x00}
	VeritasLLT = Ethertype{0xCA, 0xFE}
)
