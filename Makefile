bin: clean
	go build -o ./bin/harvester ./cmd/harvester

clean:
	rm -rf ./bin
