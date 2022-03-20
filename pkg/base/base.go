package base

import (
	"math/big"
	"strconv"
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
	Quantity string `json:"quantity"`
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

type Float32 uint32

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
