package tests

import (
	"github.com/go-pg/pg"
)

var (
	testDB             *pg.DB
	terminateContainer = func() {}
)
