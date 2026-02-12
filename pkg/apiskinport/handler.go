package apiskinport

import (
	"github.com/barkhayot/backend-test-golang-temp/pkg/appctx"
	"github.com/gofiber/fiber/v3"
)

func Handler(appCtx *appctx.AppCtx) fiber.Handler {
	return func(c fiber.Ctx) error {
		return c.JSON(appCtx.GetSkinport().GetItems())
	}
}
