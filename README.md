# e4
error utilities version 4

### Example


```
func CopyFile(src, dst string) (err error) {
  defer Handle(&err,
    WithInfo("copy %s to %s", src, dst),
  )

  r, err := os.Open(src)
  Check(err, WithInfo("open %s", src))
  defer r.Close()

  w, err := os.Create(dst)
  Check(err, WithInfo("create %s", dst))
  defer w.Close()
  defer Handle(&err, WithFunc(func() {
    os.Remove(dst)
  }))

  _, err = io.Copy(w, r)
  Check(err)
  Check(w.Close())

  return
}

```
