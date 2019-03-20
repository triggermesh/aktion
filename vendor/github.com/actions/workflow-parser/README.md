[actions]: https://github.com/features/actions/
[workflow]: https://developer.github.com/actions/creating-workflows/creating-a-new-workflow/
# The GitHub Actions Workflow Parser

This is the [language specification](language.md) and the official parser
for GitHub Actions [`main.workflow` files][workflow].  It is running in
production as part of [GitHub Actions][actions].  It is written in Go.

There are syntax-highlighting configuration files for Vim and Atom, under
the `syntax/` directory.

## Using the parser

To use the parser in your own projects, import it, and call the `Parse`
function.

```go
import "github.com/actions/workflow-parser/parser"
...
config, err := parser.Parse(reader)
```

By default, the `Parse` function validates basic syntax, type safety, and
all dependencies within a `.workflow` file.  It returns a model with
arrays of all workflows and actions defined in the file.

If there are any errors, `Parse` returns an error.  System errors are
returned as a generic `error` class, while problems in the file are
returned as a `parser.Error`.  The `parser.Error` struct has an array of
errors, each indicating a severity and a position in the file.

Warnings indicate code that might get ignored or misinterpreted.  Errors
indicate code that is incomplete or has type errors and cannot run.  Fatal
errors indicate that the file cannot be even partially displayed, due to a
syntax error or circular dependency.  Only `.workflow` files with no
warnings, errors, or fatal errors will work with Actions.

To suppress warnings or non-fatal errors, use either of the following
functions as an optional second argument to `Parse`:

```go
config, err := parser.Parse(reader, parser.WithSuppressWarnings())
// or
config, err := parser.Parse(reader, parser.WithSuppressErrors())
```

## Developing the parser

You'll need a copy of go v1.9 or higher.  You might also want a copy of
`dep`, if you plan to change `Gopkg.toml`.

On OS X, `brew install go dep` will get you there.

You can also use `docker run -it golang` on any Docker-compatible platform
to get a copy of Go.  That doesn't include `dep`.

To run the tests and build a command-line binary that validates workflow
files, run `make`.  The resulting validation binary works like so:

```
$ ./cmd/parser samples/a.workflow 
samples/a.workflow is a valid file with 9 actions and 1 workflow
```

If you would like to contribute your work back to the project, please see
[`CONTRIBUTING.md`](CONTRIBUTING.md).


## License

This project is open source, under the [MIT license](LICENSE).


## Releases

Active development happens on the `master` branch.  Releases happen by
tagging commits with `vX.Y.Z`, according to [semantic
versioning](https://semver.org/).  Increment the major version number
for:
 - Changes to the parser implementation that will break code application code.
 - Changes to the language spec that will reject existing, previously valid `.workflow` files.  If possible, use the top-level `version=` keyword to provide backward compatibility for existing files.

Increment the minor version number for:
 - New language features that prior versions of the parser would have rejected.
