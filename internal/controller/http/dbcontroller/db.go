package dbcontroller

import (
	"context"
	"datapoint/internal/controller/http/converter"
	"datapoint/internal/controller/http/model"
	"datapoint/internal/model/dbmodel"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

type Service interface {
	GetList() []*dbmodel.DB
	Add(ctx context.Context, info dbmodel.Info) (string, error)
	Edit(ctx context.Context, info dbmodel.Info, id string) error
	Delete(ctx context.Context, id string) error
	TableList(ctx context.Context, id string) ([]*dbmodel.Table, error)
	FunctionList(ctx context.Context, id string) ([]*dbmodel.Function, error)
}

type controller struct {
	s Service
	v *validator.Validate
}

func (c *controller) getList(ctx fiber.Ctx) error {
	return ctx.JSON(converter.ToDBList(c.s.GetList()))
}

func (c *controller) add(ctx fiber.Ctx) error {
	var body model.DBInfo

	err := ctx.Bind().JSON(&body)
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var id string
	if id, err = c.s.Add(ctx.Context(), converter.FromDBInfo(body)); err != nil {
		return err
	}

	return ctx.
		Status(fiber.StatusCreated).
		SendString(id)
}

func (c *controller) edit(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	err := c.v.Var(id, "uuid")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var body model.DBInfo
	if err = ctx.Bind().JSON(&body); err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	if err = c.s.Edit(ctx.Context(), converter.FromDBInfo(body), id); err != nil {
		return err
	}

	return nil
}

func (c *controller) delete(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	err := c.v.Var(id, "uuid")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	return c.s.Delete(ctx.Context(), id)
}

func (c *controller) tableList(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	err := c.v.Var(id, "uuid")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var list []*dbmodel.Table
	if list, err = c.s.TableList(ctx.Context(), id); err != nil {
		return err
	}

	return ctx.JSON(converter.ToDBTableList(list))
}

func (c *controller) functionList(ctx fiber.Ctx) error {
	id := ctx.Params("id")
	err := c.v.Var(id, "uuid")
	if err != nil {
		return fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	var list []*dbmodel.Function
	if list, err = c.s.FunctionList(ctx.Context(), id); err != nil {
		return err
	}

	return ctx.JSON(converter.ToDBFunctionList(list))
}

func (c *controller) driverList(ctx fiber.Ctx) error {
	return ctx.JSON([...]string{dbmodel.PostgreSQL})
}

func New(r fiber.Router, s Service, v *validator.Validate) {
	c := controller{s: s, v: v}
	g := r.Group("/database")
	g.Get("/", c.getList)
	g.Post("/", c.add)
	g.Patch("/:id", c.edit)
	g.Delete("/:id", c.delete)
	g.Get("/:id", c.tableList)
	g.Get("/:id", c.functionList)
	r.Get("/driver", c.driverList)
}
