module github.com/ariefdarmawan/kconfigurator

go 1.16

replace git.kanosolution.net/kano/kaos => ../../lib/kaos

replace github.com/kanoteknologi/knats => ../../lib/knats

require (
	git.kanosolution.net/kano/dbflex v1.0.15
	git.kanosolution.net/kano/kaos v0.1.1
	github.com/ariefdarmawan/datahub v0.2.0
	github.com/eaciit/toolkit v0.0.0-20210610161449-593d5fadf78e
)
