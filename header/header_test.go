package header

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHeaders(t *testing.T) {

	Header := map[string][]string{
		"X-Request-Id":    {"foo"},
		"X-Event-Key":     {"bar"},
		"X-Hub-Signature": {"baz"},
	}

	actual := New(Header)
	expect := &Headers{
		RequestId: "foo",
		EventKey:  "bar",
		Signature: Signature("baz"),
	}

	assert.Equal(t, actual, expect)
}

func TestSignature(t *testing.T) {
	cases := []struct {
		sig    Signature
		expect bool
		err    error
	}{
		{Signature("4fcc06915b43d8a49aff193441e9e18654e6a27c2c428b02e8fcc41ccc2299f9"), true, nil},
		{Signature("sha256=4fcc06915b43d8a49aff193441e9e18654e6a27c2c428b02e8fcc41ccc2299f9"), true, nil},
		{Signature("sha384=8b7e7639ef66fa2583bf5fd1c08a0b4ed3a9c1ddbc380ddfd6fa35d9ce32e4dd69213994ed2f2cc750fa48221a189dfa"), true, nil},
		{Signature("sha512=ac76d1f21ab3affcab713dcec165cc517a1d9b79b1ac21fe99619fda7dfbee98b926080dc90117a8aa600875f4dbe7d50b0f13712bbfc9db8b57d7eddb91bc0c"), true, nil},
		{Signature("sha1=7f5c0e9cb2f07137b1c0249108d5c400a3c39be5"), false, fmt.Errorf("hash not supported")},
	}

	for _, c := range cases {
		actual, err := c.sig.Validate([]byte("foobar"), "secret")
		assert.Equal(t, actual, c.expect)
		assert.Equal(t, err, c.err)
	}
}
