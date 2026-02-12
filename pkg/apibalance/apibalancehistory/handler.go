package apibalancehistory

import (
	"strconv"

	"github.com/barkhayot/backend-test-golang-temp/pkg/appctx"
	"github.com/gofiber/fiber/v3"
)

func Handler(appCtx *appctx.AppCtx) fiber.Handler {
	return func(c fiber.Ctx) error {
		idStr := c.Params("id")
		userID, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return fiber.ErrBadRequest
		}

		balance := appCtx.GetBalance()
		history, err := balance.GetHistory(c.Context(), userID)
		if err != nil {
			return err
		}

		return c.JSON(history)
	}
}
