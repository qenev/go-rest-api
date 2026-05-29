package models

import "time"

type Item struct {
    ID        int       `json:"id"`
    Name      string    `json:"name"`
    Value     string    `json:"value"`
    CreatedAt time.Time `json:"created_at"`
}
