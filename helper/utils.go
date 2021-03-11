package helper

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"github.com/gorilla/mux"
	"math/rand"
	"net/http"
	"strconv"
	"time"
)

func RandomName() string {
	return Hash(
		Md5HashSum(),
		strconv.FormatInt(time.Now().UnixNano()+rand.Int63(), 10),
	)
}

func Hash(h hasher, str string) string {
	return hex.EncodeToString(h([]byte(str)))
}

type hasher func([]byte) []byte

func Md5HashSum() hasher {
	return func(bytes []byte) []byte {
		hasher := md5.New()
		hasher.Write(bytes)
		return hasher.Sum(nil)
	}
}

func GetRequestParam(r *http.Request, param string) (uint, error) {
	vars := mux.Vars(r)
	if param, ok := vars[param]; ok {
		if uid, err := strconv.ParseUint(param, 10, 64); err == nil {
			return uint(uid), nil
		}
	}
	return 0, errors.New("parse error")
}
