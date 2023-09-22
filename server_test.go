package server

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
)

func TestUserManagement(t *testing.T) {
	tempfile, err := os.CreateTemp("", "test_db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tempfile.Name())

	dbBytes, err := os.ReadFile("./db_test.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	_, err = tempfile.Write(dbBytes)

	db, dbErr := sql.Open("sqlite3", tempfile.Name())
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	defer db.Close()

	server := NewLogicAppServer(db)

	t.Run("POST /user/signup returns 200 on successful signup", func(t *testing.T) {
		b, _ := json.Marshal(Credentials{Username: "test", Password: "test"})
		credentialBuf := bytes.NewReader(b)
		request, _ := http.NewRequest(http.MethodPost, "/user/signup", credentialBuf)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusCreated, response.Code)
	})

	t.Run("POST /user/signup returns 400 if user already exists", func(t *testing.T) {
		b, _ := json.Marshal(Credentials{Username: "user1", Password: "user1"})
		credentialBuf := bytes.NewReader(b)
		request, _ := http.NewRequest(http.MethodPost, "/user/signup", credentialBuf)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})

	t.Run("POST /user/login returns 200 on successful login", func(t *testing.T) {
		b, _ := json.Marshal(Credentials{Username: "user1", Password: "user1"})
		credentialBuf := bytes.NewReader(b)
		request, _ := http.NewRequest(http.MethodPost, "/user/login", credentialBuf)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("POST /user/login returns 400 on unsuccessful login", func(t *testing.T) {
		b, _ := json.Marshal(Credentials{Username: "user1", Password: "pwd"})
		credentialBuf := bytes.NewReader(b)
		request, _ := http.NewRequest(http.MethodPost, "/user/login", credentialBuf)
		request.Header.Set("Content-Type", "application/json")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusBadRequest, response.Code)
	})
}

func TestPlayersHandler(t *testing.T) {
	tempfile, err := os.CreateTemp("", "test_db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tempfile.Name())

	dbBytes, err := os.ReadFile("./db_test.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	_, err = tempfile.Write(dbBytes)

	db, dbErr := sql.Open("sqlite3", tempfile.Name())
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	defer db.Close()

	server := NewLogicAppServer(db)

	t.Run("GET /players/ returns 200", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/players/", nil)
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("GET /players/player/achievement returns 200 if authorized", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/players/player/achievements", nil)
		request.Header.Set("Authorization", "Bearer 2TeSZoUxu3a1qCc/23KZ8PCKgeABYGuFpRgEWMw6skQ=")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("GET /players/player/achievement returns 200 if not authorized", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/players/player/achievements", nil)
		request.Header.Set("Authorization", "Bearer unknown")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusUnauthorized, response.Code)
	})

}

func TestQuizHandler(t *testing.T) {
	tempfile, err := os.CreateTemp("", "test_db.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tempfile.Name())

	dbBytes, err := os.ReadFile("./db_test.sqlite3")
	if err != nil {
		log.Fatal(err)
	}
	_, err = tempfile.Write(dbBytes)

	db, dbErr := sql.Open("sqlite3", tempfile.Name())
	if dbErr != nil {
		log.Fatal(dbErr)
	}
	defer db.Close()

	server := NewLogicAppServer(db)

	t.Run("GET /quiz/competition returns 200", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/quiz/competition", nil)
		request.Header.Set("Authorization", "Bearer 2TeSZoUxu3a1qCc/23KZ8PCKgeABYGuFpRgEWMw6skQ=")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("GET /quiz/find returns 200", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/quiz/find?type=TAUT&difficulty=easy&numVars=2&timeLimit=1000", nil)
		request.Header.Set("Authorization", "Bearer 2TeSZoUxu3a1qCc/23KZ8PCKgeABYGuFpRgEWMw6skQ=")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})

	t.Run("GET /quiz/find returns 200", func(t *testing.T) {
		request, _ := http.NewRequest(http.MethodGet, "/quiz/find?type=SAT&difficulty=easy&numVars=5&timeLimit=1000", nil)
		request.Header.Set("Authorization", "Bearer 2TeSZoUxu3a1qCc/23KZ8PCKgeABYGuFpRgEWMw6skQ=")
		response := httptest.NewRecorder()

		server.ServeHTTP(response, request)

		assert.Equal(t, http.StatusOK, response.Code)
	})
}
