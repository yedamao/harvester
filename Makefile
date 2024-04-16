bin: clean
	go build -o ./bin/harvester ./cmd/harvester
	go build -o ./bin/planter ./cmd/planter

clean:
	rm -rf ./bin
