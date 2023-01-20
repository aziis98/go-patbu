package patbu_test

import (
	"testing"

	"github.com/aziis98/go-patbu"
	"gotest.tools/assert"
)

func TestExact(t *testing.T) {
	patbu := &patbu.Patbu{patbu.Exact{"a.txt"}}

	t.Run("Match/Good", func(t *testing.T) {
		m, err := patbu.Match("a.txt")
		if err != nil {
			t.Fatal(err)
		}

		assert.DeepEqual(t, m, map[string]string{})
	})
	t.Run("Match/Bad", func(t *testing.T) {
		m, err := patbu.Match("b.txt")
		if err == nil {
			t.Fatal(m)
		}

		t.Log(err)
	})
	t.Run("Build/Good", func(t *testing.T) {
		s, err := patbu.Build(map[string]string{})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, s, "a.txt")
	})
}

func TestFileVar(t *testing.T) {
	p := &patbu.Patbu{patbu.FileVar{"foo"}}

	t.Run("Match-Good", func(t *testing.T) {
		m, err := p.Match("a.txt")
		if err != nil {
			t.Fatal(err)
		}

		assert.DeepEqual(t, m, map[string]string{
			"foo": "a.txt",
		})
	})
	t.Run("Match/Dirs/Bad", func(t *testing.T) {
		m, err := p.Match("a/b.txt")
		if err == nil {
			t.Fatal(m)
		}

		t.Log(err)
	})
	t.Run("Build/Good", func(t *testing.T) {
		s, err := p.Build(map[string]string{
			"foo": "b.txt",
		})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, s, "b.txt")
	})
	t.Run("Build/Empty-Bad", func(t *testing.T) {
		s, err := p.Build(map[string]string{})
		if err == nil {
			t.Fatal(s)
		}

		t.Log(err)
	})
	t.Run("Build/Dirs/Bad", func(t *testing.T) {
		s, err := p.Build(map[string]string{
			"foo": "a/b.txt",
		})
		if err == nil {
			t.Fatal(s)
		}

		t.Log(err)
	})
}

func TestDirsVar(t *testing.T) {
	p := &patbu.Patbu{patbu.DirsVar{"foo"}}

	t.Run("Match/Good", func(t *testing.T) {
		m, err := p.Match("a.txt")
		if err != nil {
			t.Fatal(err)
		}

		assert.DeepEqual(t, m, map[string]string{
			"foo": "a.txt",
		})
	})
	t.Run("Build/Good", func(t *testing.T) {
		s, err := p.Build(map[string]string{
			"foo": "b.txt",
		})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, s, "b.txt")
	})
	t.Run("Build/Bad", func(t *testing.T) {
		s, err := p.Build(map[string]string{})
		if err == nil {
			t.Fatal(s)
		}

		t.Log(err)
	})
	t.Run("Build/Bad", func(t *testing.T) {
		s, err := p.Build(map[string]string{
			"foo": "a/b.txt",
		})
		if err != nil {
			t.Fatal(err)
		}

		assert.Equal(t, s, "a/b.txt")
	})
}

func TestParse(t *testing.T) {
	t.Run("Simple", func(t *testing.T) {
		p, err := patbu.Parse(`a.txt`)
		if err != nil {
			t.Fatal(err)
		}

		assert.DeepEqual(t, p, patbu.Patbu{patbu.Exact{"a.txt"}})
	})
	t.Run("SimpleDirs", func(t *testing.T) {
		p, err := patbu.Parse(`a/b/c.txt`)
		if err != nil {
			t.Fatal(err)
		}

		assert.DeepEqual(t, p, patbu.Patbu{patbu.Exact{"a/b/c.txt"}})
	})
	t.Run("FileVar", func(t *testing.T) {
		p, err := patbu.Parse(`a/b/{foo}.txt`)
		if err != nil {
			t.Fatal(err)
		}

		assert.DeepEqual(t, p, patbu.Patbu{
			patbu.Exact{"a/b/"},
			patbu.FileVar{"foo"},
			patbu.Exact{".txt"},
		})

	})
	t.Run("NearPatterns/Good", func(t *testing.T) {
		p, err := patbu.Parse(`{a}-{b}.txt`)
		if err != nil {
			t.Fatal(err)
		}

		assert.DeepEqual(t, p, patbu.Patbu{
			patbu.FileVar{"a"},
			patbu.Exact{"-"},
			patbu.FileVar{"b"},
			patbu.Exact{".txt"},
		})
	})
	t.Run("NearPatterns/Bad", func(t *testing.T) {
		p, err := patbu.Parse(`{a}{b}.txt`)
		if err == nil {
			t.Fatal(p)
		}

		t.Log(err)
	})
	t.Run("DirsVar", func(t *testing.T) {
		p, err := patbu.Parse(`src/{*route}/index.html`)
		if err != nil {
			t.Fatal(err)
		}

		assert.DeepEqual(t, p, patbu.Patbu{
			patbu.Exact{"src/"},
			patbu.DirsVar{"route"},
			patbu.Exact{"/index.html"},
		})
	})
}
