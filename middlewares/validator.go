package middlewares

import (
	"reflect"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
)

var validate = validator.New()

// ValidateBody parses the JSON request body into the given model type, applies
// defaults and normalisations, then runs go-playground/validator constraints.
//
// Usage:
//
//	api.Post("/endpoint",
//	    middlewares.ValidateBody(&models.MyRequest{}),
//	    handler,
//	)
//
// The parsed struct is stored in c.Locals("body") and retrieved in the handler via:
//
//	req := c.Locals("body").(*models.MyRequest)
func ValidateBody(modelPtr interface{}) fiber.Handler {
	return func(c *fiber.Ctx) error {
		req := reflect.New(reflect.TypeOf(modelPtr).Elem()).Interface()

		if body := c.Body(); len(body) > 0 {
			if err := c.BodyParser(req); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error":   "Invalid request body",
					"details": err.Error(),
				})
			}
		}

		applyDefaults(req)
		normalizeFields(req)

		if err := validate.Struct(req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Validation failed",
				"details": err.Error(),
			})
		}

		c.Locals("body", req)
		return c.Next()
	}
}

// ─── Internal helpers ─────────────────────────────────────────────────────────

// applyDefaults sets zero-value fields to the value specified in the `default` struct tag.
func applyDefaults(ptr interface{}) {
	val := reflect.ValueOf(ptr)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return
	}
	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return
	}
	typ := elem.Type()
	for i := 0; i < elem.NumField(); i++ {
		fv := elem.Field(i)
		ft := typ.Field(i)
		if !fv.CanSet() {
			continue
		}
		tag := ft.Tag.Get("default")
		if tag == "" {
			continue
		}
		switch fv.Kind() {
		case reflect.String:
			if fv.String() == "" {
				fv.SetString(tag)
			}
		case reflect.Bool:
			if !fv.Bool() {
				if v, err := strconv.ParseBool(tag); err == nil {
					fv.SetBool(v)
				}
			}
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if fv.Int() == 0 {
				if v, err := strconv.ParseInt(tag, 10, 64); err == nil {
					fv.SetInt(v)
				}
			}
		}
	}
}

// normalizeFields applies field-level normalisations declared via the `normalize` tag.
// Supported values: "email" (trim + lowercase), "trim".
func normalizeFields(ptr interface{}) {
	val := reflect.ValueOf(ptr)
	if val.Kind() != reflect.Ptr || val.IsNil() {
		return
	}
	elem := val.Elem()
	if elem.Kind() != reflect.Struct {
		return
	}
	typ := elem.Type()
	for i := 0; i < elem.NumField(); i++ {
		fv := elem.Field(i)
		ft := typ.Field(i)
		if !fv.CanSet() {
			continue
		}
		tag := ft.Tag.Get("normalize")
		if tag == "" {
			continue
		}
		for _, norm := range strings.Split(tag, ",") {
			norm = strings.TrimSpace(norm)
			switch norm {
			case "email":
				setStringTransform(fv, func(s string) string {
					return strings.ToLower(strings.TrimSpace(s))
				})
			case "trim":
				setStringTransform(fv, strings.TrimSpace)
			}
		}
	}
}

func setStringTransform(fv reflect.Value, fn func(string) string) {
	switch fv.Kind() {
	case reflect.String:
		if s := fv.String(); s != "" {
			fv.SetString(fn(s))
		}
	case reflect.Ptr:
		if !fv.IsNil() && fv.Type().Elem().Kind() == reflect.String {
			v := fn(fv.Elem().String())
			fv.Elem().SetString(v)
		}
	}
}
