package controllers

import (
	"context"
	"math"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/flaxrora/catchphrase-go-mongodb-rest-api/config"
	"github.com/flaxrora/catchphrase-go-mongodb-rest-api/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetAllRatings(c *fiber.Ctx) error {
	ratingsCollection := config.MI.DB.Collection("ratings")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	var ratings []models.Ratings

	filter := bson.M{}
	findOptions := options.Find()

	if s := c.Query("s"); s != "" {
		// filter = bson.M{
		// 	"$or": []bson.M{
		// 		{
		// 			"movieName": bson.M{
		// 				"$regex": primitive.Regex{
		// 					Pattern: s,
		// 					Options: "i",
		// 				},
		// 			},
		// 		},
		// 		{
		// 			"rating": bson.M{
		// 				"$regex": primitive.Regex{
		// 					Pattern: s,
		// 					Options: "i",
		// 				},
		// 			},
		// 		},
		// 	},
		// }
	}

	page, _ := strconv.Atoi(c.Query("page", "1"))
	limitVal, _ := strconv.Atoi(c.Query("limit", "100"))
	var limit int64 = int64(limitVal)

	total, _ := ratingsCollection.CountDocuments(ctx, filter)

	findOptions.SetSkip((int64(page) - 1) * limit)
	findOptions.SetLimit(limit)

	cursor, err := ratingsCollection.Find(ctx, filter, findOptions)
	defer cursor.Close(ctx)

	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Ratings Not found",
			"error":   err,
		})
	}

	for cursor.Next(ctx) {
		var rating models.Ratings
		cursor.Decode(&rating)
		ratings = append(ratings, rating)
	}

	last := math.Ceil(float64(total / limit))
	if last < 1 && total > 0 {
		last = 1
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":      ratings,
		"total":     total,
		"page":      page,
		"last_page": last,
		"limit":     limit,
	})
}

func GetRating(c *fiber.Ctx) error {
	ratingsCollection := config.MI.DB.Collection("ratings")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	var rating models.Ratings
	objId, err := primitive.ObjectIDFromHex(c.Params("id"))
	findResult := ratingsCollection.FindOne(ctx, bson.M{"_id": objId})
	if err := findResult.Err(); err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Rating Not found",
			"error":   err,
		})
	}

	err = findResult.Decode(&rating)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"success": false,
			"message": "Rating Not found",
			"error":   err,
		})
	}

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"data":    rating,
		"success": true,
	})
}

