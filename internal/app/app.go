package app

import (
	"datapoint/config"
	httpcontroller "datapoint/internal/controller/http"
	"datapoint/internal/repo/dbrepo"
	"datapoint/internal/service/dbservice"
	"datapoint/migration"
	"datapoint/pkg/database"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	jsoniter "github.com/json-iterator/go"
)

type structValidator struct{ v *validator.Validate }

func (s *structValidator) Validate(out any) error {
	return s.v.Struct(out)
}

func Run(cfg *config.Config) error {
	v := validator.New()

	app := fiber.New(fiber.Config{
		AppName:         "datapoint",
		JSONEncoder:     jsoniter.Marshal,
		JSONDecoder:     jsoniter.Unmarshal,
		StructValidator: &structValidator{v: v},
	})

	app.Use(cors.New(cors.Config{
		Next: func(c fiber.Ctx) bool { return false },
	}))

	db, err := database.New(cfg.DB.Driver, cfg.DB.DSN)
	if err != nil {
		return err
	}

	if err = migration.FromFile(db, "./migration/migration.sql"); err != nil {
		return err
	}

	dbRepo := dbrepo.New(db)

	dbService, err := dbservice.New(dbRepo, db)
	if err != nil {
		return err
	}

	httpcontroller.New(app, v, dbService)

	return app.Listen(cfg.HTTP.Addr)
}
