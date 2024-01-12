package asset

import (
	"crypto/sha1"
	"github.com/gofrs/uuid"
	"strings"
)

func CreateUUIDForName(name string) string {
	data := make([]byte, 16)
	if name != "" {
		name = strings.ToLower(strings.ReplaceAll(name, "\\", "/"))
		dataSpan := []byte(name)
		sha := sha1.New()
		sha.Write(dataSpan)
		digest := sha.Sum(nil)

		for i := 0; i < 4; i++ {
			offset := i * 4
			data[offset] = digest[i*4]
			data[offset+1] = digest[i*4+1]
			data[offset+2] = digest[i*4+2]
			data[offset+3] = digest[i*4+3]
		}

		data[8] &= 0xBF
		data[8] |= 0x80

		data[6] &= 0x5F
		data[6] |= 0x50
	}

	u, _ := uuid.FromBytes(data)
	return u.String()
}
