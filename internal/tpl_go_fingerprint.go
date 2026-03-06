package internal

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
)

const generatedFingerprintPrefix = "// sbgen:fingerprint "

func withFingerprint(content []byte) []byte {
	header := generatedFingerprintPrefix + hashContent(content) + "\n"
	return append([]byte(header), content...)
}

func splitFingerprint(content []byte) ([]byte, string, bool) {
	lines := bytes.SplitN(content, []byte("\n"), 2)
	if len(lines) < 2 {
		return nil, "", false
	}
	if !bytes.HasPrefix(lines[0], []byte(generatedFingerprintPrefix)) {
		return nil, "", false
	}

	hash := string(lines[0][len(generatedFingerprintPrefix):])
	if len(hash) == 0 {
		return nil, "", false
	}
	return lines[1], hash, true
}

func hashContent(content []byte) string {
	sum := sha256.Sum256(content)
	return hex.EncodeToString(sum[:])
}
