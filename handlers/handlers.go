package handlers

import (
    "encoding/json"
    "net/http"
    "strconv"
    "strings"
    "sync"
    "time"

    "go-rest-api/models"
)

type Store struct {
    mu     sync.RWMutex
    items  map[int]models.Item
    nextID int
}

func NewStore() *Store { return &Store{items: make(map[int]models.Item), nextID: 1} }

func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func (s *Store) HandleItems(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        s.mu.RLock(); defer s.mu.RUnlock()
        list := make([]models.Item, 0, len(s.items))
        for _, v := range s.items { list = append(list, v) }
        writeJSON(w, http.StatusOK, list)
    case http.MethodPost:
        var body struct{ Name, Value string }
        if err := json.NewDecoder(r.Body).Decode(&body); err != nil { http.Error(w, "bad request", http.StatusBadRequest); return }
        s.mu.Lock(); defer s.mu.Unlock()
        item := models.Item{ID: s.nextID, Name: body.Name, Value: body.Value, CreatedAt: time.Now()}
        s.items[s.nextID] = item; s.nextID++
        writeJSON(w, http.StatusCreated, item)
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}

func (s *Store) HandleItem(w http.ResponseWriter, r *http.Request) {
    parts := strings.Split(r.URL.Path, "/")
    id, err := strconv.Atoi(parts[len(parts)-1])
    if err != nil { http.Error(w, "invalid id", http.StatusBadRequest); return }

    switch r.Method {
    case http.MethodGet:
        s.mu.RLock(); defer s.mu.RUnlock()
        item, ok := s.items[id]
        if !ok { http.Error(w, "not found", http.StatusNotFound); return }
        writeJSON(w, http.StatusOK, item)
    case http.MethodDelete:
        s.mu.Lock(); defer s.mu.Unlock()
        if _, ok := s.items[id]; !ok { http.Error(w, "not found", http.StatusNotFound); return }
        delete(s.items, id)
        w.WriteHeader(http.StatusNoContent)
    default:
        http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
    }
}
