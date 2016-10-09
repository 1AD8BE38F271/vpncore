package enc

import (
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha1"
	"errors"
)

type Cipher string

const (
	NONE = Cipher("None")
	SALSA20 = Cipher("salsa20")
	AES256CFB = Cipher("aes256cfb")
	AES128CFB = Cipher("aes128cfb")
	SALT = "i'm salt"
)


type BlockConfig struct {
	Cipher   Cipher
	Password string
}

// BlockCrypt defines encryption/decryption methods for a given byte slice
type BlockCrypt interface {
	// Encrypt encrypts the whole block in src into dst.
	// Dst and src may point at the same memory.
	Encrypt(dst, src []byte)

	// Decrypt decrypts the whole block in src into dst.
	// Dst and src may point at the same memory.
	Decrypt(dst, src []byte)
}


func GetKey(k string, kenLen int) []byte {
	pass := pbkdf2.Key([]byte(k), []byte(SALT), 4096, kenLen, sha1.New)
	return pass
}

func NewBlock(config *BlockConfig) (BlockCrypt, error){
	switch config.Cipher {
	case SALSA20:
		pass := GetKey(config.Password, 32)
		return NewSalsa20BlockCrypt(pass)
	case AES256CFB:
		pass := GetKey(config.Password, 32)
		return NewAESBlockCrypt(pass)
	case AES128CFB:
		pass := GetKey(config.Password, 16)
		return NewAESBlockCrypt(pass)
	case NONE:
		return NewNoneBlockCrypt([]byte{})
	default:
		return nil, errors.New("Invalid type!")
	}
}
