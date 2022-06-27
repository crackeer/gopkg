package util

import (
	"errors"
	"time"

	"github.com/speps/go-hashids"
)

func newTraceIDHashData() *hashids.HashIDData {
	hd := hashids.NewData()
	hd.Salt = "#trace-id-hashids-salt@realsee#"
	hd.MinLength = 32
	return hd
}

// ParseTraceID
//  @param hostIP
//  @return string
//  @return int
//  @return error
func ParseTraceID(traceID string) (string, int, error) {
	h, _ := hashids.NewWithData(newTraceIDHashData())
	result, err := h.DecodeInt64WithError(traceID)
	if err != nil {
		return "", 0, err
	}
	if len(result) < 3 {
		return "", 0, errors.New("illegal trace_id")
	}
	hostIP := IPIntToString(int(result[0]))
	return hostIP, int(result[1]), nil

}

// GenTraceID
//  @param hostIP
//  @return string
func GenTraceID(hostIP string) string {

	ipIntValue := StringIPToInt(hostIP)

	unixNano := time.Now().UnixNano()

	ts := unixNano / 1e9
	nonce := unixNano % 1e9 / 1e3

	h, _ := hashids.NewWithData(newTraceIDHashData())
	result, _ := h.Encode([]int{ipIntValue, int(ts), int(nonce)})
	return result

}
