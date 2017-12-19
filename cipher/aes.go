package cipher

import (
	"bytes"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/binary"
	"io"
	"unsafe"
)

// PCKS padding
func padKey(key []byte) []byte {
	var targetKeySize = 16
	if len(key) <= 16 {
		targetKeySize = 16
	} else if len(key) <= 24 {
		targetKeySize = 24
	} else {
		targetKeySize = 32
	}
	if r := targetKeySize - len(key); r > 0 {
		nkey := make([]byte, targetKeySize)
		copy(nkey, key)
		copy(nkey[len(key):], bytes.Repeat([]byte{byte(r)}, r))
		key = nkey
	}
	return key[:targetKeySize]
}

// source is left unaffected.
func AesEncrypt(source []byte, keyStr string) ([]byte, error) {
	srcLength := len(source)
	var usize uint64
	headSize := int(unsafe.Sizeof(usize))
	head := make([]byte, headSize, headSize+len(source))
	source = append(head, source...)
	len0 := len(source)
	if r := len0 % aes.BlockSize; r > 0 {
		len0 += aes.BlockSize - r
		padding := make([]byte, aes.BlockSize-r)
		if _, err := io.ReadFull(rand.Reader, padding); err != nil {
			return nil, err
		}
		source = append(source, padding...)
	}
	buffer := make([]byte, len0+aes.BlockSize)

	iv := buffer[:aes.BlockSize]
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	usize = uint64(srcLength)
	binary.BigEndian.PutUint64(source, usize)
	block, ec := aes.NewCipher(padKey([]byte(keyStr)))
	if ec != nil {
		return nil, ec
	}
	mode := cipher.NewCBCEncrypter(block, iv)
	mode.CryptBlocks(buffer[aes.BlockSize:], source)
	return buffer, nil
}

func AesDecrypt(encbuffer []byte, keyStr string) ([]byte, error) {
	iv := encbuffer[:aes.BlockSize]
	out := encbuffer[aes.BlockSize:]
	block, ec := aes.NewCipher(padKey([]byte(keyStr)))
	if ec != nil {
		return nil, ec
	}
	mode := cipher.NewCBCDecrypter(block, iv)
	mode.CryptBlocks(out, out)
	usize := binary.BigEndian.Uint64(out)
	var headSize = int(unsafe.Sizeof(usize))
	return out[headSize : headSize+int(usize)], nil
}
