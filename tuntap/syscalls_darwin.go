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
	"unsafe"
	"golang.org/x/sys/unix"
	"os"
	"fmt"
	"encoding/binary"
	"syscall"
	"github.com/FTwOoO/vpncore/tcpip"
)

const (
	MAX_KCTL_NAME = 96
	SYSPROTO_CONTROL = 2
	UTUN_OPT_IFNAME = 2
	AF_SYS_CONTROL = 2        /* corresponding sub address type */
	UTUN_CONTROL_NAME = "com.apple.net.utun_control"
)

type ctl_info struct {
	ctl_id   uint32               /* Kernel Controller ID  */
	ctl_name [MAX_KCTL_NAME] byte /* Kernel Controller Name (a C string) */
}

type sockaddr_ctl struct {
	sc_len      byte   /* depends on size of bundle ID string */
	sc_family   byte   /* AF_SYSTEM */
	ss_sysaddr  uint16 /* AF_SYS_KERNCONTROL */
	sc_id       uint32 /* Controller unique identifier  */
	sc_unit     uint32 /* Developer private unit number */
	sc_reserved [5]uint32;
};

type UTunFile struct {
	file *os.File
}

func (utunF *UTunFile) Read(p []byte) (n int, err error) {
	buffer := make([]byte, 2000)
	num, err := utunF.file.Read(buffer)
	if err != nil {
		return
	}

	n = copy(p, buffer[4:num])
	return
}

func (utunF *UTunFile) Close() error {
	return utunF.file.Close()
}

func (utunF *UTunFile) Write(p []byte) (n int, err error) {
	if len(p) < 4 {
		err = fmt.Errorf("Error packet length %d length bytes %v", len(p), p)
		return
	}

	t := make([]byte, 4 + len(p))

	if tcpip.IsIPv4(p) {
		binary.BigEndian.PutUint32(t, syscall.AF_INET)
	} else if tcpip.IsIPv6(p) {
		binary.BigEndian.PutUint32(t, syscall.AF_INET6)
	}

	copy(t[4:], p)

	n, err = utunF.file.Write(t)

	return
}

func newTAP(ifName string) (ifce *Interface, err error) {
	if !checkTapName(ifName) {
		return nil, fmt.Errorf("Error name:%s", ifName)
	}
	file, err := os.OpenFile("/dev/" + ifName, os.O_RDWR, 0)
	if err != nil {
		return nil, err
	}

	ifce = &Interface{isTAP: true,
		ReadWriteCloser: file,
		name: ifName,
	}
	return
}

func newTUN(ifName string) (ifce *Interface, err error) {
	if !checkTunName(ifName) {
		return nil, fmt.Errorf("Error name:%s", ifName)
	}

	createdIfName, fd, err := openUTUN(ifName)
	if err != nil {
		return nil, err
	}

	ifce = &Interface{
		isTAP: false,
		ReadWriteCloser: &UTunFile{file:os.NewFile(fd, createdIfName)},
		name: createdIfName,
	}
	return
}

func openUTUN(ifName string) (createdTunName string, fd uintptr, err error) {
	//
	// https://github.com/OpenVPN/openvpn/blob/master/src/openvpn/tun.c:open_darwin_utun
	//
	var dev_index int

	dev_index, err = getTunTapIndex(ifName)
	if err != nil {
		return
	}

	fd, _, errno := unix.Syscall(unix.SYS_SOCKET, unix.AF_SYSTEM, unix.SOCK_DGRAM, SYSPROTO_CONTROL)
	if errno != 0 {
		err = errno
		return
	}

	var info ctl_info
	copy(info.ctl_name[:], []byte(UTUN_CONTROL_NAME))

	CTLIOCGINFO := uintptr(0xc0644e03)

	_, _, errno = unix.Syscall(unix.SYS_IOCTL, fd, CTLIOCGINFO, uintptr(unsafe.Pointer(&info)))
	if errno != 0 {

		err = errno
		return
	}

	var addr sockaddr_ctl
	addr.sc_family = unix.AF_SYSTEM
	addr.ss_sysaddr = AF_SYS_CONTROL
	addr.sc_id = info.ctl_id
	addr.sc_unit = (uint32)(dev_index)
	addr.sc_len = (byte)(unsafe.Sizeof(addr) & 0xFF)

	_, _, errno = unix.Syscall(unix.SYS_CONNECT, fd, uintptr(unsafe.Pointer(&addr)), unsafe.Sizeof(addr))
	if errno != 0 {
		err = errno
		return
	}

	var utunname [20]byte
	var utunname_len = unsafe.Sizeof(utunname)

	_, _, errno = unix.Syscall6(unix.SYS_GETSOCKOPT, fd, SYSPROTO_CONTROL, UTUN_OPT_IFNAME, uintptr(unsafe.Pointer(&utunname)), uintptr(unsafe.Pointer(&utunname_len)), 0)
	if errno != 0 {
		err = errno
		return
	}

	//When utunname returns "utun1", the len is 6, including the tailing '\0' for C string?
	//Fuck off!
	createdTunName = string(utunname[:utunname_len - 1])
	return
}

