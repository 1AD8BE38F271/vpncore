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

package conn

import (
	"testing"
	"github.com/FTwOoO/vpncore/enc"
	"fmt"
	"time"
	"io"
	"bytes"
	crand "crypto/rand"
	mrand "math/rand"
	"sync"
)

func TestNewListener(t *testing.T) {
	proto := PROTO_TCP
	password := "123456"
	port := mrand.Intn(100) + 20000
	testDatalens := []int{0x10, 0x100, 0x1000, 0x10000, 0x10000}
	testCiphers := []enc.Cipher{enc.AES128CFB, enc.AES256CFB, enc.SALSA20, enc.NONE}

	for _, testDatalen := range testDatalens {
		for _, cipher := range testCiphers {
			fmt.Printf("Test PROTOCOL[%s] with ENCRYPTION[%s] PASS[%s] DATALEN[%d]\n", proto, cipher, password, testDatalen)
			testOneConnection(t, proto, cipher, port, password, testDatalen)

		}
	}



}

func  testOneConnection (t *testing.T, proto TransProtocol, cipher enc.Cipher, port int, password string, testDatalen int) {

	blockConfig := &enc.BlockConfig{Cipher:cipher, Password:password}
	l, err := NewListener(proto, fmt.Sprintf("0.0.0.0:%d", port), blockConfig)
	if err != nil {
		t.Fatal(err)
	}

	testData := make([]byte, testDatalen)
	io.ReadFull(crand.Reader, testData)
	fmt.Printf("Test data is %v...\n", testData[:0x10])

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()

		expectedData := make([]byte, testDatalen)

		connection, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}

		areadyRead := 0
		for {
			if areadyRead == testDatalen {
				if !bytes.Equal(expectedData, testData) {
					t.Fatal("Bytes does not equal!")
				}
				return
			}

			n, err := connection.Read(expectedData[areadyRead:])
			if err != nil {
				t.Fatal(err)
			}

			fmt.Printf("Read %d bytes: %v...\n", n, expectedData[areadyRead:areadyRead + 0x10])
			areadyRead += n
		}

	}()

	<-time.After(3 * time.Second)

	go func() {
		defer wg.Done()

		connection, err := Dial(proto, fmt.Sprintf("127.0.0.1:%d", port), blockConfig)
		if err != nil {
			t.Fatal(err)
		}

		areadyWrite := 0
		for {
			if areadyWrite == testDatalen {
				break
			}
			n, err := connection.Write(testData[areadyWrite:])
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf("Write %d bytes: %v...\n", n, testData[areadyWrite:areadyWrite + 0x10])
			areadyWrite += n
		}
	}()

	wg.Wait()
	l.Close()
}