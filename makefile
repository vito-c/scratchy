
all: 
	go build -o ./bin/cscratchy ./scratchy/cliscratch.go
	go build -o ./bin/hello ./hello.go
	go build -o ./bin/scratchy ./scratchy/neoscratch.go
