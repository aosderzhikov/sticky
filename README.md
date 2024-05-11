# sticky

## Description

The storing service is pepresented with two components - `bouncer` and `keeper`.

### Keeper

`Keeper` is simple hash map service with `REST API` for `set`, `get` and `delete` operations. It stores key:value pairs in `RAM` memory, and protects pairs with `mutex`, so no worry about consistency.
Each pair could be set with `ttl` or will be used default `ttl` for the whole service. After ttl expiration entry will be automaticly removed.

### Bouncer

`Bouncer` is `load balancer` for keepers. It knows about keepers which he controls, and looking after their status via `health check`. When request to store key:value pair comes to `bouncer` it decides which `keeper` will stores this pair. It decides where to store pair by algorythms: 

- hash func from storing key
- round robin  `// todo`

If detected storage unavailable, `bouncer` will try to put pair in first alive storage. The same behaivor with updating: try to put in storage with actual key, try to put in storage detected with algorythm, try to put in first alive 

`Bouncer` stores index with pairs key:storage_index, so it knows where from to take storing value or update. 


## Deployment
### Standalone

If you need just one instance of storage it enough to deploy only one `keeper`.

Available environment variables:
- `HTTP_ADDRESS` address for `keeper` deployment, default `localhost:8181`
- `TTL` ttl for entries, uses when doesnt pass in request, default `10m`
- `DEBUG` debug mod, default `false`


To run `keeper` use

```sh
go run . cmd/keeper/main.go
```

Or build it first
```sh
go build -o keeper  cmd/keeper/main.go

./keeper
```

After it is avalaible for usage
```sh
## set value
curl -X POST 'http://localhost:8181/set?key=key1&ttl=1m' -d 'storing_value' 

## get value
curl 'http://localhost:8181/get?key=key1'

## delete value 
curl -X DELETE 'http://localhost:8181/delete?key=key1'
```

# TODO
tests
load actual data on failed for balancer
inmemory
replica
hash func
round robin func
swagger