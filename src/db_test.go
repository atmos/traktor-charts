package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitializeDatabase(t *testing.T) {
	db, _ := initializeDB("../test-tracks.db")
	assert.NotNil(t, db, "Database failed to initialize")
}
