go test -bench=".*" --cpuprofile=perf.prof && echo "web" | go tool pprof lists.test perf.prof
