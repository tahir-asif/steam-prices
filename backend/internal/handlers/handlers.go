// Package handlers is a where API request handlers go
package handlers

import (
    "database/sql"
)

type Handler struct {
    DB *sql.DB
}
