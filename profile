go build
./golang-fish -cpuprofile test.prof -i
go tool pprof golang-fish test.prof

