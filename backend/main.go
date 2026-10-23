package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

// ─── Models ──────────────────────────────────────────────────────────────────

type Entry struct {
	Date   string  `json:"date"`
	Amount float64 `json:"amount"`
	Note   string  `json:"note"`
}

type GoldEntry struct {
	ID      int     `json:"id"`
	Date    string  `json:"date"`
	Unit    string  `json:"unit"`
	Qty     float64 `json:"qty"`
	Grams   float64 `json:"grams"`
	Money   float64 `json:"money"`
	Note    string  `json:"note"`
}

// Open-Meteo response
type OpenMeteoResp struct {
	Current struct {
		Temperature float64 `json:"temperature_2m"`
		WeatherCode int     `json:"weather_code"`
		RelHumidity int     `json:"relative_humidity_2m"`
		WindSpeed   float64 `json:"wind_speed_10m"`
	} `json:"current"`
}

type WeatherResponse struct {
	Temp        float64 `json:"temp"`
	Description string  `json:"description"`
	Humidity    int     `json:"humidity"`
	WindSpeed   float64 `json:"wind_speed"`
	Icon        string  `json:"icon"`
}

var db *sql.DB

var unitToGrams = map[string]float64{
	"g":        1.0,
	"chi":      3.75,
	"domlerng": 37.5,
	"hun":      0.375,
	"baht":     15.244,
	"kg":       1000.0,
}

// ─── Main ─────────────────────────────────────────────────────────────────────

func main() {
	if err := os.MkdirAll("../data", 0755); err != nil {
		log.Fatalf("failed to create data dir: %v", err)
	}

	var err error
	db, err = sql.Open("sqlite3", "../data/savings.db")
	if err != nil {
		log.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS entries (
		date   TEXT PRIMARY KEY,
		amount REAL NOT NULL DEFAULT 0,
		note   TEXT DEFAULT ''
	)`)
	if err != nil {
		log.Fatalf("failed to create entries table: %v", err)
	}
	// Migration in case column doesn't exist
	db.Exec("ALTER TABLE entries ADD COLUMN note TEXT DEFAULT '';")

	_, err = db.Exec(`CREATE TABLE IF NOT EXISTS gold_entries (
		id      INTEGER PRIMARY KEY AUTOINCREMENT,
		date    TEXT NOT NULL,
		unit    TEXT NOT NULL,
		qty     REAL NOT NULL,
		grams   REAL NOT NULL,
		money   REAL NOT NULL DEFAULT 0,
		note    TEXT DEFAULT ''
	)`)
	if err != nil {
		log.Fatalf("failed to create gold_entries table: %v", err)
	}

	mux := http.NewServeMux()

	// Money saving routes
	mux.HandleFunc("/api/entries", withCORS(handleEntries))
	mux.HandleFunc("/api/entry", withCORS(handleEntry))
	mux.HandleFunc("/api/entry/", withCORS(handleDelete))
	mux.HandleFunc("/api/total", withCORS(handleTotal))

	// Gold routes
	mux.HandleFunc("/api/gold", withCORS(handleGold))
	mux.HandleFunc("/api/gold/", withCORS(handleGoldDelete))

	// Weather
	mux.HandleFunc("/api/weather", withCORS(handleWeather))

	// Static files
	mux.Handle("/", http.FileServer(http.Dir("../public")))

	port := "3000"
	log.Printf("✅  Savings Tracker running at http://localhost:%s", port)
	if err := http.ListenAndServe(":"+port, mux); err != nil {
		log.Fatalf("server error: %v", err)
	}
}

// ─── Middleware ───────────────────────────────────────────────────────────────

func withCORS(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}
		h(w, r)
	}
}

// ─── Money Handlers ────────────────────────────────────────────────────────────

func handleEntries(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	rows, err := db.Query("SELECT date, amount, IFNULL(note,'') FROM entries ORDER BY date DESC")
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}
	defer rows.Close()
	var entries []Entry
	for rows.Next() {
		var e Entry
		if err := rows.Scan(&e.Date, &e.Amount, &e.Note); err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
			return
		}
		entries = append(entries, e)
	}
	if entries == nil {
		entries = []Entry{}
	}
	json.NewEncoder(w).Encode(entries)
}

func handleEntry(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var e Entry
	if err := json.NewDecoder(r.Body).Decode(&e); err != nil || e.Date == "" {
		http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
		return
	}
	_, err := db.Exec("INSERT OR REPLACE INTO entries(date, amount, note) VALUES(?,?,?)", e.Date, e.Amount, e.Note)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "ok"})
}

func handleDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 || parts[3] == "" {
		http.Error(w, `{"error":"missing date"}`, http.StatusBadRequest)
		return
	}
	date := parts[3]
	_, err := db.Exec("DELETE FROM entries WHERE date = ?", date)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

func handleTotal(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	var total float64
	if err := db.QueryRow("SELECT IFNULL(SUM(amount),0) FROM entries").Scan(&total); err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]float64{"total": total})
}

// ─── Gold Handlers ─────────────────────────────────────────────────────────────

func handleGold(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rows, err := db.Query("SELECT id, date, unit, qty, grams, money, IFNULL(note,'') FROM gold_entries ORDER BY date DESC, id DESC")
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
			return
		}
		defer rows.Close()
		var entries []GoldEntry
		for rows.Next() {
			var e GoldEntry
			if err := rows.Scan(&e.ID, &e.Date, &e.Unit, &e.Qty, &e.Grams, &e.Money, &e.Note); err != nil {
				http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
				return
			}
			entries = append(entries, e)
		}
		if entries == nil {
			entries = []GoldEntry{}
		}
		json.NewEncoder(w).Encode(entries)

	case http.MethodPost:
		var e GoldEntry
		if err := json.NewDecoder(r.Body).Decode(&e); err != nil || e.Date == "" || e.Unit == "" {
			http.Error(w, `{"error":"bad request"}`, http.StatusBadRequest)
			return
		}
		// Calculate grams
		ratio, ok := unitToGrams[strings.ToLower(e.Unit)]
		if !ok {
			ratio = 1.0
		}
		e.Grams = e.Qty * ratio

		_, err := db.Exec(
			"INSERT INTO gold_entries(date, unit, qty, grams, money, note) VALUES(?,?,?,?,?,?)",
			e.Date, e.Unit, e.Qty, e.Grams, e.Money, e.Note,
		)
		if err != nil {
			http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode(map[string]string{"status": "ok"})

	default:
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
	}
}

func handleGoldDelete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	parts := strings.Split(r.URL.Path, "/")
	if len(parts) < 4 || parts[3] == "" {
		http.Error(w, `{"error":"missing id"}`, http.StatusBadRequest)
		return
	}
	id := parts[3]
	_, err := db.Exec("DELETE FROM gold_entries WHERE id = ?", id)
	if err != nil {
		http.Error(w, fmt.Sprintf(`{"error":"%v"}`, err), http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "deleted"})
}

// ─── Weather Handler ──────────────────────────────────────────────────────────

func handleWeather(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, `{"error":"method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}
	url := "https://api.open-meteo.com/v1/forecast?latitude=11.5625&longitude=104.9160" +
		"&current=temperature_2m,relative_humidity_2m,weather_code,wind_speed_10m" +
		"&timezone=Asia%2FPhnom_Penh"
	resp, err := http.Get(url)
	if err != nil {
		http.Error(w, `{"error":"weather fetch failed"}`, http.StatusInternalServerError)
		return
	}
	defer resp.Body.Close()
	var raw OpenMeteoResp
	if err := json.NewDecoder(resp.Body).Decode(&raw); err != nil {
		http.Error(w, `{"error":"decode failed"}`, http.StatusInternalServerError)
		return
	}
	result := WeatherResponse{
		Temp:        raw.Current.Temperature,
		Description: wmoDescription(raw.Current.WeatherCode),
		Humidity:    raw.Current.RelHumidity,
		WindSpeed:   raw.Current.WindSpeed,
		Icon:        wmoIcon(raw.Current.WeatherCode),
	}
	json.NewEncoder(w).Encode(result)
}

// ─── WMO Helpers ──────────────────────────────────────────────────────────────

func wmoDescription(code int) string {
	switch {
	case code == 0:
		return "Clear sky"
	case code <= 2:
		return "Partly cloudy"
	case code == 3:
		return "Overcast"
	case code <= 49:
		return "Fog"
	case code <= 59:
		return "Drizzle"
	case code <= 69:
		return "Rain"
	case code <= 79:
		return "Snow"
	case code <= 82:
		return "Rain showers"
	case code <= 84:
		return "Snow showers"
	case code <= 99:
		return "Thunderstorm"
	default:
		return "Unknown"
	}
}

func wmoIcon(code int) string {
	switch {
	case code == 0:
		return "☀️"
	case code <= 2:
		return "⛅"
	case code == 3:
		return "☁️"
	case code <= 49:
		return "🌫️"
	case code <= 69:
		return "🌧️"
	case code <= 79:
		return "❄️"
	case code <= 84:
		return "🌦️"
	case code <= 99:
		return "⛈️"
	default:
		return "🌡️"
	}
}
