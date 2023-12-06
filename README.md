# kazuyomi
A utility library to get the Japanese reading of numbers.

数字の読み仮名を取得するユーティリティ

## Installation

```
go get github.com/jiftechnify/kazuyomi
```

## Examples

```go
package main

import "github.com/jiftechnify/kazuyomi"

func main() {
    // input can contain ',' or '_' as separators.
    yomi, _ := kazuyomi.ReadString("5,000,000,000,000")
    fmt.Println(yomi) // "ゴセンチョウ"

    fmt.Println(kazuyomi.ReadInt(-3600))    // "マイナスサンゼンロッピャク"
    fmt.Println(kazuyomi.ReadUint(8e16))    // "ハッケイ"
    fmt.Println(kazuyomi.ReadFloat64(0.1))  // "レイテンイチ"
}
```

## License

MIT
