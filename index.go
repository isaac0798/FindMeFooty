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
		url := "https://api-football-v1.p.rapidapi.com/v3/leagues?country=England"

		req, _ := http.NewRequest("GET", url, nil)

		req.Header.Add("X-RapidAPI-Key", "66ec3e99f0mshb514d6b2c37db99p166d41jsn622683dcbfe1")
		req.Header.Add("X-RapidAPI-Host", "api-football-v1.p.rapidapi.com")

		res, _ := http.DefaultClient.Do(req)

		defer res.Body.Close()
		body, _ := io.ReadAll(res.Body)

		var data map[string][]byte 
		err := json.Unmarshal(body, &data)
    if err != nil {
			fmt.Printf("%v", err)
    }

		w.Write(data["response"])
	})

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}
