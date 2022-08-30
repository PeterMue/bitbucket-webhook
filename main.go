package main

import (
	"log"
	"os"

	"github.com/PeterMue/bitbucket-webhook/config"
	"github.com/PeterMue/bitbucket-webhook/handler"
	"github.com/PeterMue/bitbucket-webhook/header"
	"github.com/gofiber/fiber/v2"
)

func ping(c *fiber.Ctx) error {
	c.Status(fiber.StatusOK).SendString(`{ "status" : "OK" }`)
	return nil
}

func main() {
	config, err := config.ParseFlags(os.Args);
	if err != nil {
		log.Fatalf("Invalid config: %s", err)
	}

	// Load configured event handler
	handlers := make(map[string]handler.Handler)
	for _, hook := range config.Hooks {
		handlers[hook.EventType] = *handler.New(hook.EventType, hook.Command, hook.Args)
	}

	app := fiber.New()
	app.Post("/webhook", func(c *fiber.Ctx) error {
		h := header.New(c)

		// Diagnostics works without signature check
		if h.EventKey == "diagnostics:ping" {
			return ping(c)
		}

		// Everything else needs a valid signature
		if valid, err := h.Signature.Validate(c.Body(), config.Secret); !valid {
			log.Printf("Signature verification falied: %s\n", err)
			return fiber.NewError(fiber.ErrBadRequest.Code, "Invalid signature")
		}

		// Go, find event handler
		handler, ok := handlers[h.EventKey]
		if !ok {
			return fiber.NewError(fiber.StatusNotFound, "No handler for given eventType configured")
		}
		
		// load body and run handler
		payload := c.Body()
		if err := handler.Run(h, payload); err != nil {
			return fiber.NewError(fiber.StatusInternalServerError, "Failed to complete webhook")
		}

		return c.Status(fiber.StatusOK).SendString("Ok")
	})
	log.Fatal(app.Listen(config.Listen))
}