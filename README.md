# sdext
a simple lib for extracting stories from slashdot

Example:

``` go get github.com/linexjlin/sdext ```

```go
package main

import (
        "fmt"

        "github.com/linexjlin/sdext"
)

func main() {
        articles := sdext.Extracter("https://slashdot.org/", "/tmp/")
        for _, v := range articles {
                fmt.Println(v)
        }
}
```
