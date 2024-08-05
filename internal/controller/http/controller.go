package httpcontroller

import (
	"datapoint/internal/controller/http/dbcontroller"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

func New(
	r fiber.Router,
	v *validator.Validate,
	dbService dbcontroller.Service,
) {
	dbcontroller.New(r, dbService, v)
}
