package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/matvoy/chat_server/chat_api/repo"
	"github.com/rs/zerolog"
)

type ApiHandlers interface {
	GetProfiles(w http.ResponseWriter, r *http.Request)
	GetConversations(w http.ResponseWriter, r *http.Request)
	GetMessages(w http.ResponseWriter, r *http.Request)
	GetClients(w http.ResponseWriter, r *http.Request)
	GetUserConversations(w http.ResponseWriter, r *http.Request)
	GetAttachments(w http.ResponseWriter, r *http.Request)
}

type apiHandlers struct {
	repo repo.Repository
	log  *zerolog.Logger
}

func NewApiHandlers(repo repo.Repository, log *zerolog.Logger) ApiHandlers {
	return &apiHandlers{
		repo,
		log,
	}
}

func (a *apiHandlers) GetProfiles(w http.ResponseWriter, r *http.Request) {
	limit, page := GetDefaultQueryParams(r)
	result, err := a.repo.GetProfiles(r.Context(), limit, (page-1)*limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	body, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(body)
}

func (a *apiHandlers) GetConversations(w http.ResponseWriter, r *http.Request) {
	limit, page := GetDefaultQueryParams(r)
	result, err := a.repo.GetConversations(r.Context(), limit, (page-1)*limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	body, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(body)

}

func (a *apiHandlers) GetMessages(w http.ResponseWriter, r *http.Request) {
	limit, page := GetDefaultQueryParams(r)
	result, err := a.repo.GetMessages(r.Context(), limit, (page-1)*limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	body, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(body)

}

func (a *apiHandlers) GetClients(w http.ResponseWriter, r *http.Request) {
	limit, page := GetDefaultQueryParams(r)
	result, err := a.repo.GetClients(r.Context(), limit, (page-1)*limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	body, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(body)

}

func (a *apiHandlers) GetUserConversations(w http.ResponseWriter, r *http.Request) {
	limit, page := GetDefaultQueryParams(r)
	result, err := a.repo.GetUserConversations(r.Context(), limit, (page-1)*limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	body, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(body)

}

func (a *apiHandlers) GetAttachments(w http.ResponseWriter, r *http.Request) {
	limit, page := GetDefaultQueryParams(r)
	result, err := a.repo.GetAttachments(r.Context(), limit, (page-1)*limit)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	body, err := json.Marshal(result)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(200)
	w.Write(body)

}
