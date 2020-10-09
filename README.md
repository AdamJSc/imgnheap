# imgnheap

## About

Bulk-process your heap of images into lovely sub-folders, by way of a local Web UI

## Requirements

* Golang 1.14

## Getting started

From project root:

```
go run service/main.go
```

## Updating Templates

Requires the Pkger CLI (https://github.com/markbates/pkger)

From project root:

```
pkger -o service/views
```

## Running Tests

From project root:

```
go test ./...
```
