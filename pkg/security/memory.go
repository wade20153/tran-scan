package security

import "runtime"

// ZeroString 尝试抹除敏感字符串，防止内存 dump
func ZeroString(s *string) {
	b := []byte(*s)
	for i := range b {
		b[i] = 0
	}
	*s = ""
	runtime.GC()
}
