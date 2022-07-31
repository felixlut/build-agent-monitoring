package main

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"time"

	backoff "github.com/cenkalti/backoff/v4"
	"github.com/go-playground/webhooks/v6/github"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
)

const (
	githubPath = "/github"
)

type WorkflowJobHook struct {
	workflowRunId     int64  `db:"workflow_run_id"`
	workflowJobStatus string `db:"workflow_job_status"`
	workflowJobId     int64  `db:"workflow_job_id"`
}

func main() {
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)

	var db *sqlx.DB
	var err error
	connectToDB := func() error {
		db, err = sqlx.Connect("postgres", os.Getenv("DATABASE_URL"))
		return err
	}

	log.Info("Connecting to database...")
	err = backoff.RetryNotify(connectToDB, backoff.NewExponentialBackOff(), func(err error, duration time.Duration) {
		log.Info(fmt.Sprintf("Error encountered, retrying in %f seconds", duration.Seconds()))
		log.Info(err)
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	log.Info("Connection established")

	githubHook, _ := github.New(github.Options.Secret(os.Getenv("GITHUB_WEBHOOK_SECRET")))

	http.HandleFunc(githubPath, func(w http.ResponseWriter, r *http.Request) {
		log.Debug("Webhook received")
		payload, err := githubHook.Parse(r, github.WorkflowJobEvent)
		if err != nil {
			if errors.Is(err, github.ErrEventNotFound) {
				log.Error("Does not handle this type of webhook, ignoring... ", err)
				return
			}
			log.Error("Invalid payload, ignoring... \n", err)
			return
		}
		switch load := payload.(type) {

		case github.WorkflowJobPayload:
			log.Debug(fmt.Sprintf("Recieved webhook with action: %s", load.Action))

			log.Debug("Initiating transaction")

			hoooooook := &WorkflowJobHook{load.WorkflowJob.RunID, load.Action, load.WorkflowJob.ID}
			log.Debug(hoooooook)
			tx := db.MustBegin()
			tx.MustExec(`
				INSERT INTO webhooks_workflow_job 
				(workflow_run_id, workflow_job_status, workflow_job_id) 
				VALUES ($1, $2, $3)`,
				load.WorkflowJob.RunID, load.Action, load.WorkflowJob.ID)
			err = tx.Commit()
			if err != nil {
				log.Error(err)
			}
			log.Debug("Transaction successful")
		}
	})
	log.Info("Listening...")
	http.ListenAndServe(":3000", nil)
}
