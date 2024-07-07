## Contributing

Thanks for your interest in contribute to dynmgrm!


### Getting Started

1. [Fork `miyamo2/dynmgrm`](https://github.com/miyamo2/dynmgrm/fork)

2. Clone your fork repository locally

```sh
git clone https://github.com/{YOUR_USERNAME}/dynmgrm.git
```
3. Check out to the working branch

```sh
git checkout -b cool_branch_name
```

### Unit Test

```sh
go test -v ./... 
```

### Integration Test with DynamoDB Local

```sh
cd ./integrationtest
make test
```
