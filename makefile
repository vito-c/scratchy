
run: all
	./bin/cliscratchy < test.json

all: 
	go build -o ./bin/cliscratchy ./cli/scratchy.go
	go build -o ./bin/scratchy ./nvim/neoscratch.go
	# go build -o ./bin/demo ./test/demo.go
	# go build -o ./bin/widgets ./test/wdemo.go
