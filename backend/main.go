package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
	"strings"
)

type Response struct {
	Status     bool   `json:"status,omitempty"`
	StatusCode int    `json:"status_code,omitempty"`
	Message    string `json:"message,omitempty"`
}

func removePrefixOfBase64(base64 string) string {
	willRemovePrefix := regexp.MustCompile(`^data:.+\/.+;base64,`)
	return willRemovePrefix.ReplaceAllString(base64, "")
}

func removePrefixOfMovieRequestQuery(url string) string {
	return strings.Replace(url, "/movies/", "", -1)
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

	if r.Method == "POST" {
		var f map[string]string
		json.NewDecoder(r.Body).Decode(&f)

		saveMovie(f["movieBase64"], 1)
	}

	if r.Method == "GET" {
		url := r.URL.String()
		movieNumber := removePrefixOfMovieRequestQuery(url)
		fmt.Println(movieNumber)

		// movie_{movieNumber}.mp4 を取得する
		file, err := ioutil.ReadFile(fmt.Sprintf("www/movies/movie_%s.mp4", movieNumber))
		if err != nil {
			log.Fatal(err)
		}

		// 取得したmp4をbase64エンコード
		movieToBase64 := "data:movie/mp4;base64," + base64.StdEncoding.EncodeToString(file)

		responseBeforeMarshal := Response{
			Status:     true,
			StatusCode: 200,
			Message:    movieToBase64,
		}

		responseAfterMarshal, _ := json.Marshal(responseBeforeMarshal)
		w.Write(responseAfterMarshal)
	}
}

func main() {
	http.HandleFunc("/movies/", handler)
	http.ListenAndServe(":5555", nil)
}
