package main

import (
	"fmt"
	"net/http"
	"io"
	"encoding/json"
	"time"

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
		fixtures := getTodaysFixtures()
		w.Write(fixtures)
	})

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}

func getLeagueIDs() ([]int64) {
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

	leagueIDs := make([]int64, 0)

	for _, a := range leagues.Response {
		fmt.Printf("%v", a)
		leagueIDs = append(leagueIDs, a.League.ID)
	}

	return leagueIDs
}

func getTodaysFixtures() ([]byte) {	
	url := fmt.Sprintf("https://api-football-v1.p.rapidapi.com/v3/fixtures?date=%v",
		time.Now().Format("2006-01-02"))

	fmt.Printf(url)

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("X-RapidAPI-Key", "66ec3e99f0mshb514d6b2c37db99p166d41jsn622683dcbfe1")
	req.Header.Add("X-RapidAPI-Host", "api-football-v1.p.rapidapi.com")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	fixtures, err := UnmarshalFixtures(body)
	
	if err != nil {
		fmt.Printf("%v", err)
	}

	leagues := getLeagueIDs()
	englishFixtures := make([]Fixture, 0)

	for _, a := range fixtures.Response {
		fmt.Printf("%v", a)
		if (stringInSlice(a.League.ID, leagues)) {
			englishFixtures = append(englishFixtures, a.Fixture)
		}
	}

	response, _ := json.Marshal(englishFixtures)
	
	return response
}

func stringInSlice(a int64, list []int64) bool {
    for _, b := range list {
        if b == a {
            return true
        }
    }
    return false
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

func UnmarshalFixtures(data []byte) (FFixtures, error) {
	var r FFixtures
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *FFixtures) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type FFixtures struct {
	Get        string        `json:"get"`
	Parameters Parameters    `json:"parameters"`
	Errors     []interface{} `json:"errors"`
	Results    int64         `json:"results"`
	Paging     Paging        `json:"paging"`
	Response   []FResponse    `json:"response"`
}

type FPaging struct {
	Current int64 `json:"current"`
	Total   int64 `json:"total"`
}

type FParameters struct {
	Date string `json:"date"`
}

type FResponse struct {
	Fixture Fixture `json:"fixture"`
	League  FLeague  `json:"league"`
	Teams   Goals   `json:"teams"`
	Goals   Goals   `json:"goals"`
	Score   Score   `json:"score"`
}

type Fixture struct {
	ID        int64       `json:"id"`
	Referee   interface{} `json:"referee"`
	Timezone  Timezone    `json:"timezone"`
	Date      string      `json:"date"`
	Timestamp int64       `json:"timestamp"`
	Periods   Periods     `json:"periods"`
	Venue     Venue       `json:"venue"`
	Status    Status      `json:"status"`
}

type Periods struct {
	First  interface{} `json:"first"`
	Second interface{} `json:"second"`
}

type Status struct {
	Long    Long        `json:"long"`
	Short   Short       `json:"short"`
	Elapsed interface{} `json:"elapsed"`
}

type Venue struct {
	ID   *int64  `json:"id"`
	Name *string `json:"name"`
	City *string `json:"city"`
}

type Goals struct {
	Home *Away `json:"home"`
	Away *Away `json:"away"`
}

type Away struct {
	ID     int64       `json:"id"`
	Name   string      `json:"name"`
	Logo   string      `json:"logo"`
	Winner interface{} `json:"winner"`
}

type FLeague struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Country string `json:"country"`
	Logo    string `json:"logo"`
	Flag    string `json:"flag"`
	Season  int64  `json:"season"`
	Round   string `json:"round"`
}

type Score struct {
	Halftime  Goals `json:"halftime"`
	Fulltime  Goals `json:"fulltime"`
	Extratime Goals `json:"extratime"`
	Penalty   Goals `json:"penalty"`
}

type Long string

const (
	MatchPostponed  Long = "Match Postponed"
	NotStarted      Long = "Not Started"
	TimeToBeDefined Long = "Time to be defined"
)

type Short string

const (
	NS  Short = "NS"
	Pst Short = "PST"
	Tbd Short = "TBD"
)

type Timezone string

const (
	UTC Timezone = "UTC"
)