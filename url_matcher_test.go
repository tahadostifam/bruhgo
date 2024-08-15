package mux

import (
	"log"
	"os"
	"runtime/pprof"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestUrlMatchesPattern(t *testing.T) {
	testCases := []struct {
		pattern   string
		url       string
		matches   bool
		paramsMap ParamsMap
	}{
		{
			pattern: "/main",
			url:     "/",
			matches: false,
		},
		{
			pattern:   "/api/{name}/provider/{git}",
			url:       "/api/mux/provider/github",
			matches:   true,
			paramsMap: ParamsMap{"name": "mux", "git": "github"},
		},
		{
			pattern:   "/{slug}/{name}/{age}",
			url:       "/hello_world/mux/123",
			matches:   true,
			paramsMap: ParamsMap{"slug": "hello_world", "name": "mux", "age": "123"},
		},
		{
			pattern: "/{slug}/{name}/{age}",
			url:     "/hello_world/mux",
			matches: false,
		},
		{
			pattern: "/{slug}/{name}/{age}",
			url:     "/hello_world/mux/123/extra",
			matches: false,
		},
		{
			pattern: "/{slug}/{name}",
			url:     "/hello_world/mux/123",
			matches: false,
		},
		{
			pattern:   "/{slug}/{name}",
			url:       "/hello_world/mux",
			matches:   true,
			paramsMap: ParamsMap{"slug": "hello_world", "name": "mux"},
		},
		{
			pattern:   "/{slug}/",
			url:       "/hello_world/",
			matches:   true,
			paramsMap: ParamsMap{"slug": "hello_world"},
		},
		{
			pattern:   "/",
			url:       "/",
			matches:   true,
			paramsMap: ParamsMap{},
		},
		{
			pattern: "",
			url:     "",
			matches: false,
		},
		{
			pattern:   "/{slug}",
			url:       "/hello world",
			matches:   true,
			paramsMap: ParamsMap{"slug": "hello world"},
		},
		{
			pattern: "/{slug}",
			url:     "/hello/world",
			matches: false,
		},
	}

	paramsMap := make(ParamsMap)
	for _, tc := range testCases {
		matches, resultMap, err := urlMatchesPattern(tc.pattern, tc.url, paramsMap)
		if tc.matches {
			require.NoError(t, err)
		}

		require.Equal(t, matches, tc.matches, tc.url)
		require.Equal(t, resultMap, tc.paramsMap)

		// Clear the map for the next iteration
		for k := range paramsMap {
			delete(paramsMap, k)
		}
	}
}

func BenchmarkUrlMatchesPattern(b *testing.B) {
	f, err := os.Create("url_matcher_allocs.pprof")
	if err != nil {
		log.Fatalln("Could not create file", err)
	}
	defer f.Close()

	pprof.Lookup("allocs").WriteTo(f, 0)

	paramsMap := make(ParamsMap)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		urlMatchesPattern("/{slug}", "hello_world", paramsMap)
		// Clear the map for the next iteration
		for k := range paramsMap {
			delete(paramsMap, k)
		}
	}
}
