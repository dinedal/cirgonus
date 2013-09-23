all: cirgonus cstat

clean:
	rm -f cirgonus cstat

cirgonus: cirgonus.go src/*/*/*.go src/*/*.go
	GOPATH=$(PWD) go build cirgonus.go

cstat: cstat.go
	GOPATH=$(PWD) go build cstat.go

cirgonus.tar.gz: cirgonus cstat
	tar cvzf cirgonus.tar.gz cirgonus cstat >/dev/null

dist: all cirgonus.tar.gz clean

distclean: clean
	rm -f cirgonus.tar.gz

run: cirgonus stop
	sh -c './cirgonus test.json &'

stop: 
	(pkill cirgonus || exit 0)
