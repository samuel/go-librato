include $(GOROOT)/src/Make.inc

TARG=librato
GOFILES=\
	librato.go\
	metrics.go\
	query.go\
	services.go\
	users.go\

include $(GOROOT)/src/Make.pkg