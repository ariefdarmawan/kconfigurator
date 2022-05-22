module github.com/ariefdarmawan/kconfigurator

go 1.16

// replace git.kanosolution.net/kano/kaos => ../../lib/kaos

// replace github.com/kanoteknologi/knats => ../../lib/knats

replace git.kanosolution.net/koloni/crowd => git.kanosolution.net/koloni/crowd v0.0.1

require (
	git.kanosolution.net/kano/dbflex v1.0.16
	git.kanosolution.net/kano/kaos v0.2.0
	github.com/ariefdarmawan/datahub v0.2.0
	github.com/eaciit/toolkit v0.0.0-20210610161449-593d5fadf78e
)
