package header

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"reflect"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

type Signature string

type Headers struct {
	RequestId string `http:"X-Request-Id"`
	EventKey string `http:"X-Event-Key"`
	Signature Signature `http:"X-Hub-Signature"`
}

func (s Signature) Validate(message []byte, secret string) (valid bool, err error) {
	mac := hmac.New(sha256.New, []byte(secret))
	mac.Write(message)
	actual := mac.Sum(nil)
	expected, err := hex.DecodeString(string(s))
	if(err != nil) {
		valid = false
		return
	}
	if hmac.Equal(actual, expected) {
		valid = true
		err = nil
		return
	}
	valid = false
	err = errors.New("Signature mismatch")
	return
}

func (h *Headers) parse(ctx *fiber.Ctx) {
	t := reflect.TypeOf(h).Elem()
	v := reflect.ValueOf(h).Elem()
	for i := 0; i< t.NumField(); i++ {
		f := t.Field(i)
		if name, ok := f.Tag.Lookup("http"); ok {
			if tf := v.FieldByName(f.Name); tf.IsValid() && tf.CanSet() {
				tf.SetString(utils.ImmutableString(ctx.Get(name)))
			}
		}
	}
}

func New (ctx *fiber.Ctx) *Headers {
	h := &Headers{}
	h.parse(ctx)
	return h
}