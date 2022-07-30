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
	githubHook, _ := github.New(github.Options.Secret("MyGitHubSuperSecretSecrect...?"))

	http.HandleFunc(githubPath, func(w http.ResponseWriter, r *http.Request) {
		payload, err := githubHook.Parse(r, github.ReleaseEvent, github.PullRequestEvent)
		log.Info("Payload:\n", payload)
		if err != nil {
			if err == github.ErrEventNotFound {
				log.Warning(err)
				// ok event wasn;t one of the ones asked to be parsed
			}
		}

		switch payload.(type) {

		case github.ReleasePayload:
			release := payload.(github.ReleasePayload)
			// Do whatever you want from here...
			fmt.Printf("%+v", release)

		case github.PullRequestPayload:
			pullRequest := payload.(github.PullRequestPayload)
			// Do whatever you want from here...
			fmt.Printf("%+v", pullRequest)

		default:
			log.Warning(fmt.Sprintf("Payload type not handled - Type: %T", payload))
		}
	})
	log.Info("Listening...")
	http.ListenAndServe(":3000", nil)
}
