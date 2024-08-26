package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/go-resty/resty/v2"

	"test/generated/client"
	"test/generated/message"
	"test/generated/models"
)

func noError(err error) {
	if err != nil {
		panic(err)
	}
}

func empty[T any](v *T) {
	if v != nil {
		panic("empty")
	}
}

func notEmpty[T any](v *T) {
	if v == nil {
		panic("not empty")
	}
}

func equal(a, b any) {
	aJson, err := json.Marshal(a)
	noError(err)
	bJson, err := json.Marshal(b)
	noError(err)

	aStr := string(aJson)
	bStr := string(bJson)

	if aStr != bStr {
		panic(fmt.Sprintf("%s != %s", aStr, bStr))
	}
}

var rd *_rd = &_rd{
	rand: rand.New(rand.NewSource(time.Now().UnixNano())),
}

type _rd struct {
	rand *rand.Rand
}

func (rd *_rd) bool() bool {
	return rd.rand.Intn(2) == 1
}

func (rd *_rd) int() int {
	return rd.rand.Int()
}

func (rd *_rd) float64() float64 {
	return rd.rand.Float64()
}

func (rd *_rd) string() string {
	return strconv.FormatInt(rd.rand.Int63(), 16)
}

func list[T any](f func() T, n int) (r []T) {
	r = make([]T, 0, n)

	for i := 0; i < n; i++ {
		r = append(r, f())
	}

	return r
}

func optional[T any](f func() T) *T {
	if rd.bool() {
		return nil
	}
	v := f()
	return &v
}

func main() {
	base := resty.New()
	base.SetBaseURL("http://127.0.0.1:30435/")

	c := client.New(base)

	// 																			op1

	{
		rspE := &message.Op1Rsp200{
			Uri1: rd.string(),
			Uri2: rd.int(),
			Qry1: rd.string(),
			Qry2: rd.int(),
			Qryo: nil,
			Req1: rd.int(),
			Req2: list(rd.string, 16),
		}

		rspA, err := c.TestTag1.Op1(
			&message.Op1Uri{
				Uri1: rspE.Uri1,
				Uri2: rspE.Uri2,
			},
			&message.Op1Qry{
				Qry1: rspE.Qry1,
				Qry2: rspE.Qry2,
			},
			&message.Op1Req{
				Req1: rspE.Req1,
				Req2: rspE.Req2,
			},
		)
		noError(err)
		equal(rspE, rspA)
	}

	// 																			op2
	{
		rdObj2 := func() models.Object2 {
			return models.Object2{
				RequiredField: rd.string(),
				OptionalField: optional(rd.string),
			}
		}

		rspE := &message.Op2Rsp200{
			StringField: rd.string(),
			IntField:    rd.int(),
			FloatField:  rd.float64(),
			ArrayField1: list(rd.int, 16),
			ArrayField2: list(rdObj2, 4),
			ObjectField1: struct {
				IntField *int "json:\"intField,omitempty\""
			}{
				IntField: optional(rd.int),
			},
			ObjectField2: rdObj2(),
		}

		rspA, err := c.TestTag1.Op2(
			&message.Op2Uri{
				Uri1: rd.string(),
				Uri2: rd.int(),
			},
			&message.Op2Qry{
				Qry1: rd.string(),
			},
			rspE,
		)
		noError(err)
		equal(rspE, rspA)
	}

	// 																			op3
	{
		rspE := rd.string()

		rspA, err := c.TestTag1.Op3(
			&message.Op2Uri{
				Uri1: rd.string(),
				Uri2: rd.int(),
			},
			&message.Op2Qry{
				Qry1: rd.string(),
			},
			&rspE,
		)
		noError(err)
		equal(rspE, rspA)
	}

	// 																			op4
	{
		rsp, err := c.TestTag2.Op4()
		noError(err)
		notEmpty(rsp)
	}

	// 																			op5
	{
		rsp, err := c.TestTag2.Op5()
		noError(err)
		notEmpty(rsp)
	}

	// 																			op6
	{
		code := 200 + rd.rand.Intn(4)

		rsp200, rsp201, rsp202, rsp203, err := c.TestTag2.Op6(&message.Op6Req{
			Code: code,
		})

		noError(err)
		switch code {
		case 200:
			notEmpty(rsp200)
			empty(rsp201)
			empty(rsp202)
			empty(rsp203)
		case 201:
			empty(rsp200)
			notEmpty(rsp201)
			empty(rsp202)
			empty(rsp203)
		case 202:
			empty(rsp200)
			empty(rsp201)
			notEmpty(rsp202)
			empty(rsp203)
		case 203:
			empty(rsp200)
			empty(rsp201)
			empty(rsp202)
			notEmpty(rsp203)
		}
	}
}
