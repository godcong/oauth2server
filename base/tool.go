package base

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/http"
	"strings"
	"time"
)

type RandomKind int

const (
	T_RAND_NUM      RandomKind = iota // 纯数字
	T_RAND_LOWER                      // 小写字母
	T_RAND_UPPER                      // 大写字母
	T_RAND_LOWERNUM                   // 数字、小写字母
	T_RAND_UPPERNUM                   // 数字、大写字母
	T_RAND_ALL                        // 数字、大小写字母
)

var (
	RandomString = map[RandomKind]string{
		T_RAND_NUM:      "0123456789",
		T_RAND_LOWER:    "abcdefghijklmnopqrstuvwxyz",
		T_RAND_UPPER:    "ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		T_RAND_LOWERNUM: "0123456789abcdefghijklmnopqrstuvwxyz",
		T_RAND_UPPERNUM: "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ",
		T_RAND_ALL:      "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ",
	}
)

// 随机字符串
func GenerateRandomString2(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{{10, 48}, {26, 97}, {26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0

	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

func GenerateRandomString(size int, kind ...RandomKind) []byte {

	bytes := RandomString[T_RAND_ALL]
	if kind != nil {
		if k, b := RandomString[kind[0]]; b == true {
			bytes = k
		}
	}
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < size; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return result
}

func GeneratePassword(pass, salt string) string {
	return GenerateMD5(pass, salt)
}

func GenerateMD5(args ...string) string {
	m5 := md5.New()
	for _, v := range args {
		m5.Write([]byte(v))
	}
	st := m5.Sum(nil)
	return hex.EncodeToString(st)
}

//get remote ip address
func ObtainClientIP(r *http.Request) string {
	addr := r.Header.Get("remote_addr")
	if addr == "" {
		addr = r.RemoteAddr
	}

	if ips := strings.Split(addr, ":"); len(ips) > 0 && len(ips) <= 2 {
		return ips[0]
	}
	//[::]:1234 return 127.0.0.1
	return "127.0.0.1"

}

//only take len == 2 strings
func StitchAddress(args ...[]string) string {
	var tmp []string
	for _, v := range args {
		if len(v) == 2 {
			fmt.Println(v)
			tmp = append(tmp, strings.Join([]string{v[0], v[1]}, "="))
		}
	}
	return strings.Join(tmp, "&")
}

func JsonEncode(v interface{}) (string, error) {
	b, e := json.Marshal(v)
	return string(b), e
}

func JsonDecode(p string, v interface{}) error {
	if err := json.Unmarshal(([]byte)(p), &v); err != nil {
		return err
	}
	return nil
}
