startfile = "cmd/hn.go"
run:
	go run $(startfile) 

install:
	go install $(startfile) 