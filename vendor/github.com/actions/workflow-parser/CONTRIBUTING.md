## Contributing

[fork]: https://github.com/actions/workflow-parser/fork
[pr]: https://github.com/github/workflow-parser/compare
[code-of-conduct]: CODE_OF_CONDUCT.md
[effective-go]: https://golang.org/doc/effective_go.html
[actions]: https://developer.github.com/actions/

Hi there! We're thrilled that you'd like to contribute to this project. Your help is essential for keeping it great.

Contributions to this project are [released](https://help.github.com/articles/github-terms-of-service/#6-contributions-under-repository-license) to the public under the [MIT license](LICENSE.md).

Please note that this project is released with a [Contributor Code of Conduct][code-of-conduct]. By participating in this project you agree to abide by its terms.

## Submitting a pull request

0. [Fork][fork] and clone the repository
0. Make sure the build and tests succeed on your machine: `make`
0. Create a new branch: `git checkout -b my-branch-name`
0. Make your change, add tests, and make sure the tests still pass
0. Push to your fork and [submit a pull request][pr]
0. Pat your self on the back and wait for your pull request to be reviewed and merged.

Here are a few things you can do that will increase the likelihood of your pull request being accepted:

- Write [effective Go][effective-go].
- Run `make fmt` to format your code.
- Write tests.
- Keep your change as focused as possible. If there are multiple changes you would like to make that are not dependent upon each other, consider submitting them as separate pull requests.
- Write a [good commit message](http://tbaggery.com/2008/04/19/a-note-about-git-commit-messages.html).

## Notes

We are eager to make the parser and the Actions workflow language useful to as many people in the Actions community as possible. To that end, we particularly welcome improvements to the implementation, versions of the parser in additional languages, and syntax-highlighting files for more editors.

However, this parser is only a small part of the overall Actions project. Pull requests that add features to the language are less likely to be merged, because new features only make sense if we can implement them throughout the Actions platform. If you want to talk about feature-level work, please reach out through the feedback channels [here][actions].

## Resources

- [How to Contribute to Open Source](https://opensource.guide/how-to-contribute/)
- [Using Pull Requests](https://help.github.com/articles/about-pull-requests/)
- [GitHub Help](https://help.github.com)

