package main

import (
	"database/sql"
)

type Storage struct {
	Conn *sql.DB
}

type Pair struct {
	ID       string
	JwtToken string
	CACert   []byte
}

func (s *Storage) GetPair(deviceID string) (*Pair, error) {
	// Implementation here
	return nil, nil
}

func (s *Storage) AddPair(deviceId string, jwtToken string, caCert []byte) error {
	// Implementation here
	return nil
}

func (s *Storage) RemovePair(deviceId string) error {
	// Implementation here
	return nil
}
