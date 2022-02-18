package util

import (
	"fmt"
	"strings"
)

func JoinEntries(m map[string]string, glue string) []string {
	var vsf []string
	for k, v := range m {
		vsf = append(vsf, fmt.Sprintf("%s%s%s", strings.TrimSpace(k), glue, v))
	}
	return vsf
}
