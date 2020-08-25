
echo [tool fix]
go tool fix -diff .

echo
echo [fmt]
go fmt ./...

echo
echo [lint]
go get golang.org/x/lint/golint >/dev/null 2>/dev/null
golint -min_confidence 0 ./... |
  grep -Ev "exported (.+)?should have comment (.+)?or be unexported" |
  grep -Ev "should have a package comment, unless it's in another file for this package"

echo
echo [staticcheck]
go get honnef.co/go/tools/cmd/staticcheck >/dev/null 2>/dev/null
staticcheck ./...

echo
echo [vet]
go vet ./...

echo
echo [mod tidy]
go mod tidy -v 2>&1 |
  grep -v 'unused golang.org/x/lint' |
  grep -v 'unused honnef.co/go/tools'

echo
echo [mod verify]
go mod verify

echo
echo [test]
go test -v ./...
