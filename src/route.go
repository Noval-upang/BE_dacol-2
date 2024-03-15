package src

import (
	"net/http"

	"github.com/Hy-Iam-Noval/dacol-2/src/ctrl"
	"github.com/Hy-Iam-Noval/dacol-2/src/helpers"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/filesystem"
	_ "github.com/joho/godotenv"
)

type App = *fiber.App

func Route() App {
	r := fiber.New()

	// config

	r.Static("/public", "./public")
	r.Use(filesystem.New(filesystem.Config{
		Root: http.Dir("./public"),
	}))

	r.Use(cors.New(cors.Config{
		AllowOrigins: "http://fedacol2-production.up.railway.app",
		AllowMethods: "GET, POST, DELETE",
	}))
	r.Use(ctrl.Acceptable)

	// route
	r.Post("/login", ctrl.Login)
	r.Post("/register", ctrl.Register)

	// Group /auth
	auth := r.Group("/auth").Use(helpers.AuthWare)
	{
		auth.Get("/user", ctrl.ParseToken)

		// Post /upload/:path
		// :path path folder where ]file will be put
		auth.Post("/upload/:path", ctrl.FileUpload)

		// Group /auth/product
		auth.Post("/product_add", ctrl.AddProd)
		auth.Post("/product_update", ctrl.UpdateProd)
		auth.Get("/product_all", ctrl.AllProd)

		// Group /auth/product/:id
		prodIdR := auth.Group("/product/:id")
		{
			prodIdR.Get("/", ctrl.FindProd)
			prodIdR.Delete("/delete", ctrl.DeleteProd)
		}

		// Group /auth/selling
		auth.Post("/selling_add", ctrl.AddSelling)
		sellingR := auth.Group("/selling")
		{
			sellingR.Get("/", ctrl.AllSelling)
			sellingR.Get("/:id", ctrl.AllByIDSelling)
		}

	}
	return r
}
