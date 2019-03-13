# Git autosquash

Determines automatically on which commit to rebase when using [git-commit](https://git-scm.com/docs/git-commit) `--fixup` or `--squash` and calls `rebase --autosquash -i <hash> <command-line arguments>`.
