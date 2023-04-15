# matchers

`matchers` is a simple package to parse Prometheus-like matchers.

## Example

```go
package main

import (
	"fmt"
	"os"

	"github.com/grobinson-grafana/matchers"
	"github.com/prometheus/common/model"
)

func main() {
	labels := model.LabelSet{"foo": "bar"}
	m, err := matchers.Parse("{foo=~\"[a-z]+\"}")
	if err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
	ok := m.Matches(labels)
	fmt.Println(ok)
}
```