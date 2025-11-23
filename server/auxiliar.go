package server

import (
	"crypto/sha256"
	"encoding/hex"
)

const _SHORTEN_LEN = 10

// createShortenURL takes a string and returns a shorten version of it (hash)
func createShortenURL(url string) string {
	h := sha256.New()
	h.Write([]byte(url))
	hashBytes := h.Sum(nil)

	// Encode the hash to a hexadecimal string
	hexHash := hex.EncodeToString(hashBytes)

	// Return a truncated portion of the hex hash
	if len(hexHash) > _SHORTEN_LEN {
		return hexHash[:_SHORTEN_LEN]
	}
	return hexHash
}
