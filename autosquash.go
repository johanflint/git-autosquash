package main

import (
	"fmt"

	"gopkg.in/src-d/go-git.v4"
)

func main() {
	_, err := git.PlainOpenWithOptions(".", &git.PlainOpenOptions{DetectDotGit: true})
	if err == git.ErrRepositoryNotExists {
		printf("No git repository found")
		return
	} else if err != nil {
		printf("Unable to open git repository: %s", err)
		return
	}

	printf("Repository found")
}

func printf(format string, args ...interface{}) {
	fmt.Printf("%s\n", fmt.Sprintf(format, args...))
}
