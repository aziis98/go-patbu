# Patbu

Generalized path glob **pat**terns that are also **bu**ilders. The main idea is that an expression like

<p align="center">
<code>a/b/{x}.txt</code>
</p>

can be interpreted as a pattern or a template at the same time, for example the following expressions have the following intuitive meaning

- `a/b/c.txt` corresponds to the literal file `a/b/c.txt`

- `{year}-{month}-{day}.txt` will match something like `2023-01-20.txt`

- `a/{var1}/c.txt` won't match `a/b/b/c.txt` because the capture `{var1}` doesn't match `/`, the equivalent is that this throws an error when used as a template and `var1` contains slashes.

- `{*module}/package.json`, here the `*` modifier will also match `/`.

so this syntax can be used to

- match _against a given pattern_ and return a dictionary of captures

- but _also as a path template or builder_  that is given a dictionary of bindings

This duality can be used to modify paths for example. Using the provided CLI for this library we can do something like the following

```bash
echo 'blog/src/projects/go/patbu.html' 
| patbu --stdin 'blog/src/{*route}/{page}.html' 'blog/dist/{*route}/{page}/index.html'
```

## Setup & Install

```bash shell
$ go build -v -o ./bin/patbu ./cmd/patbu
$ go install -v ./cmd/patbu
```

## Notes & Todos

Other name ideas: blueprint, pathbp, template, patt

Add some flags to use as a bulk rename utility
