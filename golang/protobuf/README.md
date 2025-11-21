# protobuf tutorial in Go

Link to [tutorial](https://protobuf.dev/getting-started/gotutorial/).

## Compiling .proto

You have to install `protoc`, ie: `brew install protobuf`.

Then install protoc-gen-go, ie: `go install google.golang.org/protobuf/cmd/protoc-gen-go@latest`

Now run protobuf compiler:

```bash
protoc -I=. --go_out=. addressbook.proto
```