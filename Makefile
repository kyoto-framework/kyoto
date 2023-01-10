

# Build client code
client:
	(cd client && npm i && npm run client)

# Generic test
test:
	go test -race .

# Coverage
testcov:
	go test -cover .

# Coverage with web report
testcovweb:
	go test -coverprofile /tmp/kyoto-coverage.out .
	go tool cover -html=/tmp/kyoto-coverage.out
	sleep 3 && rm /tmp/kyoto-coverage.out

# Serve docs
doc:
	(sleep 1 && open http://localhost:8000/pkg/github.com/kyoto-framework/kyoto/v2) &
	godoc -http=:8000
