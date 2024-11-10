package model

type User struct {
    ID        string `json:"id"`
    Email     string `json:"email"`
    Username  string `json:"username"`
    CreatedAt string `json:"created_at"`
    UpdatedAt string `json:"updated_at"`
} 