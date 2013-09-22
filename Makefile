cirgonus: cirgonus.go
	GOPATH=$(PWD) go build cirgonus.go

cstat: cstat.go
	GOPATH=$(PWD) go build cstat.go

all: cirgonus cstat

cirgonus.tar.gz: cirgonus cstat
	tar cvzf cirgonus.tar.gz cirgonus cstat >/dev/null

dist: all cirgonus.tar.gz
