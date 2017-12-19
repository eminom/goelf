package cipher

import (
	"bytes"
	"log"
	"testing"
)

func encTest(source []byte, key string, t *testing.T) {
	ebuff, err := AesEncrypt(source, key)
	if err != nil {
		panic(err)
	}
	dbuff, err := AesDecrypt(ebuff, key)
	if err != nil {
		panic(err)
	}
	log.Printf("<%v> : len = %v (encrypted length=%v)\n",
		string(dbuff),
		len(dbuff),
		len(ebuff),
	)
	log.Printf("<%v>\n", ebuff)

	if !bytes.Equal(source, dbuff) {
		t.Errorf("test for<%v> with key=<%v> failed", source, key)
	}
}

func TestAes(t *testing.T) {
	encTest([]byte("hasta victoria siempre"), "kl1", t)
	encTest([]byte(""), "a", t)
	encTest([]byte(""), "ss", t)
	encTest([]byte("a"), "", t)
	encTest([]byte("a"), "aa", t)
}
