# Git autosquash

Determines automatically on which commit to rebase when using [git-commit](https://git-scm.com/docs/git-commit) `--fixup` or `--squash` and calls `rebase --autosquash -i <hash> <command-line arguments>`.

## Long version

Git supports two ways of combining commits during a rebase, `fixup` and `squash`. The [git-commit](https://git-scm.com/docs/git-commit) command has two corresponding options, called `--fixup` and `--squash`, to indicate that the new commit should be squashed or fixed up with the specified commit.

The [git-rebase](https://git-scm.com/docs/git-rebase) command supports `--autosquash`, which automatically modifies the todo list of the rebase so that commits marked for squashing or fixup comes right after the commit to be modified. It also changes the action from `pick` to `fixup` or `squash`. 

Instead of calling `rebase` manually with the correct target commit, `autosquash` automatically determines on which commit to rebase and calls git rebase.

## Installation

* Install [Go](https://golang.org/doc/install).
* Set the `GOBIN` environment variable to where `go install` will install a command.
* Install `autosquash` by running `go install autosquash.go`.

## Usage

Run `autosquash` in any git repository that has commits to fixup. Any command line arguments passed to `autosquash` are passed to `git rebase`.

## License

Autosquash is MIT licensed, see the [LICENSE](LICENSE) file.
