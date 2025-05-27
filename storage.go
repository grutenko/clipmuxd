package main

import "database/sql"

type Storage struct {
	Conn *sql.DB
}
