
all: 
	go build -o ./bin/cliscratchy ./scratchy/cliscratch.go
	go build -o ./bin/tliscratchy ./scratchy/tliscratch.go
	go build -o ./bin/hello ./hello.go
	go build -o ./bin/scratchy ./scratchy/neoscratch.go
	go build -o ./bin/demo ./scratchy/demo.go
	go build -o ./bin/widgets ./scratchy/widgets.go
