# Patbu

Generalized path glob **pat**terns that are also **bu**ilders. The main idea is that all of the following expressions

- `a/b/c.txt` -- the literal file `a/b/c.txt`

- `{year}-{month}-{day}.txt` -- for example `2023-01-01.txt`

- `a/{var1}/c.txt` -- this won't match `a/b/b/c.txt` because the capture `{var1}` doesn't match `/`

- `{*module}/package.json` -- the `*` modifier will also match `/`

can be used to match _against a given pattern_ and return an optional dictionary of captures but _also as a path builder_ that is given a dictionary of bindings. This duality can be used to modify paths. For example using the provided CLI for this library we can do

```bash
echo 'src/projects/go/patbu.html' 
| patbu --stdin 'src/{*route}/{page}.html' 'dist/{*route}/{page}/index.html'
```

## Setup & Install

```bash shell
$ go build -v -o ./bin/patbu ./cmd/patbu
$ go install -v ./cmd/patbu
```

## Notes

Other name ideas: blueprint, pathbp, template, patt