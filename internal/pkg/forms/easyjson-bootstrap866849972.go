// +build ignore

// TEMPORARY AUTOGENERATED FILE: easyjson bootstapping code to launch
// the actual generator.

package main

import (
  "fmt"
  "os"

  "github.com/mailru/easyjson/gen"

  pkg "failless/internal/pkg/forms"
)

func main() {
  g := gen.NewGenerator("form_meta_easyjson.go")
  g.SetPkg("forms", "failless/internal/pkg/forms")
  g.Add(pkg.EasyJSON_exporter_MetaForm(nil))
  if err := g.Run(os.Stdout); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
