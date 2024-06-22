module subscriber

go 1.22.4

require (
	github.com/go-redis/redis/v8 v8.11.5
	project/common v0.0.0-00010101000000-000000000000
)

require (
	github.com/fsnotify/fsnotify v1.6.0 // indirect
	golang.org/x/net v0.24.0 // indirect
	golang.org/x/sys v0.20.0 // indirect
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
)

replace project/common => ../common
