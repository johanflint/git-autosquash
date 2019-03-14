package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"gopkg.in/src-d/go-git.v4"
	"gopkg.in/src-d/go-git.v4/plumbing/object"
	"gopkg.in/src-d/go-git.v4/plumbing/storer"
)

const fixupPrefix = "fixup!"
const squashPrefix = "squash!"

func main() {
	repository, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	if err == git.ErrRepositoryNotExists {
		printf("No git repository found")
		return
	} else if err != nil {
		printf("Unable to open git repository: %s", err)
		return
	}

	var upstreamCommit *object.Commit
	autosquashCommitMessages := make([]string, 0)

	log, err := repository.Log(&git.LogOptions{})
	log.ForEach(func(commit *object.Commit) error {
		if commit.NumParents() > 1 {
			return storer.ErrStop
		}

		if matchesAutosquashCommit(commit, autosquashCommitMessages) {
			upstreamCommit = commit
		}

		if strings.HasPrefix(commit.Message, fixupPrefix) || strings.HasPrefix(commit.Message, squashPrefix) {
			messageWithoutPrefix := commitTitle(trimPrefix(commit.Message))
			autosquashCommitMessages = append(autosquashCommitMessages, messageWithoutPrefix)
		}

		return nil
	})

	if len(autosquashCommitMessages) == 0 {
		printf("Nothing to autosquash: no commits found that start with %s or %s", fixupPrefix, squashPrefix)
		return
	}

	upstreamParentCommit, err := upstreamCommit.Parent(0)
	if err == object.ErrParentNotFound {
		printf("Unable to find the parent commit of %s", upstreamCommit.Hash)
		return
	}

	gitArguments := []string{"rebase", "--autosquash", "-i", fmt.Sprintf("%s", upstreamParentCommit.Hash)}
	if len(os.Args) > 1 {
		gitArguments = append(gitArguments, os.Args[1:]...)
	}

	printf("Executing git %s", strings.Join(gitArguments, " "))
	command := exec.Command("git", gitArguments...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Run()
}

func printf(format string, args ...interface{}) {
	fmt.Printf("%s\n", fmt.Sprintf(format, args...))
}

func commitTitle(commitMessage string) string {
	return strings.Split(strings.ReplaceAll(commitMessage, "\r\n", "\n"), "\n")[0]
}

func trimPrefix(commitMessage string) string {
	trimmed := strings.TrimPrefix(commitMessage, fixupPrefix)
	trimmed = strings.TrimPrefix(trimmed, squashPrefix)

	return strings.TrimLeft(trimmed, " ")
}

func matchesAutosquashCommit(commit *object.Commit, autosquashCommitMessages []string) bool {
	for _, commitMessage := range autosquashCommitMessages {
		if strings.Contains(commit.Message, commitMessage) {
			return true
		}
	}

	return false
}
