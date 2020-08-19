package main

import "testing"

func TestResolveStaticPath(t *testing.T) {
	type testCase struct {
		root   string
		prefix string
		strip  bool
		path   string
		output string
	}

	var cases []testCase = []testCase{
		{"www", "/foo/", false, "/foo/bleem.txt", "www/foo/bleem.txt"},
		{"www", "/foo/", true, "/foo/bleem.txt", "www/bleem.txt"},

		{"files", "/", false, "/a/b/c.txt", "files/a/b/c.txt"},
		{"files", "/", true, "/a/b/c.txt", "files/a/b/c.txt"},

		{".", "/", false, "/a/b/c.txt", "a/b/c.txt"},
		{".", "/", true, "/a/b/c.txt", "a/b/c.txt"},
	}

	for ix, c := range cases {
		h := NewStaticHandler(&RouteConfig{
			Static: &StaticRouteConfig{
				Path:        c.root,
				StripPrefix: c.strip,
			},
			Prefix: c.prefix,
		})

		result, _ := h.resolvePath(c.path)
		if result != c.output {
			t.Errorf("TestResolveStaticPath() - case %d failed (expect=%s, actual=%s)", ix, c.output, result)
			t.Fail()
		}
	}
}
