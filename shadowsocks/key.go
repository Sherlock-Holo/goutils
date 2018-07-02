package shadowsocks

import "crypto/md5"

func EncodeKey(key string) []byte {
	keyBytes := []byte(key)

	sum := md5.Sum(keyBytes)

	b := sum[:]

	sum = md5.Sum(append(sum[:], keyBytes...))

	b = append(b, sum[:]...)

	return b
}
