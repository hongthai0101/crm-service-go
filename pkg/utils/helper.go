package utils

import (
	"encoding/json"
	"math/rand"
	"regexp"
	"time"
)

func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

func ContainsInt(s []int, e int) bool {
	for _, a := range s {
		if a == e {
			return true
		}
	}
	return false
}

func Random(min int, max int) int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(max-min) + min
}

func Pick(input interface{}, fields []string) map[string]interface{} {
	b, _ := json.Marshal(&input)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	output := make(map[string]interface{})
	for k, v := range m {
		if Contains(fields, k) {
			output[k] = v
		}
	}

	return output
}

func Omit(input interface{}, fields []string) map[string]interface{} {
	b, _ := json.Marshal(&input)
	var m map[string]interface{}
	_ = json.Unmarshal(b, &m)
	output := make(map[string]interface{})
	for k, v := range m {
		if !Contains(fields, k) {
			output[k] = v
		}
	}

	return output
}

func MaskStr(input string, startNo uint8, endNo int, maskChar string) string {
	start := input[:startNo]
	re := regexp.MustCompile(".")
	mask := input[startNo : len(input)-endNo]
	mask = re.ReplaceAllString(mask, maskChar)
	end := input[endNo:]

	return start + mask + end
}

func Map[K, V any](s []K, transform func(K) V) []V {
	rs := make([]V, len(s))
	for i, v := range s {
		rs[i] = transform(v)
	}
	return rs
}

type SignedInteger interface {
	int | int8 | int16 | int32 | int64
}

func Smallest[T SignedInteger](s []T) T {
	r := s[0]
	for _, v := range s[1:] {
		if r > v {
			r = v
		}
	}
	return r
}

func TypeConverter[R any](data any) (*R, error) {
	var result R
	b, err := json.Marshal(&data)
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(b, &result)
	if err != nil {
		return nil, err
	}
	return &result, err
}

func Get[T interface{}](inputs []T, index int) *T {
	if len(inputs) > index {
		return &inputs[index]
	}
	return nil
}
