package users

import "time"

type User struct {
    ID int    `json:"id" db:"id"`
    Name        string `json:"name" db:"name"`
    Email       string `json:"email" db:"email"`
    Password    string `json:"-" db:"password_hash"`
    PhoneNumber string `json:"phone_number" db:"phone_number"`
    CreatedAt time.Time `db:"created_at" json:"-"`
    RoleID      int    `json:"role_id" db:"role_id"`
}
