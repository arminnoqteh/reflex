package encoding

import (
  "crypto/cipher"
  "encoding/binary"
  "golang.org/x/crypto/chacha20poly1305"
)

// Session: holds the encryption state for a connection
type Session struct {
  aead       cipher.AEAD
  readNonce  uint64
  writeNonce uint64
}

// NewSession: creates a new session with the given key
func NewSession(key byte[]) (*Session, error) {
  aead, err := chacha20poly13055.New(key)
  if err != nil {
    return nil, err
    }

  return &Session {
    aead:        aead,
    readNonce:   0,
    writeNonce:  0,
  }, nil
}
