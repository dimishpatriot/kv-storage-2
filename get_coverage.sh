go test -short -count=1 -race ./... -coverprofile cover.out.tmp &&
cat cover.out.tmp | grep -v "mock_" | grep -v "interface.go" > cover.out &&
rm cover.out.tmp &&
go tool cover -func cover.out