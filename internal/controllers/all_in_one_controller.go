package controllers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	shippingevent "go-native-webserver/internal/service/shipping_event"

	"github.com/go-playground/validator/v10"
)

type AllInOneController interface {
	PostSubscription(w http.ResponseWriter, r *http.Request)
	GetPublicHoliday(w http.ResponseWriter, r *http.Request)
	PutShippingEvent(w http.ResponseWriter, r *http.Request)
}

type allInOneController struct {
	shippingEventService shippingevent.ShippingEventService
	validator            *validator.Validate
}

func NewAllInOneController() AllInOneController {
	return &allInOneController{
		shippingEventService: shippingevent.NewShippingEventService(),
		validator:            validator.New(),
	}
}

// Holiday represents a public holiday from the API
type Holiday struct {
	Date        string   `json:"date"`
	LocalName   string   `json:"localName"`
	Name        string   `json:"name"`
	CountryCode string   `json:"countryCode"`
	Fixed       bool     `json:"fixed"`
	Global      bool     `json:"global"`
	Counties    []string `json:"counties,omitempty"`
	LaunchYear  *int     `json:"launchYear,omitempty"`
	Types       []string `json:"types,omitempty"`
}

// QueryParams for validation
type QueryParams struct {
	Year      int      `validate:"required,min=1900,max=2100"`
	Countries []string `validate:"required,min=1,dive,required,len=2"`
}

// func (c *allInOneController) PostLogin(w http.ResponseWriter, r *http.Request) {

// }

func (c *allInOneController) PostSubscription(w http.ResponseWriter, r *http.Request) {
	err := c.shippingEventService.AddUserToShippingEventSubscription(context.Background(), 1, 1)
	if err != nil {
		ResponseError(w, err) // Handle error with fallback catching to 500 Internal Server Error
		return
	}
	ResponseSuccessJSON(w, http.StatusAccepted, map[string]string{"status": "subscribed"})
}

// ugh, there is hidden api need to make the logic
func (c *allInOneController) PutShippingEvent(w http.ResponseWriter, r *http.Request) {
	// Implement the logic for updating a shipping event.
	w.WriteHeader(http.StatusNotImplemented)
	w.Write([]byte("PutShippingEvent not implemented"))
}

// ***** This API is written by AI ***** I am hungry for dinner
func (c *allInOneController) GetPublicHoliday(w http.ResponseWriter, r *http.Request) {
	// Parse query params
	yearStr := r.URL.Query().Get("year")
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		ResponseError(w, fmt.Errorf("invalid year: %v", err))
		return
	}

	countries := r.URL.Query()["country"]
	if len(countries) == 0 {
		ResponseError(w, fmt.Errorf("at least one country required"))
		return
	}

	params := QueryParams{
		Year:      year,
		Countries: countries,
	}

	// Validate
	err = c.validator.Struct(params)
	if err != nil {
		ResponseError(w, fmt.Errorf("validation error: %v", err))
		return
	}

	// Fan out: concurrent requests
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	type result struct {
		country  string
		holidays []Holiday
		err      error
	}

	results := make(chan result, len(countries))
	var wg sync.WaitGroup

	for _, country := range countries {
		wg.Add(1)
		go func(country string) {
			defer wg.Done()
			holidays, err := c.fetchHolidays(ctx, year, country)
			results <- result{country: country, holidays: holidays, err: err}
		}(country)
	}

	// Close channel after all goroutines done
	go func() {
		wg.Wait()
		close(results)
	}()

	// Fan in: aggregate results
	aggregated := make(map[string][]Holiday)
	for res := range results {
		if res.err != nil {
			// Log error, but continue aggregating others
			fmt.Printf("Error fetching for %s: %v\n", res.country, res.err)
			continue
		}
		aggregated[res.country] = res.holidays
	}

	// Respond with 200 OK
	ResponseSuccessJSON(w, http.StatusOK, aggregated)
}

// fetchHolidays makes HTTP request to the API
func (c *allInOneController) fetchHolidays(ctx context.Context, year int, country string) ([]Holiday, error) {
	url := fmt.Sprintf("https://date.nager.at/api/v3/PublicHolidays/%d/%s", year, country)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API returned status %d", resp.StatusCode)
	}

	var holidays []Holiday
	err = json.NewDecoder(resp.Body).Decode(&holidays)
	if err != nil {
		return nil, err
	}

	return holidays, nil
}
