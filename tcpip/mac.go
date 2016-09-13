/*
Package waterutil provides utility functions for interpreting TUN/TAP MAC farme headers and IP packet headers.
 It defines some constants such as protocol numbers and ethernet frame types.
Use waterutil along with package water to work with TUN/TAP interface data.

Frames/packets are interpreted in following format (as in TUN/TAP devices):

TAP - MAC Frame:
   No Tagging
  +-----------------------------------------------------------------------------
  | Octet |00|01|02|03|04|05|06|07|08|09|10|11|12|13|14|15|16|17|18|19|20|21|...
  +-----------------------------------------------------------------------------
  | Field | MAC Destination |   MAC  Source   |EType| Payload
  +-----------------------------------------------------------------------------

   Single-Tagged -- Octets [12,13] == {0x81, 0x00}
  +-----------------------------------------------------------------------------
  | Octet |00|01|02|03|04|05|06|07|08|09|10|11|12|13|14|15|16|17|18|19|20|21|...
  +-----------------------------------------------------------------------------
  | Field | MAC Destination |   MAC  Source   |    Tag    | Payload
  +-----------------------------------------------------------------------------

   Double-Tagged -- Octets [12,13] == {0x88, 0xA8}
  +-----------------------------------------------------------------------------
  | Octet |00|01|02|03|04|05|06|07|08|09|10|11|12|13|14|15|16|17|18|19|20|21|...
  +-----------------------------------------------------------------------------
  | Field | MAC Destination |   MAC  Source   | Outer Tag | Inner Tag | Payload
  +-----------------------------------------------------------------------------

TUN - IPv4 Packet:
  +---------------------------------------------------------------------------------------------------------------+
  |       | Octet |           0           |           1           |           2           |           3           |
  | Octet |  Bit  |00|01|02|03|04|05|06|07|08|09|10|11|12|13|14|15|16|17|18|19|20|21|22|23|24|25|26|27|28|29|30|31|
  +---------------------------------------------------------------------------------------------------------------+
  |   0   |   0   |  Version  |    IHL    |      DSCP       | ECN |                 Total  Length                 |
  +---------------------------------------------------------------------------------------------------------------+
  |   4   |  32   |                Identification                 | Flags  |           Fragment Offset            |
  +---------------------------------------------------------------------------------------------------------------+
  |   8   |  64   |     Time To Live      |       Protocol        |                Header Checksum                |
  +---------------------------------------------------------------------------------------------------------------+
  |  12   |  96   |                                       Source IP Address                                       |
  +---------------------------------------------------------------------------------------------------------------+
  |  16   |  128  |                                    Destination IP Address                                     |
  +---------------------------------------------------------------------------------------------------------------+
  |  20   |  160  |                                     Options (if IHL > 5)                                      |
  +---------------------------------------------------------------------------------------------------------------+
  |  24   |  192  |                                                                                               |
  |  30   |  224  |                                            Payload                                            |
  |  ...  |  ...  |                                                                                               |
  +---------------------------------------------------------------------------------------------------------------+

*/

package waterutil

import (
	"net"
)

type Tagging int
// Indicating whether/how a MAC frame is tagged. The value is number of bytes taken by tagging.
const (
	NotTagged Tagging = 0
	Tagged Tagging = 4
	DoubleTagged Tagging = 8
)

func IsBroadcast(addr net.HardwareAddr) bool {
	return addr[0] == 0xff && addr[1] == 0xff && addr[2] == 0xff && addr[3] == 0xff && addr[4] == 0xff && addr[5] == 0xff
}

func IsIPv4Multicast(addr net.HardwareAddr) bool {
	return addr[0] == 0x01 && addr[1] == 0x00 && addr[2] == 0x5e
}

type MACPacket []byte

func (macFrame MACPacket)MACDestination() net.HardwareAddr {
	return net.HardwareAddr(macFrame[:6])
}

func (macFrame MACPacket)MACSource() net.HardwareAddr {
	return net.HardwareAddr(macFrame[6:12])
}

func (macFrame MACPacket)MACTagging() Tagging {
	if macFrame[12] == 0x81 && macFrame[13] == 0x00 {
		return Tagged
	} else if macFrame[12] == 0x88 && macFrame[13] == 0xa8 {
		return DoubleTagged
	}
	return NotTagged
}

func (macFrame MACPacket)MACEthertype() Ethertype {
	ethertypePos := 12 + macFrame.MACTagging()
	return Ethertype{macFrame[ethertypePos], macFrame[ethertypePos + 1]}
}

func (macFrame MACPacket)MACPayload() []byte {
	return macFrame[12 + macFrame.MACTagging() + 2:]
}


