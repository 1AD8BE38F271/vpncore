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

package aont

import (
	"bytes"
	"crypto/rand"
	"io"
	"testing"
	"testing/quick"
)

var (
	testAontKey *[AontKeySize]byte = new([AontKeySize]byte)
)

func init() {
	io.ReadFull(rand.Reader, testAontKey[:])
}

func TestAontSymmetric(t *testing.T) {
	f := func(data []byte) bool {
		encoded, err := AontEncode(testAontKey, data)
		if err != nil {
			return false
		}
		if len(encoded) != len(data)+ AontKeySize + AontHashSize {
			return false
		}
		decoded, err := AontDecode(encoded)
		if err != nil {
			return false
		}
		return bytes.Compare(decoded, data) == 0
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func TestAontSmallSize(t *testing.T) {
	_, err := AontDecode([]byte("foobar"))
	if err == nil {
		t.Fail()
	}
}

func TestTampered(t *testing.T) {
	f := func(data []byte, index int) bool {
		if len(data) == 0 {
			return true
		}
		encoded, _ := AontEncode(testAontKey, data)
		encoded[len(data)%index] ^= byte('a')
		_, err := AontDecode(encoded)
		if err == nil {
			return false
		}
		return true
	}
	if err := quick.Check(f, nil); err != nil {
		t.Error(err)
	}
}

func BenchmarkEncode(b *testing.B) {
	data := make([]byte, 128)
	io.ReadFull(rand.Reader, data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AontEncode(testAontKey, data)
	}
}

func BenchmarkDecode(b *testing.B) {
	data := make([]byte, 128)
	io.ReadFull(rand.Reader, data)
	encoded, _ := AontEncode(testAontKey, data)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		AontDecode(encoded)
	}
}
