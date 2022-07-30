package main

import (
	"fmt"
	"net/http"

	"github.com/go-playground/webhooks/v6/github"
	log "github.com/sirupsen/logrus"
)

const (
	githubPath = "/github"
)

func main() {
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
			fmt.Printf("%+v\n", workflowJob)

		default:
			log.Warning(fmt.Sprintf("Payload type not handled - Type: %T", load))
		}
	})
	log.Info("Listening...")
	http.ListenAndServe(":3000", nil)
}
