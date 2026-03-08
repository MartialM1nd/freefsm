PREFIX?=	/usr/local
PROG=		freefsm

all: build

build:
	go build -o ${PROG} ./cmd/freefsm

run:
	go run ./cmd/freefsm

migrate:
	go run ./cmd/freefsm -migrate

clean:
	rm -f ${PROG}

install: build
	install -m 755 ${PROG} ${PREFIX}/bin/
	install -d ${PREFIX}/share/freefsm/templates
	install -d ${PREFIX}/share/freefsm/static
	cp -R ui/templates/* ${PREFIX}/share/freefsm/templates/
	cp -R ui/static/* ${PREFIX}/share/freefsm/static/

uninstall:
	rm -f ${PREFIX}/bin/${PROG}
	rm -rf ${PREFIX}/share/freefsm

.PHONY: all build run migrate clean install uninstall
