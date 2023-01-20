package patbu

import (
	"fmt"
	"io"
	"strings"
	"unicode/utf8"
)

type Binding struct{ Key, Value string }

type Part interface {
	Match(s string, next Part) (string, *Binding, error)
	Build(w io.Writer, m map[string]string) error
}

type Patbu []Part

func (patbu Patbu) Match(s string) (map[string]string, error) {
	var (
		binding *Binding
		err     error
	)

	m := map[string]string{}

	for i, part := range patbu {
		var next Part = nil
		if i < len(patbu)-1 {
			next = patbu[i+1]
		}

		s, binding, err = part.Match(s, next)
		if err != nil {
			return nil, err
		}
		if binding != nil {
			m[binding.Key] = binding.Value
		}
	}

	if len(s) > 0 {
		return nil, fmt.Errorf(`expected end of string but got "...%s"`, s)
	}

	return m, nil
}

func (patbu Patbu) Build(m map[string]string) (string, error) {
	sb := &strings.Builder{}

	for _, part := range patbu {
		if err := part.Build(sb, m); err != nil {
			return "", err
		}
	}

	return sb.String(), nil
}

var _ Part = Exact{}

type Exact struct{ Value string }

func (p Exact) Match(s string, _ Part) (string, *Binding, error) {
	if !strings.HasPrefix(s, p.Value) {
		return s, nil, fmt.Errorf(`expected "%s" but got "...%s"`, p.Value, s)
	}

	return s[len(p.Value):], nil, nil
}

func (p Exact) Build(w io.Writer, m map[string]string) error {
	fmt.Fprint(w, p.Value)
	return nil
}

var _ Part = FileVar{}

type FileVar struct{ Name string }

func (p FileVar) Match(s string, next Part) (string, *Binding, error) {
	// this capture is at the end of the pattern
	if next == nil {
		i := strings.IndexRune(s, '/')
		if i != -1 {
			return "", nil, fmt.Errorf(`expected file part but got "/" at "...%s"`, s[i:])
		}

		return "", &Binding{Key: p.Name, Value: s}, nil
	}

	nextExact, ok := next.(Exact)
	if !ok {
		// this cannot happen, the parser will already reject "...{a}{b}..."
		panic(`cannot have two captures one next to the other`)
	}

	untilRune, _ := utf8.DecodeRuneInString(nextExact.Value)
	i := strings.IndexAny(s, string(untilRune)+"/")
	if i == -1 {
		return s, nil, fmt.Errorf(`expected "%v" but got "...%s"`, nextExact.Value[0], s)
	}

	return s[i:], &Binding{Key: p.Name, Value: s[:i]}, nil
}

func (p FileVar) Build(w io.Writer, m map[string]string) error {
	value, ok := m[p.Name]
	if !ok {
		return fmt.Errorf(`missing key "%s"`, p.Name)
	}
	if strings.ContainsRune(value, '/') {
		return fmt.Errorf(`a file part interpolation cannot contain "/", got "%s"`, value)
	}

	fmt.Fprint(w, value)
	return nil
}

var _ Part = DirsVar{}

type DirsVar struct{ Name string }

func (p DirsVar) Match(s string, next Part) (string, *Binding, error) {
	// this capture is at the end of the pattern
	if next == nil {
		return "", &Binding{Key: p.Name, Value: s}, nil
	}

	nextExact, ok := next.(Exact)
	if !ok {
		// this cannot happen, the parser will already reject "...{a}{b}..."
		panic(`cannot have two captures one next to the other`)
	}

	i := -1
	if nextExact.Value == "/" {
		// FIX: This is not really how regex work (?)

		// DirsVar will do a greedy match from the end...
		i = strings.LastIndex(s, "/")
	} else {
		// DirsVar will do a greedy match on the full next exact value...
		i = strings.Index(s, nextExact.Value)
	}

	if i == -1 {
		return s, nil, fmt.Errorf(`expected "%v" but got "...%s"`, nextExact.Value[0], s)
	}

	return s[i:], &Binding{Key: p.Name, Value: s[:i]}, nil
}

func (p DirsVar) Build(w io.Writer, m map[string]string) error {
	value, ok := m[p.Name]
	if !ok {
		return fmt.Errorf(`missing key "%s"`, p.Name)
	}

	fmt.Fprint(w, value)
	return nil
}
