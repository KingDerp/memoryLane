package main

import (
	"context"
	"flag"
	"net/http"

	"memoryLane/database"
	"memoryLane/handlers"

	"github.com/sirupsen/logrus"
	"github.com/zeebo/errs"
)

var (
	addressFlag = flag.String(
		"address",
		":3000",
		"the address MemoryLane binds to")
	dbFlag = flag.String(
		"db",
		"postgres://memorylane:something_stupid@localhost/memorylane",
		"the connection string to the desired database")
)

func main() {
	flag.Parse()

	err := run(context.Background())
	if err != nil {
		logrus.Fatalf("%+v\n", err)
	}
}

func run(ctx context.Context) error {

	db, err := database.Open("postgres", dbFlag)
	if err != nil {
		return err
	}

	handler := handlers.NewHandler(db)

	logrus.WithField("address", *addressFlag).Info("server listening")
	return errs.Wrap(http.ListenAndServe(*addressFlag, handler))
}
