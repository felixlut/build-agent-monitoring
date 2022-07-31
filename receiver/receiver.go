package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/jackc/pgx"
	log "github.com/sirupsen/logrus"
)

const (
	githubPath = "/github"
)

func main() {
	// retryAttempts := 10
	// retryDelay := 5

	log.Info("Parsing URI")
	connConfig, err := pgx.ParseURI(os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}

	log.Info("Connecting to database...")
	var conn *pgx.Conn
	connectToDB := func() error {
		conn, err = pgx.Connect(connConfig)
		return err
	}
	err = backoff.RetryNotify(connectToDB, backoff.NewExponentialBackOff(), func(err error, duration time.Duration) {
		log.Info(fmt.Sprintf("Error encountered, retrying in %f seconds", duration.Seconds()))
		log.Info(err)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close()
	log.Info("Connection established")

	githubHook, _ := github.New(github.Options.Secret("MyGitHubSuperSecretSecrect"))

	http.HandleFunc(githubPath, func(w http.ResponseWriter, r *http.Request) {
		payload, err := githubHook.Parse(r, github.WorkflowJobEvent)
		if err != nil {
			if err == github.ErrEventNotFound {
				log.Warning(err)
				// ok event wasn;t one of the ones asked to be parsed
			}
			log.Error(err)
		}

		switch load := payload.(type) {

		case github.WorkflowJobPayload:
			workflowJob := load
			// Do whatever you want from here...
			log.Info("%+v\n", workflowJob)

		default:
			log.Warning(fmt.Sprintf("Payload type not handled - Type: %T", load))
		}
	})
	log.Info("Listening...")
	http.ListenAndServe(":3000", nil)
}
