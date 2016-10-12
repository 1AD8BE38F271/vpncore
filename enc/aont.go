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

// All-Or-Nothing-Transform, based on OAEP.
//
// This package implements OAEP (Optimal Asymmetric Encryption Padding)
// (http://cseweb.ucsd.edu/~mihir/papers/oaep.html)
// used there as All-Or-Nothing-Transformation
// (http://theory.lcs.mit.edu/~cis/pubs/rivest/fusion.ps).
// We do not fix OAEP parts length, instead we add hash-based
// checksum like in SAEP+
// (http://crypto.stanford.edu/~dabo/abstracts/saep.html).
//
// AONT takes 128-bit random r, data M to be encoded and produce the
// package PKG:
//
//     PKG = P1 || P2
//      P1 = Salsa20(key=r, nonce=0x00, M || BLAKE2b(r || M) )
//      P2 = BLAKE2b(P1) XOR r
package enc

import (
	"crypto/subtle"
	"errors"
	"golang.org/x/crypto/salsa20"
	"github.com/flynn/noise"
)

const (
	AontHashSize = 64
	AontKeySize = 32
)

var (
	dummyNonce []byte = make([]byte, 8)
)

// Encode the data, produce AONT package. Data size will be larger than
// the original one for 48 bytes.
func AontEncode(key *[AontKeySize]byte, in []byte) ([]byte, error) {
	out := make([]byte, len(in)+ AontHashSize + AontKeySize)
	copy(out, in)
	h := noise.HashBLAKE2b.Hash()
	h.Write(key[:])
	h.Write(in)
	copy(out[len(in):], h.Sum(nil))

	salsa20.XORKeyStream(out, out, dummyNonce, key)
	h.Reset()
	h.Write(out[:len(in)+ AontHashSize])
	for i, b := range h.Sum(nil)[:AontKeySize] {
		out[len(in)+ AontHashSize +i] = b ^ key[i]
	}
	return out, nil
}

// Decode the data from AONT package. Data size will be smaller than the
// original one for 48 bytes.
func AontDecode(in []byte) ([]byte, error) {
	if len(in) < AontHashSize + AontKeySize {
		return nil, errors.New("Too small input buffer")
	}
	h := noise.HashBLAKE2b.Hash()
	h.Write(in[:len(in)- AontKeySize])
	salsaKey := new([AontKeySize]byte)
	for i, b := range h.Sum(nil)[:AontKeySize] {
		salsaKey[i] = b ^ in[len(in)- AontKeySize +i]
	}
	h.Reset()
	h.Write(salsaKey[:AontKeySize])
	out := make([]byte, len(in)- AontKeySize)
	salsa20.XORKeyStream(out, in[:len(in)- AontKeySize], dummyNonce, salsaKey)
	h.Write(out[:len(out)- AontHashSize])
	if subtle.ConstantTimeCompare(h.Sum(nil), out[len(out)- AontHashSize:]) != 1 {
		return nil, errors.New("Invalid checksum")
	}
	return out[:len(out)- AontHashSize], nil
}
