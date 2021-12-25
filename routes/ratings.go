package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/flaxrora/catchphrase-go-mongodb-rest-api/controllers" // replace
)

func RatingsRoute(route fiber.Router) {
	route.Get("/", controllers.GetAllRatings)
	route.Get("/:id", controllers.GetRating)
}
