package storage

import (
	"context"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	. "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage/memory"
)

type Git struct {
	URL  string
	Auth transport.AuthMethod

	namespace string
	repo      *Repository
}

func (g *Git) Use(ctx context.Context, namespace string) (err error) {
	if g.namespace == "" {
		if g.repo, err = g.clone(ctx, namespace); err == nil {
			g.namespace = namespace
		}
		return
	}
	return
}

func (g *Git) Read(path string) (content string, err error) {
	var (
		wt   *Worktree
		fd   billy.File
		data []byte
	)
	if wt, err = g.repo.Worktree(); err != nil {
		return
	}
	if fd, err = wt.Filesystem.Open(path); err != nil {
		return
	}
	if data, err = ioutil.ReadAll(fd); err != nil {
		return
	}
	return string(data), nil
}

func (g *Git) clone(ctx context.Context, namespace string) (repo *Repository, err error) {
	opts := &CloneOptions{
		URL:          g.URL,
		Auth:         g.Auth,
		SingleBranch: true,
		Depth:        1,
		Progress:     os.Stdout,
	}
	if opts.ReferenceName, err = g.branchRef(namespace); err != nil {
		return
	}

	return CloneContext(ctx, memory.NewStorage(), memfs.New(), opts)
}

func (g *Git) branchRef(namespace string) (ref plumbing.ReferenceName, err error) {
	if namespace == "" {
		err = fmt.Errorf("namespace is required")
		return
	}

	ref = plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", namespace))
	return
}
