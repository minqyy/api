package hasher

import (
	"crypto/sha256"
	"fmt"
)

type Hasher struct {
	salt []byte
}

func New(salt string) *Hasher {
	return &Hasher{salt: []byte(salt)}
}

func (h *Hasher) Create(buffer []byte) string {
	hash := sha256.New()

	hash.Write(append(buffer, h.salt...))
	return fmt.Sprintf("%x", hash.Sum(nil))
}
