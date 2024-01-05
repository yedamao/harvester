# harvester
redis scan tool. dump data to file

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

## function
1. dump keys with idle time
2. dump `Strings` as key urlencode(val)
3. dump `Hashes` as key field urlencode(val)
4. dump `Sorted sets` as key score urlencode(member)
5. Other types TODO
