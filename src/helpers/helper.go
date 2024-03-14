package helpers

import (
	_ "github.com/go-sql-driver/mysql"

	"github.com/gofiber/fiber/v2"
)

type (
	Ctx    = *fiber.Ctx
	ImgReq struct {
		Name string `json:"name"`
		Size uint   `json:"size"`
	}
)

const (
	Invalid = fiber.StatusForbidden
	UnAuth  = fiber.StatusUnauthorized
)
