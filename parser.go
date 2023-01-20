package patbu

import (
	"fmt"
	"io"
	"unicode/utf8"
)

func nextRune(ctx *string) (rune, error) {
	r, size := utf8.DecodeRuneInString(*ctx)
	if r == utf8.RuneError && size == 0 {
		return 0, io.EOF
	}
	if r == utf8.RuneError && size == 1 {
		return 0, fmt.Errorf(`invalid rune at "%s"`, *ctx)
	}

	*ctx = (*ctx)[size:]

	return r, nil
}

func Parse(template string) (Patbu, error) {
	full := template
	ctx := &template

	parts := Patbu{}

	exactPart := []rune{}

	for len(*ctx) > 0 {
		r, err := nextRune(ctx)
		if err != nil {
			return nil, err
		}

		if r == '{' {
			if len(exactPart) > 0 {
				parts = append(parts, Exact{string(exactPart)})
				exactPart = exactPart[:0]
			} else if len(parts) > 0 {
				switch parts[len(parts)-1].(type) {
				case FileVar, DirsVar:
					return nil, fmt.Errorf(`two captures cannot be next to each other in "%s"`, full)
				}
			}

			capture, err := parseCapture(ctx)
			if err != nil {
				return nil, err
			}

			parts = append(parts, capture)
		} else {
			exactPart = append(exactPart, r)
		}
	}

	if len(exactPart) > 0 {
		parts = append(parts, Exact{string(exactPart)})
	}

	return parts, nil
}

// parseCapture will parse "...foo}rest..." and leave with "rest..."
func parseCapture(ctx *string) (Part, error) {
	isDirs := false
	r, err := nextRune(ctx)
	if err != nil {
		return nil, err
	}

	name := ""

	if r == '*' {
		isDirs = true
	} else {
		name += string(r)
	}

	for {
		r, err = nextRune(ctx)
		if err != nil {
			return nil, err
		}
		if r == '}' {
			break
		}
		name += string(r)
	}

	if isDirs {
		return DirsVar{name}, nil
	} else {
		return FileVar{name}, nil
	}
}
