# Testing

## Mocking

Requied pack gomock. Package have 'mockgen' to generate mocks.

Example, how generate mocks with 'mockgen':
``` bash
mockgen -destination=mocks/actuator.go -package=mocks github.com/gardener/gardener/extensions/pkg/controller/extension Actuator
mockgen -destination=mocks/client.go -package=mocks sigs.k8s.io/controller-runtime/pkg/client Client
```

## Coverage

Below example how generate coverage for tests. Second line generate readable website. Best for check html from remote server is LiveServer, author ritwickdey.
``` go
go test -coverprofile=coverage.out ./pkg/...
go tool cover -html=coverage.out -o coverage.html           
```
