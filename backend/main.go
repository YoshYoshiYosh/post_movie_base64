package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"regexp"
)

type Response struct {
	Status     bool
	StatusCode int
	Message    string
}

func removePrefixOfBase64(base64 string) string {
	willRemovePrefix := regexp.MustCompile(`^data:.+\/.+;base64,`)
	return willRemovePrefix.ReplaceAllString(base64, "")
}

func saveMovie(movieBase64 string, itemId int) {
	base64Formatted := removePrefixOfBase64(movieBase64)

	decodedMovie, err := base64.StdEncoding.DecodeString(base64Formatted)

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(decodedMovie)

	newMovieFile, _ := os.Create(fmt.Sprintf("www/movies/movie_%d.mp4", itemId))
	newMovieFile.Write(decodedMovie)
}

func handler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Headers", "*")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
	w.Header().Set("Content-Type", "application/json")

	if r.Method == "OPTIONS" {
		responseBeforeMarshal := Response{
			Status:     true,
			StatusCode: 200,
			Message:    "success!",
		}

		responseAfterMarshal, _ := json.Marshal(responseBeforeMarshal)
		w.Write(responseAfterMarshal)
		return
	}

	var f map[string]string
	json.NewDecoder(r.Body).Decode(&f)

	saveMovie(f["movieBase64"], 1)
}

func main() {
	http.HandleFunc("/movies", handler)
	http.ListenAndServe(":5555", nil)
}
