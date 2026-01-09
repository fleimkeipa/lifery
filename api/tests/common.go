package tests

import (
	"github.com/go-pg/pg/v10"
)

var (
	testDB             *pg.DB
	terminateContainer = func() {}
)
