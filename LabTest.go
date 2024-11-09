package main

import (
  "fmt"
  "io"
  "math/rand"
  "os"
  "sync"
  "time"
)

func gen(f string, sizeMB int) error {
  cs := sizeMB * 1024 * 1024
  file, err := os.Create(f)
  if err != nil {
    return err
  }
  defer file.Close()

  for cs > 0 {
    c := make([]byte, cs)
    rand.Read(c)
    if _, err := file.Write(c); err != nil {
      file.Close()
      return err
    }
    cs -= len(c)
  }

  return nil
}

func r(f string) int {
  file, err := os.Open(f)
  if err != nil {
    return 0
  }
  defer file.Close()

  to := time.After(2 * time.Second)
  list := make([]byte, 0)
  d := make([]byte, 1024)

  for {
    n, err := file.Read(d)
    if err == io.EOF {
      break
    }
    select {
    case <-to:
      fmt.Println("время вышло")
      return 0
    default:
      list = append(list, d[:n]...)
    }
  }

  return len(list)
}

func main() {
  f := "r.bin"
  sizeMB := 30
  s1 := time.Now()
  var wg sync.WaitGroup
  wg.Add(1)

  go func() {
    defer wg.Done()
    err := gen(f, sizeMB)
    if err != nil {
      fmt.Println("ошибка генерации чисел:", err)
      return
    }

    l := r(f)
    fmt.Println("было просчитано с файла", l, "байт")
  }()

  wg.Wait()
  s2 := time.Now()
  fmt.Printf("время работы: %v сек.\n", s2.Sub(s1).Seconds())
}