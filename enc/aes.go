package enc

import (
	"crypto/cipher"
	"crypto/aes"
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha256"
)


// AESBlockCrypt implements BlockCrypt
type AESBlockCrypt struct {
	block  cipher.Block
	cfb    cipher.Stream
	cfbdec cipher.Stream
}

// NewAESBlockCrypt initates BlockCrypt by the given key
func NewAESBlockCrypt(key []byte) (BlockCrypt, error) {
	c := new(AESBlockCrypt)

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	iv := pbkdf2.Key(key, []byte(key), 1, block.BlockSize(), sha256.New)

	c.block = block
	c.cfb = cipher.NewCFBEncrypter(block, iv)
	c.cfbdec = cipher.NewCFBDecrypter(block, iv)

	return c, nil
}

// Encrypt implements Encrypt interface
func (c *AESBlockCrypt) Encrypt(dst, src []byte) {
	c.cfb.XORKeyStream(dst, src)

}

// Decrypt implements Decrypt interface
func (c *AESBlockCrypt) Decrypt(dst, src []byte) {
	c.cfbdec.XORKeyStream(dst, src)

}
