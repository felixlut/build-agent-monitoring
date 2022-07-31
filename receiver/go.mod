module build-agent-monitor

go 1.18

require (
	github.com/cenkalti/backoff/v4 v4.1.3
	github.com/go-playground/webhooks/v6 v6.0.1
	github.com/jmoiron/sqlx v1.3.5
	github.com/lib/pq v1.10.6
	github.com/sirupsen/logrus v1.9.0
)

require golang.org/x/sys v0.0.0-20220730100132-1609e554cd39 // indirect
