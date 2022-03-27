package base

import (
	"math/big"
	"strconv"
	"strings"
	"time"
)

type BaseRows struct {
	More    bool   `json:"more"`
	NextKey string `json:"next_key"`
}

type Authorization struct {
	Actor      string `json:"actor"`
	Permission string `json:"permission"`
}

type BaseAction struct {
	Account       string          `json:"account"`
	Name          string          `json:"name"`
	Authorization []Authorization `json:"authorization"`
}

type ExtendedAsset struct {
	Quantity Asset  `json:"quantity"`
	Contract string `json:"contract"`
}

type BigInt struct {
	*big.Int
}

func (w *BigInt) UnmarshalJSON(i []byte) error {
	s := string(i)
	if s == "null" {
		return nil
	}

	var z big.Int
	if s[0] == '"' {
		err := z.UnmarshalText(i[1 : len(i)-1])
		if err != nil {
			return err
		}
	} else {
		err := z.UnmarshalText(i)
		if err != nil {
			return err
		}
	}

	w.Int = &z
	return nil
}

type Bool bool

func (b *Bool) UnmarshalJSON(i []byte) error {
	s := string(i)
	if s == "null" {
		return nil
	}

	bo, err := strconv.ParseBool(s)
	if err != nil {
		return err
	}

	*b = Bool(bo)
	return nil
}

type UInt64 uint64

func (w *UInt64) UnmarshalJSON(i []byte) error {
	s := string(i)
	if s == "null" {
		return nil
	}

	if s[0] == '"' {
		z, err := strconv.ParseUint(s[1:len(s)-1], 10, 64)
		if err != nil {
			return err
		}

		*w = UInt64(z)
	} else {
		z, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}

		*w = UInt64(z)
	}

	return nil
}

type Int64 int64

func (w *Int64) UnmarshalJSON(i []byte) error {
	s := string(i)
	if s == "null" {
		return nil
	}

	if s[0] == '"' {
		z, err := strconv.ParseInt(s[1:len(s)-1], 10, 64)
		if err != nil {
			return err
		}

		*w = Int64(z)
	} else {
		z, err := strconv.ParseUint(s, 10, 64)
		if err != nil {
			return err
		}

		*w = Int64(z)
	}

	return nil
}

type Float32 float32

func (w *Float32) UnmarshalJSON(i []byte) error {
	s := string(i)
	if s == "null" {
		return nil
	}

	if s[0] == '"' {
		z, err := strconv.ParseFloat(s[1:len(s)-1], 32)
		if err != nil {
			return err
		}

		*w = Float32(z)
	} else {
		z, err := strconv.ParseFloat(s, 32)
		if err != nil {
			return err
		}

		*w = Float32(z)
	}

	return nil
}

type Float64 float64

func (w *Float64) UnmarshalJSON(i []byte) error {
	s := string(i)
	if s == "null" {
		return nil
	}

	if s[0] == '"' {
		z, err := strconv.ParseFloat(s[1:len(s)-1], 64)
		if err != nil {
			return err
		}

		*w = Float64(z)
	} else {
		z, err := strconv.ParseFloat(s, 64)
		if err != nil {
			return err
		}

		*w = Float64(z)
	}

	return nil
}

type Asset string

func (a Asset) Quantity() float64 {
	f, _ := strconv.ParseFloat(a.RawQuantity(), 64)
	return f
}

func (a Asset) RawQuantity() string {
	return strings.TrimSpace(strings.Split(string(a), " ")[0])
}

func (a Asset) Precision() int {
	splitted := strings.SplitN(a.RawQuantity(), ".", 2)
	if len(splitted) < 2 {
		return 0
	}

	return len(strings.TrimSpace(splitted[1]))
}

func (a Asset) Symbol() string {
	split := strings.SplitN(string(a), " ", 2)
	if len(split) < 2 {
		return ""
	}
	return strings.TrimSpace(split[1])
}

type Symbol string

func (s Symbol) Name() string {
	splitted := strings.SplitN(string(s), ",", 2)
	if len(splitted) < 2 {
		return ""
	}

	return strings.TrimSpace(splitted[1])
}

func (s Symbol) Precision() int {
	splitted := strings.SplitN(string(s), ",", 2)
	if len(splitted) < 2 {
		return 0
	}

	n, _ := strconv.Atoi(strings.TrimSpace(splitted[0]))
	return n
}

type TimePoint time.Time

func (w *TimePoint) UnmarshalJSON(i []byte) error {
	s := string(i)
	if s == "null" {
		return nil
	}

	t, err := time.ParseInLocation("2006-01-02T15:04:05.999", s[1:len(s)-1], time.UTC)
	if err != nil {
		return err
	}

	*w = TimePoint(t)

	return nil
}

type TimePointSec time.Time

func (w *TimePointSec) UnmarshalJSON(i []byte) error {
	s := string(i)
	if s == "null" {
		return nil
	}

	t, err := time.ParseInLocation("2006-01-02T15:04:05", s[1:len(s)-1], time.UTC)
	if err != nil {
		return err
	}

	*w = TimePointSec(t)

	return nil
}

type BlockTimestampType TimePoint
