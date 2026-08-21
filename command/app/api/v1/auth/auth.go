package auth

import (
	"github.com/gofiber/fiber/v3"
	"github.com/z46-dev/florr-legacy/command/app/api/v1/auth/discord"
)

func Init(parent fiber.Router) {
	var authRouter = parent.Group("/auth")
	discord.Init(authRouter)
}
