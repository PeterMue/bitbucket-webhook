package header

import (
	"crypto"
	"crypto/hmac"
	"encoding/hex"
	"errors"
	"fmt"
	"reflect"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

type Signature string

type Headers struct {
	RequestId string    `http:"X-Request-Id"`
	EventKey  string    `http:"X-Event-Key"`
	Signature Signature `http:"X-Hub-Signature"`
}

var DefaultHash crypto.Hash = crypto.SHA256

func (s Signature) hashName() string {
	str := string(s)
	if strings.Contains(str, "=") {
		return string(s[0:strings.IndexRune(str, '=')])
	}
	return ""
}

func (s Signature) Hash() (crypto.Hash, error) {
	name := s.hashName()
	if name == "" {
		return DefaultHash, nil
	}
	switch strings.ToUpper(name) {
	case "SHA256":
		return crypto.SHA256, nil
	case "SHA384":
		return crypto.SHA384, nil
	case "SHA512":
		return crypto.SHA512, nil
	}
	return DefaultHash, fmt.Errorf("Hash not supported")
}

func (s Signature) Digest() string {
	a := strings.SplitN(string(s), "=", 2)
	return a[len(a)-1]
}

func (s Signature) Validate(message []byte, secret string) (valid bool, err error) {
	hash, err := s.Hash()
	if err != nil {
		return false, err
	}
	mac := hmac.New(hash.New, []byte(secret))
	mac.Write(message)
	actual := mac.Sum(nil)
	expected, err := hex.DecodeString(s.Digest())
	if err != nil {
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
	for i := 0; i < t.NumField(); i++ {
		f := t.Field(i)
		if name, ok := f.Tag.Lookup("http"); ok {
			if tf := v.FieldByName(f.Name); tf.IsValid() && tf.CanSet() {
				tf.SetString(utils.ImmutableString(ctx.Get(name)))
			}
		}
	}
}

func New(ctx *fiber.Ctx) *Headers {
	h := &Headers{}
	h.parse(ctx)
	return h
}
