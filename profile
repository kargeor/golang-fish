go build golang-fish.go
./golang-fish -cpuprofile test.prof -i
go tool pprof golang-fish test.prof

