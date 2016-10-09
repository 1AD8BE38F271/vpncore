package enc

import "golang.org/x/crypto/salsa20"


// Salsa20BlockCrypt implements BlockCrypt
type Salsa20BlockCrypt struct {
	key [32]byte
}

// NewSalsa20BlockCrypt initates BlockCrypt by the given key
func NewSalsa20BlockCrypt(key []byte) (BlockCrypt, error) {
	c := new(Salsa20BlockCrypt)
	copy(c.key[:], key)
	return c, nil
}

// Encrypt implements Encrypt interface
func (c *Salsa20BlockCrypt) Encrypt(dst, src []byte) {
	salsa20.XORKeyStream(dst[:], src[:], c.key[:8], &c.key)
	copy(dst, src)
}

// Decrypt implements Decrypt interface
func (c *Salsa20BlockCrypt) Decrypt(dst, src []byte) {
	salsa20.XORKeyStream(dst[:], src[:], c.key[:8], &c.key)
	copy(dst, src)
}
