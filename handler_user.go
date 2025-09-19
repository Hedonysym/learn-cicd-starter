package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/bootdotdev/learn-cicd-starter/internal/database"
	"github.com/google/uuid"
)

func (cfg *apiConfig) handlerUsersCreate(w http.ResponseWriter, r *http.Request) {
	type parameters struct {
		Name string `json:"name"`
	}
	decoder := json.NewDecoder(r.Body)
	params := parameters{}
	err := decoder.Decode(&params)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't decode parameters", err)
		return
	}

	apiKey, err := generateRandomSHA256Hash()
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't gen apikey", err)
		return
	}

	err = cfg.DB.CreateUser(r.Context(), database.CreateUserParams{
		ID:        uuid.New().String(),
		CreatedAt: time.Now().UTC().Format(time.RFC3339),
		UpdatedAt: time.Now().UTC().Format(time.RFC3339),
		Name:      params.Name,
		ApiKey:    apiKey,
	})
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, fmt.Sprintf("Couldn't create user: %v", err), err)
		return
	}

	user, err := cfg.DB.GetUser(r.Context(), apiKey)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't get user", err)
		return
	}
	w.WriteHeader(201)

	userResp, err := databaseUserToUser(user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert user", err)
		return
	}
	respondWithJSON(w, 201, userResp)
}

func generateRandomSHA256Hash() (string, error) {
	randomBytes := make([]byte, 32)
	_, err := rand.Read(randomBytes)
	if err != nil {
		return "", err
	}
	hash := sha256.Sum256(randomBytes)
	hashString := hex.EncodeToString(hash[:])
	return hashString, nil
}

func (cfg *apiConfig) handlerUsersGet(w http.ResponseWriter, r *http.Request, user database.User) {

	userResp, err := databaseUserToUser(user)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, "Couldn't convert user", err)
		return
	}

	respondWithJSON(w, http.StatusOK, userResp)
}
