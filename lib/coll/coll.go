package coll

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

func TransformKeys(coll map[string]string, f func(s string) string) map[string]string {
	collTransformed := make(map[string]string, len(coll))
	for k, v := range coll {
		collTransformed[f(k)] = v
	}
	return collTransformed
}
