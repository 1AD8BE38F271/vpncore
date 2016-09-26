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
	"github.com/FTwOoO/go-enc"
	"fmt"
	"time"
	"io"
	"bytes"
	crand "crypto/rand"
	"sync"
)

func TestNewListener(t *testing.T) {
	proto := PROTO_TCP
	cipher := enc.SALSA20
	port := 20001
	password := "123456"
	testDataLen := 0x100000

	blockConfig := &enc.BlockConfig{Cipher:cipher, Password:password}
	l, err := NewListener(proto, fmt.Sprintf("0.0.0.0:%d", port), blockConfig)
	if err != nil {
		t.Fatal(err)
	}

	testData := make([]byte, testDataLen)
	io.ReadFull(crand.Reader, testData)
	fmt.Printf("Test data is %v...\n", testData[:0x10])

	var wg sync.WaitGroup
	wg.Add(2)

	go func(testData []byte, testDataLen int) {
		defer wg.Done()

		expectedData := make([]byte, testDataLen)

		connection, err := l.Accept()
		if err != nil {
			t.Fatal(err)
		}

		areadyRead := 0

		for {
			if areadyRead == testDataLen {
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

	}(testData, testDataLen)

	<-time.After(3 * time.Second)

	go func(testData []byte, testDataLen int) {
		defer wg.Done()

		connection, err := Dial(proto, fmt.Sprintf("127.0.0.1:%d", port), blockConfig)
		if err != nil {
			t.Fatal(err)
		}

		areadyWrite := 0

		for {
			if areadyWrite == testDataLen {
				break
			}
			n, err := connection.Write(testData[areadyWrite:])
			if err != nil {
				t.Fatal(err)
			}
			fmt.Printf("Write %d bytes: %v...\n", n, testData[areadyWrite:areadyWrite + 0x10])
			areadyWrite += n
		}
	}(testData, testDataLen)

	wg.Wait()

}