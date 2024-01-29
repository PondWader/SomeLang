package standardlibrary

import "strings"

func Len[A any, V string | []A](v V) int64 {
  return int64(len(v))
}

func Trim(v string) string {
  return strings.TrimSpace(v)
}