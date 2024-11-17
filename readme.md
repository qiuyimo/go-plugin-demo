```bash
go build -o kv
go build -o kv-go-grpc ./plugin
export KV_PLUGIN="./kv-go-grpc"
./kv put hello world
./kv get hello 
```