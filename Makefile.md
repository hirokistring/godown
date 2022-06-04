# Makefile of `godown`

This is the `Makefile` written in markdown for `godown`.

This `Makefile.md` file can be run by [makedown](https://github.com/hirokistring/makedown).

## About `makedown:`

Note that this `Makefile.md` assumes that you have `gmake` 3.82+ to use `.ONESHELL`.

Make sure you have a required version of `make`,
especially on MacOS. Then, type:

`export MAKEDOWN_MAKE_COMMAND=gmake`

```
#!/usr/local/bin/gmake -f

.ONESHELL:
```

## How to `generate:`

This generates `*.go` from `*.md`.

The debug logs are saved in `verbose.log`.
Check it out.

```sh
./godown --verbose 2&> verbose.log
```

## `diff:`

Compare the generated `*.go` and `*.go.expected` files.

```sh
diff godown.go godown.go.expected || true
diff main.go main.go.expected || true
```

## How to just `build:`

```sh
go build
```

The generated `godown` executable file is not ready to be released. It has to be notarized by [gon](https://github.com/mitchellh/gon) for MacOS.

## Another way to build, `build-by-godown:`

Let's try to build `godown` by `godown` itself.

```sh
./godown build
```

## How to `build-and-notarize:`

This makes binaries for each platform. Then, it signs and notarizes the binary for MacOS.

```sh
goreleaser build --snapshot --rm-dist
```

### How to `check-notarized:`

```sh
@echo Check the binary is signed
codesign --display -vvv dist/macos_darwin_amd64_v1/godown

@echo
@echo Check the binary is notarized
spctl --assess --type install -vvv dist/macos_darwin_amd64_v1/godown
```

Note that `godown` is just a binary, not an .app.

## How to try `build-and-release` with `build-and-release-snapshot:`

Type the command bellow.

```sh
goreleaser release --snapshot --rm-dist
```

## How to`build-and-release:`

Type the command bellow.

```sh
goreleaser release --rm-dist
```

## Build and Release Tools

`godown` uses [goreleaser](https://github.com/goreleaser) to build and release the binary.

`godown` uses [gon](https://github.com/mitchellh/gon) to notarize the binary for Mac.

## Tests

Make sure that you have a make command 3.82+ before running the tests.

### `hello:`

```sh
cd tests/hello
# clean up
if [ -e hello.go ]; then
  rm hello.go
fi
if [ -e hello ]; then
  rm hello
fi
# generate *.go files
../../godown
# check diff
diff hello.go hello.go.expected > hello.go.diff
cat hello.go.diff
# build the generated *.go files
../../godown build
# run it
./hello
# another way to run
../../godown run hello.go.md
```
