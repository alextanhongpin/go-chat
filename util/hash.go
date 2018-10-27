package util

import (
	"crypto/md5"
	"encoding/base64"
)

func HashID(id string) (string, error) {
	h := md5.New()
	_, err := h.Write([]byte(id))
	if err != nil {
		return "", err
	}
	hash := h.Sum(nil)
	return base64.URLEncoding.EncodeToString(hash), nil
}
