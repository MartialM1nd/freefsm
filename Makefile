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
	install -d ${PREFIX}/share/freefsm/static
	cp -R ui/static/* ${PREFIX}/share/freefsm/static/
	install -m 755 deploy/freebsd/freefsm ${PREFIX}/etc/rc.d/
	install -m 644 deploy/freebsd/freefsm.conf.sample ${PREFIX}/share/freefsm/

uninstall:
	rm -f ${PREFIX}/bin/${PROG}
	rm -f ${PREFIX}/etc/rc.d/freefsm
	rm -rf ${PREFIX}/share/freefsm

.PHONY: all build run migrate clean install uninstall
