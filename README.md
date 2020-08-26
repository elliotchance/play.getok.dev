# getok.dev

This repository represents the [getok.dev](https://getok.dev) domain.

## Deploy

```bash
make deploy
```

## play.getok.dev

```
sls invoke -f run --data '{"body":"func main() { print(\"hello\") }"}'
```

```
curl --data 'func main() { print("hello") }' \
  https://fh0504kns1.execute-api.us-east-1.amazonaws.com/dev/run
```
