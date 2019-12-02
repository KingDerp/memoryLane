//Makai is a doodle woodle!!!
package main

import (
	"context"
	"flag"
	"memoryLane/database"
	"memoryLane/handlers"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/zeebo/errs"
)

var (
	addressFlag = flag.String(
		"address",
		":3000",
		"the address MemoryLane binds to")
)

func main() {
	flag.Parse()

	err := run(context.Background())
	if err != nil {
		logrus.Errorf("%+v\n", err)
		os.Exit(1)
	}
}

type Verse struct {
	VerseNumber int
	Text        string
}

type Scripture struct {
	Id      string
	Chapter int
	Verse   Verse
	Book    string
	Hint    string
}

func run(ctx context.Context) error {

	db, err := database.Open("postgres",
		"postgres://memorylane:something_stupid@localhost/memorylane")
	if err != nil {
		return err
	}

	handler := handlers.NewHandler(db)

	logrus.Infof("server listening on address %s\n", *addressFlag)
	return errs.Wrap(http.ListenAndServe(*addressFlag, handler))
}
