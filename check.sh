
echo [mod tidy]
go mod tidy -v

echo;echo [fmt]
go fmt ./...

echo;echo [mod verify]
go mod verify

echo;echo [lint]
go get golang.org/x/lint/golint >/dev/null 2>/dev/null
golint -min_confidence 0 ./... \
  | grep -Ev "exported (.+)?should have comment (.+)?or be unexported" \
  | grep -Ev "should have a package comment, unless it's in another file for this package"
go mod tidy

echo;echo [vet]
go vet ./...

echo;echo [tool fix]
go tool fix -diff .

echo;echo [test]
go test -v ./...
