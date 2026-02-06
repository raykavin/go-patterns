# Go Patterns Examples (Functional)

This repo contains **36 runnable examples** of idiomatic Go patterns.

## Run a single example
From the repo root:

```bash
go run ./concurrency/fanout_fanin
```

(or any other folder)

## Run all examples quickly (one by one)
```bash
for d in $(find . -mindepth 2 -maxdepth 2 -type d); do
  echo "==> $d"
  go run "$d" || exit 1
done
```

Each example is a standalone `package main` with its own folder.
