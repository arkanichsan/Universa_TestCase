package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	_ "fmt"
)

const (
	DB_USER     = "postgres"
	DB_PASSWORD = "gitong1999"
	DB_NAME     = "music"
)

func setupDB() *sql.DB {
	dbinfo := fmt.Sprintf("user=%s password=%s dbname=%s sslmode=disable", DB_USER, DB_PASSWORD, DB_NAME)
	db, err := sql.Open("postgres", dbinfo)

	checkErr(err)

	return db
}

type Music struct {
	MusicID     string    `json:"musicid"`
	MusicName   string    `json:"musicname"`
	MusicAlbum  string    `json:"musicalbum"`
	MusicArt    string    `json:"musicart"`
	MusicSinger string    `json:"musicinger"`
	Musicdate   time.Time `json:"musicdate"`
}

type JsonResponse struct {
	Type    string  `json:"type"`
	Data    []Music `json:"data"`
	Message string  `json:"message"`
}

func main() {

	router := mux.NewRouter()

	router.HandleFunc("/music/", Getmusic).Methods("GET")

	router.HandleFunc("/music/", Createmusic).Methods("POST")

	router.HandleFunc("/music/{musicid}", Deletemusic).Methods("DELETE")

	router.HandleFunc("/music/", Deletemusic).Methods("DELETE")

	// serve the app
	fmt.Println("Server at 8080")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func printMessage(message string) {
	fmt.Println("")
	fmt.Println(message)
	fmt.Println("")
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}

func Getmusic(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Getting music...")

	rows, err := db.Query("SELECT * FROM music")

	checkErr(err)

	var music []Music

	for rows.Next() {
		var id int
		var musicID string
		var musicName string

		err = rows.Scan(&id, &musicID, &musicName)

		checkErr(err)

		music = append(music, Music{MusicID: musicID, MusicName: musicName})
	}

	var response = JsonResponse{Type: "success", Data: music}

	json.NewEncoder(w).Encode(response)
}

func Createmusic(w http.ResponseWriter, r *http.Request) {
	musicID := r.FormValue("musicid")
	musicName := r.FormValue("musicname")

	var response = JsonResponse{}

	if musicID == "" || musicName == "" {
		response = JsonResponse{Type: "error", Message: "You are missing musicID or musicName parameter."}
	} else {
		db := setupDB()

		printMessage("Inserting music into DB")

		fmt.Println("Inserting new music with ID: " + musicID + " and name: " + musicName)

		var lastInsertID int
		err := db.QueryRow("INSERT INTO music(musicID, musicName) VALUES($1, $2) returning id;", musicID, musicName).Scan(&lastInsertID)

		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The music has been inserted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

func deletemusic(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)

	musicID := params["musicid"]

	var response = JsonResponse{}

	if musicID == "" {
		response = JsonResponse{Type: "error", Message: "You are missing musicID parameter."}
	} else {
		db := setupDB()

		printMessage("Deleting music from DB")

		_, err := db.Exec("DELETE FROM music where musicID = $1", musicID)

		checkErr(err)

		response = JsonResponse{Type: "success", Message: "The music has been deleted successfully!"}
	}

	json.NewEncoder(w).Encode(response)
}

func Deletemusic(w http.ResponseWriter, r *http.Request) {
	db := setupDB()

	printMessage("Deleting all music...")

	_, err := db.Exec("DELETE FROM music")

	checkErr(err)

	printMessage("All music have been deleted successfully!")

	var response = JsonResponse{Type: "success", Message: "All music have been deleted successfully!"}

	json.NewEncoder(w).Encode(response)
}
