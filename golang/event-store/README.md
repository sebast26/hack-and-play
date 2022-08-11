# Event Store based on DynamoDB

## To run

```bash
$(npm bin)/sls package
$(npm bin)/sls deploy
```

## To execute integration tests
```bash
make local-dynamo && make test
```