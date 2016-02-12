package main

import (
 "os"
 "io/ioutil"
 "fmt"
 "net/http"
 "path/filepath"
 "github.com/jmars/go-duktape"
)

type DuktapeHandler struct {
  name string
  httpCtx []byte
  wsCtx []byte
}

func make_context(c []byte) []byte {
  ctx := duktape.New()
  ctx.EvalString(string(c))
  ctx.DumpFunction()
  outsize := 0
  bc := ctx.GetBuffer(-1, &outsize)
  ctx.DestroyHeap()
  return bc
}

func (h *DuktapeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  name := h.name
  // Write it back to the client.
  fmt.Fprintf(w, "hi %s!\n", name)
}

func main() {
  if len(os.Args) == 1 {
    println("Must provide a package path.")
    return
  }

  folder := os.Args[1]

  httpCode, err := ioutil.ReadFile(filepath.Join(folder, "http.js"))
  if err != nil {
    println("Couldn't find http.js in provided folder.")
    return
  }

  websocketCode, err := ioutil.ReadFile(filepath.Join(folder, "websocket.js"))
  if err != nil {
    println("Couldn't find websocket.js in provided folder.")
    return
  }

  ctx := DuktapeHandler{
    name: "test",
    httpCtx: make_context(httpCode),
    wsCtx: make_context(websocketCode),
  }

  // Register our handler.
  http.ListenAndServe(":8080", &ctx)
}