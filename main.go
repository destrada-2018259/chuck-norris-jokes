package main

import (
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

type ChuckNorrisJoke struct {
	Id        string `json:"id"`
	Url       string `json:"url"`
	Value     string `json:"value"`
	IconUrl   string `json:"icon_url"`
	CreatedAt string `json:"created_at"`
}

func main() {

	http.HandleFunc("/api/jokes", func(w http.ResponseWriter, r *http.Request) {

		if r.Method != http.MethodGet {
			http.Error(w, "method not supported", http.StatusNotImplemented)
			return
		}

		jokes, err := getJokesList(25)

		if err != nil {
			http.Error(w, "Error requesting joke", http.StatusInternalServerError)
			return
		}

		jokeJson, err := json.Marshal(jokes)

		if err != nil {
			http.Error(w, "Error turning jokes into json", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.Write(jokeJson)

	})

	server := http.Server{
		Addr: ":8080",
	}

	err := server.ListenAndServe()

	if err != nil {
		panic(err)
	}
}

var existingIds = make(map[string]bool)

func getJokesList(count int) ([]ChuckNorrisJoke, error) {

	var wg sync.WaitGroup
	var mu sync.Mutex

	var list []ChuckNorrisJoke
	var jokesChannel = make(chan ChuckNorrisJoke, count)

	wg.Add(count)
	for i := 0; i < count; i++ {
		go func() {
			defer wg.Done()
			for {
				joke, err := getChuckNorrisJoke()
				if err != nil {
					return
				}

				mu.Lock()
				if !existingIds[joke.Id] {
					existingIds[joke.Id] = true
					mu.Unlock()

					jokesChannel <- joke
					return
				}
				mu.Unlock()
			}
		}()
	}

	go func() {
		wg.Wait()
		close(jokesChannel)
	}()

	for joke := range jokesChannel {
		list = append(list, joke)
	}

	return list, nil
}

func getChuckNorrisJoke() (ChuckNorrisJoke, error) {
	resp, err := http.Get("https://api.chucknorris.io/jokes/random")
	if err != nil {
		return ChuckNorrisJoke{}, err
	}
	defer resp.Body.Close()

	var joke ChuckNorrisJoke
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return ChuckNorrisJoke{}, err
	}

	err = json.Unmarshal(body, &joke)
	if err != nil {
		return ChuckNorrisJoke{}, err
	}

	return joke, nil
}
