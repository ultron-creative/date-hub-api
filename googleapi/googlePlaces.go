package googleapi

import (
	"date-hub-api/server"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
)

func getFood(w http.ResponseWriter, r *http.Request) {

	type request struct {
		Lat     string
		Lon     string
		Type    string
		Keyword string
	}

	type response struct {
		HTMLAttributions []interface{} `json:"html_attributions"`
		NextPageToken    string        `json:"next_page_token"`
		Results          []struct {
			BusinessStatus string `json:"business_status"`
			Geometry       struct {
				Location struct {
					Lat float64 `json:"lat"`
					Lng float64 `json:"lng"`
				} `json:"location"`
				Viewport struct {
					Northeast struct {
						Lat float64 `json:"lat"`
						Lng float64 `json:"lng"`
					} `json:"northeast"`
					Southwest struct {
						Lat float64 `json:"lat"`
						Lng float64 `json:"lng"`
					} `json:"southwest"`
				} `json:"viewport"`
			} `json:"geometry"`
			Icon         string `json:"icon"`
			Name         string `json:"name"`
			OpeningHours struct {
				OpenNow bool `json:"open_now"`
			} `json:"opening_hours"`
			Photos []struct {
				Height           int      `json:"height"`
				HTMLAttributions []string `json:"html_attributions"`
				PhotoReference   string   `json:"photo_reference"`
				Width            int      `json:"width"`
			} `json:"photos"`
			PlaceID  string `json:"place_id"`
			PlusCode struct {
				CompoundCode string `json:"compound_code"`
				GlobalCode   string `json:"global_code"`
			} `json:"plus_code"`
			PriceLevel       int      `json:"price_level,omitempty"`
			Rating           float64  `json:"rating"`
			Reference        string   `json:"reference"`
			Scope            string   `json:"scope"`
			Types            []string `json:"types"`
			UserRatingsTotal int      `json:"user_ratings_total"`
			Vicinity         string   `json:"vicinity"`
		} `json:"results"`
		Status string `json:"status"`
	}

	handler := func(req request) *response {
		URL := fmt.Sprintf(`https://maps.googleapis.com/maps/api/place/nearbysearch/json?location=%s,%s&radius=8000&type=%s&keyword=%s&key=%s`, req.Lat, req.Lon, req.Type, req.Keyword, os.Getenv("GOOGLE_KEY"))

		r, err := http.NewRequest("GET", URL, nil)
		if err != nil {
			server.PanicWithStatus(err, http.StatusBadRequest)
		}

		client := &http.Client{}
		resp, err := client.Do(r)
		if err != nil {
			server.PanicWithStatus(err, http.StatusBadRequest)
		}
		defer resp.Body.Close()

		var res response
		if err = json.NewDecoder(resp.Body).Decode(&res); err != nil {
			log.Println(err)
		}

		if resp.StatusCode > 204 {
			panic(fmt.Errorf("expected status code 200 or 204 but got %d", resp.StatusCode))
		}

		return &res
	}

	var req request = request{Lat: server.GetStringParam(r, "lat"), Lon: server.GetStringParam(r, "lon"), Type: server.GetStringParam(r, "type"), Keyword: server.GetStringParam(r, "keyword", true)}

	res := handler(req)

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(res); err != nil {
		panic(err)
	}
}

// amusement_park
// aquarium
// art_gallery
// beauty_salon
// bowling_alley
// campground
// casino
// gym
// movie_theater
// museum
// night_club
// park
// spa
// zoo
