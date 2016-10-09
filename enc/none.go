package enc

// NoneBlockCrypt simple returns the plaintext
type NoneBlockCrypt struct {
	xortbl []byte
}

// NewNoneBlockCrypt initate by the given key
func NewNoneBlockCrypt(key []byte) (BlockCrypt, error) {
	return new(NoneBlockCrypt), nil
}

// Encrypt implements Encrypt interface
func (c *NoneBlockCrypt) Encrypt(dst, src []byte) {
	copy(dst, src)
}

// Decrypt implements Decrypt interface
func (c *NoneBlockCrypt) Decrypt(dst, src []byte) {
	copy(dst, src)
}
