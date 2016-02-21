package main

import (
 "os"
 "io/ioutil"
 "fmt"
 "net/http"
 "path/filepath"
 "github.com/jmars/go-duktape"
 "unsafe"
)

type DuktapeHandler struct {
  httpCode []byte
  wsCode []byte
  uid int
}

func copy_duk_buffer(__cbuf__ unsafe.Pointer, outSize int) []byte {
  bufc := (*[1<<30]byte)(__cbuf__)
  buf := make([]byte, outSize)
  for i := 0; i < outSize; i++ {
    buf[i] = bufc[i]
  }
  return buf
}

func make_context(c []byte) []byte {
  ctx := duktape.New()
  ctx.EvalString("(" + string(c) + ")")
  ctx.DumpFunction()
  outsize := 0
  bc := copy_duk_buffer(ctx.GetBuffer(-1, &outsize), outsize)
  ctx.DestroyHeap()
  return bc
}

func CallBuffer(c *duktape.Context, b []byte) {
  c.PushFixedBuffer(len(b))
  outsize := 0
  newbuf := (*[1<<30]byte)(c.GetBuffer(-1, &outsize))
  for i := 0; i < outsize; i++ {
    newbuf[i] = b[i]
  }
  c.LoadFunction()
  c.Call(0)
}

func (h *DuktapeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
  //uid := h.uid
  h.uid = h.uid + 1
  c := duktape.New()
  CallBuffer(c, h.httpCode)
  res := c.GetLstring(-1, 0)
  c.DestroyHeap()
  fmt.Fprintf(w, "%s", res)
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
    httpCode: make_context(httpCode),
    wsCode: make_context(websocketCode),
    uid: 0,
  }

  //c := duktape.New()
  //c.PushFixedBuffer(len(ctx.httpCtx))
  //outsize := 0
  //newbuf := (*[1<<30]byte)(c.GetBuffer(-1, &outsize))
  //for i := 0; i < outsize; i++ {
  //  newbuf[i] = ctx.httpCtx[i]
  //}
  //c.DestroyHeap()


  // Register our handler.
  http.ListenAndServe(":8080", &ctx)
}