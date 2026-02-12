package apibalancecharge

import (
	"strconv"

	"github.com/barkhayot/backend-test-golang-temp/pkg/appctx"
	"github.com/barkhayot/backend-test-golang-temp/pkg/balance"
	"github.com/gofiber/fiber/v3"
)

type ErrorResponse struct {
	Error string `json:"error"`
}

func Handler(appCtx *appctx.AppCtx) fiber.Handler {
	return func(c fiber.Ctx) error {
		idStr := c.Params("id")
		id, err := strconv.ParseInt(idStr, 10, 64)
		if err != nil {
			return fiber.ErrBadRequest
		}

		var req struct {
			Amount float64 `json:"amount"`
		}

		if err := c.Bind().Body(&req); err != nil {
			return fiber.ErrBadRequest
		}

		h, err := appCtx.GetBalance().Charge(c.Context(), id, req.Amount)
		if err == balance.ErrInsufficientBalance {
			return c.Status(fiber.StatusBadRequest).JSON(ErrorResponse{
				Error: err.Error(),
			})
		}
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(ErrorResponse{
				Error: "internal server error",
			})
		}

		return c.JSON(h)
	}
}
