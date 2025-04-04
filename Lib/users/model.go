package users

import "time"

type User struct {
    ID        int       `db:"id" json:"id"`
    Name      string    `db:"name" json:"name"`
    Email     string    `db:"email" json:"email"`
    CreatedAt time.Time `db:"created_at" json:"created_at"`
    RoleID    int       `db:"role_id" json:"role_id"`
}

