# harvester
redis scan tool

## Install
```
go install github.com/yedamao/harvester/cmd/harvester@latest
```

## Usage
```
harvester -h
Usage of harvester:
  -action string
        dump the key data (default "dump")
  -matchPattern string
        scan match pattern (default "*")
  -parallel int
        the number of worker to run parallel (default 1)
  -redisUrl string
        eg: redis://<user>:<password>@<host>:<port>/<db_number>. separated by commas.
```
