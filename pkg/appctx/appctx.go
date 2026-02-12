package appctx

import (
	"github.com/barkhayot/backend-test-golang-temp/pkg/balance"
	"github.com/barkhayot/backend-test-golang-temp/pkg/skinport"
	"github.com/gofiber/fiber/v3"
)

type AppCtx struct {
	fiberApp        *fiber.App
	skinportService *skinport.Service
	balance         *balance.Service
}

type Option func(*AppCtx)

func WithFiberApp(app *fiber.App) Option {
	return func(ac *AppCtx) {
		ac.fiberApp = app
	}
}

func WithSkinportService(svc *skinport.Service) Option {
	return func(ac *AppCtx) {
		ac.skinportService = svc
	}
}

func WithBalanceService(svc *balance.Service) Option {
	return func(ac *AppCtx) {
		ac.balance = svc
	}
}

func NewAppCtx(opts ...Option) *AppCtx {
	ac := &AppCtx{}

	for _, opt := range opts {
		opt(ac)
	}
	return ac
}

func (ac *AppCtx) GetFiberApp() *fiber.App {
	return ac.fiberApp
}

func (ac *AppCtx) GetSkinport() *skinport.Service {
	return ac.skinportService
}

func (ac *AppCtx) GetBalance() *balance.Service {
	return ac.balance
}
