package service

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/shaileshhb/websockets/db"
	"github.com/shaileshhb/websockets/db/model"
)

// ListingService will contain reference to db.
type ListingService struct {
	db *db.Database
}

// NewListingService will create new instance of ListingService
func NewListingService(db *db.Database) *ListingService {
	return &ListingService{db: db}
}

// GetExerciseListings will get the exercise listings from the MuscleWiki API
func (ls *ListingService) GetExerciseListings(exerciseResponse *model.ExerciseResponse, limit string, offset string, search string) error {
	url := "https://musclewiki-api.p.rapidapi.com/exercises"
	queryparams := ""

	if limit == "" {
		queryparams += fmt.Sprintf("limit=%d&", 20)
	} else {
		queryparams += fmt.Sprintf("limit=%s&", limit)
	}

	if offset == "" {
		queryparams += fmt.Sprintf("offset=%d&", 0)
	} else {
		queryparams += fmt.Sprintf("offset=%s&", offset)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", url, queryparams), nil)
	if err != nil {
		return err
	}

	req.Header.Add("x-rapidapi-key", os.Getenv("RAPID_API_KEY"))
	req.Header.Add("x-rapidapi-host", os.Getenv("RAPID_API_HOST"))

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(exerciseResponse)
	if err != nil {
		return err
	}

	return nil
}

// GetExerciseById will get the exercise detail from the MuscleWiki API
func (ls *ListingService) GetExerciseById(exerciseDetailID string, exerciseDetail *model.ExerciseDetail) error {
	url := fmt.Sprintf("https://musclewiki-api.p.rapidapi.com/exercises/%s?detail=true", exerciseDetailID)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	req.Header.Add("x-rapidapi-key", os.Getenv("RAPID_API_KEY"))
	req.Header.Add("x-rapidapi-host", os.Getenv("RAPID_API_HOST"))

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(exerciseDetail)
	if err != nil {
		return err
	}

	return nil
}

// ProxyVideo proxies video requests to MuscleWiki API with authentication
func (ls *ListingService) ProxyVideo(videoPath string, rangeHeader string) (*http.Response, error) {
	url := fmt.Sprintf("https://musclewiki-api.p.rapidapi.com/media/videos/%s", videoPath)

	client := &http.Client{
		Timeout: 60 * time.Second, // Longer timeout for video streaming
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Add RapidAPI authentication headers
	req.Header.Add("x-rapidapi-key", os.Getenv("RAPID_API_KEY"))
	req.Header.Add("x-rapidapi-host", os.Getenv("RAPID_API_HOST"))

	// Forward Range header for video seeking/streaming
	if rangeHeader != "" {
		req.Header.Add("Range", rangeHeader)
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch video: %w", err)
	}

	return res, nil
}

func (ls *ListingService) GetExerciseSearch(exerciseSearch *[]model.ExerciseList, limit, search string) error {
	url := "https://musclewiki-api.p.rapidapi.com/search"
	queryparams := ""

	if limit == "" {
		queryparams += fmt.Sprintf("limit=%d&", 20)
	} else {
		queryparams += fmt.Sprintf("limit=%s&", limit)
	}

	if search != "" && len(search) >= 3 {
		queryparams += fmt.Sprintf("q=%s", search)
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s?%s", url, queryparams), nil)
	if err != nil {
		return err
	}

	req.Header.Add("x-rapidapi-key", os.Getenv("RAPID_API_KEY"))
	req.Header.Add("x-rapidapi-host", os.Getenv("RAPID_API_HOST"))

	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()

	err = json.NewDecoder(res.Body).Decode(exerciseSearch)
	if err != nil {
		return err
	}

	return nil
}
