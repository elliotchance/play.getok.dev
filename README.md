# play.getok.dev

This repository represents the [play.getok.dev](https://play.getok.dev) domain.

## Deploy

```bash
make deploy
```

## Testing

```
sls invoke -f run --data '{"body":"func main() { print(\"hello\") }"}'
```

```
curl --data 'func main() { print("hello") }' \
  https://fh0504kns1.execute-api.us-east-1.amazonaws.com/dev/run
```
