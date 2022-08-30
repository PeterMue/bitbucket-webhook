package header

import (
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"
)

func TestHeaders(t *testing.T) {

	app := fiber.New()
	ctx := app.AcquireCtx(&fasthttp.RequestCtx{})

	ctx.Request().Header.Set("X-Request-Id", "foo")
	ctx.Request().Header.Set("X-Event-Key", "bar")
	ctx.Request().Header.Set("X-Hub-Signature", "baz")

	actual := New(ctx)
	expect := &Headers{
		RequestId: "foo",
		EventKey: "bar",
		Signature: Signature("baz"),
	}

	assert.Equal(t, actual, expect)
}

func TestSignature(t *testing.T) {
	sig := Signature("5aa613733e957bce8409d7e59b9851b4662147cbd72a76268868395b5acc4f10")
	if valid, err := sig.Validate([]byte("foobar"), "$ecret"); !valid {
		t.Errorf("Signature expect to be valid but was invalid: %s", err)
	}
}