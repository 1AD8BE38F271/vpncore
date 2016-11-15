/*
go test -v ./... -bench '^Bench*' -run ^$

BenchmarkAES128-4   	  100000	     11934 ns/op
BenchmarkAES192-4   	  100000	     13296 ns/op
BenchmarkAES256-4   	  100000	     14143 ns/op
BenchmarkTEA-4      	   50000	     34950 ns/op
BenchmarkSimpleXOR-4	 2000000	       774 ns/op
BenchmarkBlowfish-4 	   20000	     59895 ns/op
BenchmarkNone-4     	20000000	       109 ns/op
BenchmarkCast5-4    	   20000	     69818 ns/op
BenchmarkTripleDES-4	    2000	   1008371 ns/op
BenchmarkTwofish-4  	   20000	     89424 ns/op
BenchmarkXTEA-4     	   20000	     80140 ns/op
BenchmarkSalsa20-4  	  300000	      5118 ns/op

 */


package enc

import (
	"bytes"
	mrand "math/rand"
	crand "crypto/rand"
	"io"
	"testing"
    . "github.com/1AD8BE38F271/vpncore"
)

func EnryptionOne(t *testing.T, encrytion Cipher, testKey string, testDataLen int) {
	Logger.Infof("Test %s for EnryptionOne with key[%s] and test data length %d\n", encrytion, testKey, testDataLen)

	bc, err := NewBlock(&BlockConfig{Cipher:encrytion, Password:testKey})
	if err != nil {
		t.Fatal(err)
	}
	data := make([]byte, testDataLen)
	io.ReadFull(crand.Reader, data)
	dec := make([]byte, testDataLen)
	enc := make([]byte, testDataLen)
	bc.Encrypt(enc, data)
	bc.Decrypt(dec, enc)
	if !bytes.Equal(data, dec) {
		t.Fail()
	}
}

func EnryptionStreaming(t *testing.T, encrytion Cipher, testKey string, testDataLen int) {

	bc1, err := NewBlock(&BlockConfig{Cipher:encrytion, Password:testKey})
	bc2, err := NewBlock(&BlockConfig{Cipher:encrytion, Password:testKey})

	if err != nil {
		t.Fatal(err)
	}

	len1 := mrand.Intn(testDataLen) + testDataLen
	len2 := mrand.Intn(testDataLen) + testDataLen
	len3 := mrand.Intn(testDataLen) + testDataLen
	Logger.Infof("Test %s for EnryptionStreaming with data length %d-%d-%d\n", encrytion, len1, len2, len3)

	data1 := make([]byte, len1)
	data2 := make([]byte, len2)
	data3 := make([]byte, len3)

	io.ReadFull(crand.Reader, data1)
	io.ReadFull(crand.Reader, data2)
	io.ReadFull(crand.Reader, data3)

	alldata := make([]byte, len1+len2+len3)
	copy(alldata[:len1], data1)
	copy(alldata[len1:len1+len2], data2)
	copy(alldata[len1+len2:], data3)

	dec := make([]byte, len(alldata))
	dec2 := make([]byte, len(alldata))
	enc := make([]byte, len(alldata))
	enc2 := make([]byte, len(alldata))

	bc1.Encrypt(enc[:len1], data1)
	bc1.Encrypt(enc[len1:len1+len2], data2)
	bc1.Encrypt(enc[len1+len2:], data3)
	bc2.Encrypt(enc2, alldata)

	if !bytes.Equal(enc2, enc) {
		t.Fatalf("Not streaming consistent encryption!")
	}

	bc2.Decrypt(dec2[:len1], enc2[:len1])
	bc2.Decrypt(dec2[len1:len1+len2], enc2[len1:len1+len2])
	bc2.Decrypt(dec2[len1+len2:], enc2[len1+len2:])
	bc1.Decrypt(dec, enc)

	if !bytes.Equal(dec2, dec) {
		t.Fatalf("Error decryption 1!")
	}

	if !bytes.Equal(alldata, dec2) {
		t.Fatalf("Error decryption 2!")
	}
}

func TestAll(t *testing.T) {
	password := "I'm test key"

	testDatalens := []int{0x10, 0x100, 0x1000, 0x10000, 0x10000}
	testCiphers := []Cipher{AES128CFB, AES256CFB, SALSA20, NONE}

	for _, testDatalen := range testDatalens {
		for _, cipher := range testCiphers {
			EnryptionOne(t, cipher, password, testDatalen)
			EnryptionStreaming(t, cipher, password, testDatalen)

		}
	}

}

func BenchmarkSalsa20(b *testing.B) {
	var testDataLen = 2047

	pass := make([]byte, 32)
	io.ReadFull(crand.Reader, pass)
	bc, err := NewSalsa20BlockCrypt(pass)
	if err != nil {
		b.Fatal(err)
	}

	data := make([]byte, testDataLen)
	io.ReadFull(crand.Reader, data)
	dec := make([]byte, testDataLen)
	enc := make([]byte, testDataLen)

	for i := 0; i < b.N; i++ {
		bc.Encrypt(enc, data)
		bc.Decrypt(dec, enc)
	}
}
