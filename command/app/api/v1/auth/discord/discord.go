package discord

import (
	"github.com/gofiber/fiber/v3"
	"github.com/z46-dev/florr-legacy/command/config"
)

func Init(parent fiber.Router) {
	var discordRouter = parent.Group("/discord")
	discordRouter.Get("/config", GetDiscordConfig)
}

func GetDiscordConfig(ctx fiber.Ctx) (err error) {
	err = ctx.JSON(fiber.Map{
		"client_id":           config.Config.Discord.ClientID,
		"client_redirect_uri": config.Config.Discord.ClientRedirect,
		"scopes":              "identify email",
	})

	return
}
