# Patbu

Generalized path glob **pat**terns that are also **bu**ilders. The main idea is that all of the following expressions

<p align="center">
<code>a/b/c.txt</code> ~> the literal file <code>a/b/c.txt</code>
<br>
<code>{year}-{month}-{day}.txt</code> ~> for example <code>2023-01-01.txt</code>
<br>
<code>a/{var1}/c.txt</code> ~> this won't match <code>a/b/b/c.txt</code> because the capture <code>{var1}</code> doesn't match <code>/</code>
<br>
<code>{*module}/package.json</code> ~> the <code>*</code> modifier will also match <code>/</code>
</p>

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