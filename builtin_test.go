package logecho

import (
	"testing"

	"github.com/pkg/profile"
)

func TestBuildParams(t *testing.T) {
	t.Run("should build json string from single param", func(t *testing.T) {
		c := NewContext(
			WithParams(map[string]string{"name": "john"}),
		)

		params := buildParams(c)
		expected := `{"name":"john"}`

		if string(params) != expected {
			t.Fatal("expected params as", expected, "but is", string(params))
		}
	})

	t.Run("should build json string from params", func(t *testing.T) {
		c := NewContext(
			WithParams(map[string]string{
				"name": "john",
				"id":   "1",
			}),
		)

		params := buildParams(c)
		expected := `{"name":"john","id":"1"}`

		if string(params) != expected {
			t.Fatal("expected params as", expected, "but is", string(params))
		}
	})
}

func BenchmarkBuildParams(b *testing.B) {
	c := NewContext(
		WithParams(map[string]string{
			"name": "john",
			"id":   "1",
		}),
	)

	defer profile.Start(
		profile.MemProfile,
		profile.MemProfileAllocs,
		profile.MemProfileHeap,
		profile.ProfilePath("."),
	).Stop()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		buildParams(c)
	}
}
