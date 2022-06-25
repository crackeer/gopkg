package util

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha1"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"math/rand"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"unicode"
)

var (
	md5Regex *regexp.Regexp
)

func init() {
	md5Regex = regexp.MustCompile(`[0-9A-F]`)
}

// SortStringMap ...
func SortStringMap(data map[string]string) [][]string {
	var (
		ret      [][]string
		rawSlice []string
	)
	for k := range data {
		rawSlice = append(rawSlice, k)
	}
	sort.Strings(rawSlice)

	for _, k := range rawSlice {
		v, _ := data[k]
		ret = append(ret, []string{k, v})
	}

	return ret
}

// MD5 ...
func MD5(input string) string {
	sum := md5.Sum([]byte(input))
	return hex.EncodeToString(sum[:])
}

// MD5Bytes ...
func MD5Bytes(input []byte) string {
	hash := md5.New()
	hash.Write(input)
	return hex.EncodeToString(hash.Sum(nil))
}

// Sha1 Return sha1 result for a string
func Sha1(input string) string {
	hash := sha1.Sum([]byte(input))
	return hex.EncodeToString(hash[:])
}

// Sha256 Return sha256 result for a string
func Sha256(input string) string {
	hash := sha256.Sum256([]byte(input))
	return hex.EncodeToString(hash[:])
}

// HMACSha1 ...
func HMACSha1(secret, input string) string {
	mac := hmac.New(sha1.New, []byte(secret))
	mac.Write([]byte(input))
	return string(hex.EncodeToString(mac.Sum(nil)))
}

// ParseVersion Parse version represent by string
// Only support 4 slice version
func ParseVersion(version string) int {
	var (
		versionInt = 0
		base       = 100
	)

	// Normalize to 4 slice version
	rawVersion := strings.Split(version, ".")
	lenOfVersion := len(rawVersion)
	if lenOfVersion > 4 {
		return versionInt
	}

	for i := 0; i < (4 - lenOfVersion); i++ {
		rawVersion = append(rawVersion, "0")
	}

	for idx, versionSegment := range rawVersion {
		versionSegmentInt, _ := strconv.Atoi(versionSegment)
		if versionSegmentInt <= 0 {
			versionSegmentInt = 0
		}

		switch idx {
		case 0:
			versionInt += base * base * base * versionSegmentInt
		case 1:
			versionInt += base * base * versionSegmentInt
		case 2:
			versionInt += base * versionSegmentInt
		case 3:
			versionInt += versionSegmentInt
		}
	}

	return versionInt
}

func Bs64Decode2StrMap(encodeStr string) (map[string]string, error) {
	if len(encodeStr) < 1 {
		return nil, errors.New("invalid encode string")
	}

	decodeBytes, err := base64.StdEncoding.DecodeString(encodeStr)
	if err != nil {
		return nil, err
	}

	params := map[string]string{}

	err = json.Unmarshal(decodeBytes, &params)

	if err != nil {
		return nil, err
	}

	return params, nil
}

// IsContain ...
func IsContain(strList string, str string, seq string) bool {
	rlist := strings.Split(strList, seq)
	for _, s := range rlist {
		if s == str {
			return true
		}
	}
	return false
}

// FuckOffURLProtocol fuck off protocol prefix for involving ST-API-ADAPT-PROJECT
func FuckOffURLProtocol(url string) string {
	protocols := []string{"http://", "https://"}

	for _, protocol := range protocols {
		if strings.Contains(url, protocol) {
			return strings.TrimPrefix(url, protocol)
		}
	}

	return url
}

var defaultLetters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func RandomString(n int, allowedChars ...[]rune) string {
	var letters []rune

	if len(allowedChars) == 0 {
		letters = defaultLetters
	} else {
		letters = allowedChars[0]
	}

	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}

	return string(b)
}

func ToNumber(str string) int {
	retData := 0
	for _, c := range str {
		retData += int(c)
	}
	return retData
}

// UcFirst 字符串首字母转化成大写
//
// Author : fuhaixu@ke.com<付海旭>
//
// Date : 下午11:03 2021/4/8
func UcFirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToUpper(v))
		res := u + str[len(u):]
		return res
	}

	return ""
}

// 字符串首字母转化成小写...
//
// Author : fuhaixu@ke.com<付海旭>
//
// Date : 下午7:35 2021/4/11
func LcFirst(str string) string {
	for _, v := range str {
		u := string(unicode.ToLower(v))
		res := u + str[len(u):]
		return res
	}

	return ""
}

// ToCamel ...
func ToCamel(str string) string {
	parts := strings.Split(str, "_")
	if len(parts) < 2 {
		return str
	}
	for i, v := range parts {
		if i == 0 {
			continue
		}
		parts[i] = UcFirst(v)
	}

	return strings.Join(parts, "")
}

// IsMD5String ...
func IsMD5String(s string) bool {
	if len(s) < 1 {
		return false
	}
	if len(md5Regex.FindAllString(s, -1)) == 32 {
		return true
	}
	return false
}

func ToString(src interface{}) string {
	switch v := src.(type) {
	case string:
		return v
	case []byte:
		return string(v)
	}
	rv := reflect.ValueOf(src)
	switch rv.Kind() {
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return strconv.FormatInt(rv.Int(), 10)
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return strconv.FormatUint(rv.Uint(), 10)
	case reflect.Float64:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 64)
	case reflect.Float32:
		return strconv.FormatFloat(rv.Float(), 'g', -1, 32)
	case reflect.Bool:
		return strconv.FormatBool(rv.Bool())
	}
	bytes, _ := Marshal(src)
	return string(bytes)
}
