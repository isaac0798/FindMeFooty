package main

import (
	"fmt"
	"net/http"
	"io"
	"encoding/json"

	"github.com/a-h/templ"
)

func main() {
	component := hello()
	
	http.Handle("/", templ.Handler(component))
	http.HandleFunc("/result", func(w http.ResponseWriter, r *http.Request) {
		teamName := r.URL.Query().Get("team")
		results(teamName).Render(r.Context(), w)
	})

	http.HandleFunc("/fixtures", func(w http.ResponseWriter, r *http.Request) {
		//json.NewEncoder(w).Encode(p)
	})

	http.HandleFunc("/leagues", func(w http.ResponseWriter, r *http.Request) {
		leagues := getLeagueIDs()
		w.Write(leagues)
	})

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}

func getLeagueIDs() ([]byte) {
	url := "https://api-football-v1.p.rapidapi.com/v3/leagues?country=England"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "66ec3e99f0mshb514d6b2c37db99p166d41jsn622683dcbfe1")
	req.Header.Add("X-RapidAPI-Host", "api-football-v1.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	leagues, err := UnmarshalLeagues(body)
	
	if err != nil {
		fmt.Printf("%v", err)
	}

	response, _ := json.Marshal(leagues.Response)

	leagueIDs := make([]int64, 0)

	for _, a := range leagues.Response {
		fmt.Printf("%v", a)
		leagueIDs = append(leagueIDs, a.League.ID)
	}

	response, _ = json.Marshal(leagueIDs)

	return response
}

func UnmarshalLeagues(data []byte) (Leagues, error) {
	var r Leagues
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Leagues) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Leagues struct {
	Get        string        `json:"get"`
	Parameters Parameters    `json:"parameters"`
	Errors     []interface{} `json:"errors"`
	Results    int64         `json:"results"`
	Paging     Paging        `json:"paging"`
	Response   []Response    `json:"response"`
}

type Paging struct {
	Current int64 `json:"current"`
	Total   int64 `json:"total"`
}

type Parameters struct {
	Country CountryEnum `json:"country"`
}

type Response struct {
	League  LeagueClass  `json:"league"`
	Country CountryClass `json:"country"`
	Seasons []Season     `json:"seasons"`
}

type CountryClass struct {
	Name CountryEnum `json:"name"`
	Code Code        `json:"code"`
	Flag string      `json:"flag"`
}

type LeagueClass struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Type Type   `json:"type"`
	Logo string `json:"logo"`
}

type Season struct {
	Year     int64    `json:"year"`
	Start    string   `json:"start"`
	End      string   `json:"end"`
	Current  bool     `json:"current"`
	Coverage Coverage `json:"coverage"`
}

type Coverage struct {
	Fixtures    Fixtures `json:"fixtures"`
	Standings   bool     `json:"standings"`
	Players     bool     `json:"players"`
	TopScorers  bool     `json:"top_scorers"`
	TopAssists  bool     `json:"top_assists"`
	TopCards    bool     `json:"top_cards"`
	Injuries    bool     `json:"injuries"`
	Predictions bool     `json:"predictions"`
	Odds        bool     `json:"odds"`
}

type Fixtures struct {
	Events             bool `json:"events"`
	Lineups            bool `json:"lineups"`
	StatisticsFixtures bool `json:"statistics_fixtures"`
	StatisticsPlayers  bool `json:"statistics_players"`
}

type CountryEnum string

const (
	England CountryEnum = "England"
)

type Code string

const (
	GB Code = "GB"
)

type Type string

const (
	Cup    Type = "Cup"
	League Type = "League"
)