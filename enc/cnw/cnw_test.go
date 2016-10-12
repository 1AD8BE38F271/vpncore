/*
GoVPN -- simple secure free software virtual private network daemon
Copyright (C) 2014-2016 Sergey Matveev <stargrave@stargrave.org>

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU General Public License as published by
the Free Software Foundation, either version 3 of the License, or
(at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU General Public License for more details.

You should have received a copy of the GNU General Public License
along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package cnw

import (
	"bytes"
	"crypto/rand"
	"encoding/binary"
	"io"
	"testing"
	"testing/quick"
)

var (
	testCnwKey *[32]byte = new([32]byte)
)

func init() {
	io.ReadFull(rand.Reader, testCnwKey[:])
}

func TestChaffWindowSymmetric(t *testing.T) {
	nonce := make([]byte, 8)
	f := func(data []byte, pktNum uint64) bool {
		if len(data) == 0 {
			return true
		}
		binary.BigEndian.PutUint64(nonce, pktNum)
		chaffed := Chaff(testCnwKey, nonce, data)
		if len(chaffed) != len(data)*EnlargeFactor {
			return false
		}
		decoded, err := Winnow(testCnwKey, nonce, chaffed)
		if err != nil {
			return false
		}
		return bytes.Compare(decoded, data) == 0
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestChaffWindowSmallSize(t *testing.T) {
	_, err := Winnow(testCnwKey, []byte("foobar12"), []byte("foobar"))
	if err == nil {
		t.Fail()
	}
}

func BenchmarkChaff(b *testing.B) {
	nonce := make([]byte, 8)
	data := make([]byte, 16)
	io.ReadFull(rand.Reader, nonce)
	io.ReadFull(rand.Reader, data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Chaff(testCnwKey, nonce, data)
	}
}

func BenchmarkWinnow(b *testing.B) {
	nonce := make([]byte, 8)
	data := make([]byte, 16)
	io.ReadFull(rand.Reader, nonce)
	io.ReadFull(rand.Reader, data)
	chaffed := Chaff(testCnwKey, nonce, data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Winnow(testCnwKey, nonce, chaffed)
	}
}
