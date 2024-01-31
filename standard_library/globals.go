package standardlibrary

import (
	"fmt"
	"os"
	"strings"
  "bufio"
)

func Len[A any, V string | []A](v V) int64 {
  return int64(len(v))
}

func Trim(v string) string {
  return strings.TrimSpace(v)
}

func Print(args ...any) {
  fmt.Println(args...)
}

func Input() string {
  reader := bufio.NewReader(os.Stdin)
  line, _, _ := reader.ReadLine()
  return string(line)
}