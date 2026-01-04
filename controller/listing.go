package controller

import (
	"io"

	"github.com/gofiber/fiber/v2"
	"github.com/rs/zerolog"
	"github.com/shaileshhb/websockets/db/model"
	"github.com/shaileshhb/websockets/service"
)

type ListingController struct {
	listingService *service.ListingService
	log            zerolog.Logger
}

func NewListingController(listingService *service.ListingService, log zerolog.Logger) *ListingController {
	return &ListingController{listingService: listingService}
}

// RegisterRoute registers all endpoints to router.
func (controller *ListingController) RegisterRoute(router fiber.Router) {
	router.Get("/exercises", controller.GetExerciseListings)
	router.Get("/search", controller.GetExerciseSearch)
	router.Get("/exercises/:id", controller.GetExerciseById)
	router.Get("/videos/*", controller.ProxyVideo)
	controller.log.Info().Msg("Quiz routes registered")
}

// GetExerciseListings will get the exercise listings from the MuscleWiki API
func (lc *ListingController) GetExerciseListings(c *fiber.Ctx) error {
	exerciseResponse := &model.ExerciseResponse{}
	limit := c.Query("limit")
	offset := c.Query("offset")
	search := c.Query("search")
	err := lc.listingService.GetExerciseListings(exerciseResponse, limit, offset, search)
	if err != nil {
		lc.log.Error().Err(err).Msg("Failed to get exercise listings")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get exercise listings",
		})
	}
	return c.Status(fiber.StatusOK).JSON(exerciseResponse)
}

// GetExerciseSearch will get the exercise search from the MuscleWiki API
func (lc *ListingController) GetExerciseSearch(c *fiber.Ctx) error {
	exerciseSearch := &[]model.ExerciseList{}
	limit := c.Query("limit")
	search := c.Query("search")
	err := lc.listingService.GetExerciseSearch(exerciseSearch, limit, search)
	if err != nil {
		lc.log.Error().Err(err).Msg("Failed to get exercise search")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get exercise search",
		})
	}
	return c.Status(fiber.StatusOK).JSON(exerciseSearch)
}

// GetExerciseById will get the exercise detail from the MuscleWiki API
func (lc *ListingController) GetExerciseById(c *fiber.Ctx) error {
	exerciseDetailID := c.Params("id")
	exerciseDetail := &model.ExerciseDetail{}
	err := lc.listingService.GetExerciseById(exerciseDetailID, exerciseDetail)
	if err != nil {
		lc.log.Error().Err(err).Msg("Failed to get exercise detail")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to get exercise detail",
		})
	}
	return c.Status(fiber.StatusOK).JSON(exerciseDetail)
}

// ProxyVideo will proxy the video request to the MuscleWiki API
func (lc *ListingController) ProxyVideo(c *fiber.Ctx) error {
	// Get the video path (everything after /videos/)
	videoPath := c.Params("*")

	if videoPath == "" {
		lc.log.Error().Msg("Video path is required")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Video path is required",
		})
	}

	lc.log.Info().Str("path", videoPath).Msg("Proxying video request")

	// Get Range header from request for video seeking
	rangeHeader := c.Get("Range")

	// Proxy the video request through the service
	resp, err := lc.listingService.ProxyVideo(videoPath, rangeHeader)
	if err != nil {
		lc.log.Error().Err(err).Msg("Failed to proxy video")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to fetch video",
		})
	}
	defer resp.Body.Close()

	// Check if the upstream request was successful
	if resp.StatusCode >= 400 {
		lc.log.Error().Int("status", resp.StatusCode).Msg("Upstream API returned error")
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"message": "Failed to fetch video from upstream",
		})
	}

	// Forward the response status
	if resp.StatusCode >= 400 {
		lc.log.Error().Int("status", resp.StatusCode).Msg("Upstream API returned error")
		return c.Status(resp.StatusCode).JSON(fiber.Map{
			"message": "Failed to fetch video from upstream",
		})
	}

	// Forward the response status
	c.Status(resp.StatusCode)

	// Forward essential headers from the proxied response
	if contentType := resp.Header.Get("Content-Type"); contentType != "" {
		c.Set("Content-Type", contentType)
	}
	if contentLength := resp.Header.Get("Content-Length"); contentLength != "" {
		c.Set("Content-Length", contentLength)
	}
	if contentRange := resp.Header.Get("Content-Range"); contentRange != "" {
		c.Set("Content-Range", contentRange)
	}
	if acceptRanges := resp.Header.Get("Accept-Ranges"); acceptRanges != "" {
		c.Set("Accept-Ranges", acceptRanges)
	}

	// Set cache control header
	c.Set("Cache-Control", "public, max-age=3600")

	// Read the entire response body (like curl does)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		lc.log.Error().Err(err).Msg("Failed to read video response")
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Failed to read video",
		})
	}

	lc.log.Info().Int("bytes", len(body)).Msg("Successfully read video")

	// Send the video content
	return c.Send(body)
}
