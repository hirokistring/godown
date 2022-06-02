# hello.go

Compare the generated [hello.go](./hello.go) source file from this file.

_Code Blocks_ under a _Heading_ that ends with `.go` are translated to the global sections in go source files.

```go
package main
import "fmt"
```

```go
const WHO = "godown"
```

## func `main`

_Code Blocks_ under a _Heading_ that begins with the keyword `func` are translated to the functions in go source files.

```go
msg := sayhello(WHO)
fmt.Println(msg)
```

## func `sayhello`

- parameters
  - `who string` : someone that will be greeted
- returns
  - `string` : greeting message

```go
msg := fmt.Sprintf("Hello, %s!", who)
return msg
```

The first _Code_ in an _Item_ under _Item_ named `parameters` is translated to the function parameters.

The first _Code_ in an _Item_ under _Item_ named `returns` is translated to the return type of a function.
