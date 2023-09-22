package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/template/html/v2"
	"github.com/kokweikhong/goklse/klse"
	"github.com/kokweikhong/goklse/types"
)

func main() {
    data := types.Data{}

	engine := html.New("./views", ".html")
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/static", "./static")

	app.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Hello, World!",
		})
	})

	app.Get("/stock", func(c *fiber.Ctx) error {
        if data.StockList == nil {
            data.StockList = klse.GetStockListing()
        }
		return c.Render("stock-list-table", fiber.Map{
			"StockList": data.StockList,
		})
	})

	log.Fatal(app.Listen(":3000"))

}
