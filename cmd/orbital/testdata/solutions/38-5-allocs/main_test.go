package main

import "testing"

var assemblySink string

func BenchmarkAssemble(b *testing.B) {
	parts := []string{"orbit", "relay", "dock"}
	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		assemblySink = Assemble(parts)
	}
}
