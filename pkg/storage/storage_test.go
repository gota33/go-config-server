package storage

import (
	"context"
	"io/ioutil"
	"testing"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	"github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/transport/http"
	"github.com/go-git/go-git/v5/storage/memory"
	"github.com/stretchr/testify/assert"
)

const repo = `https://github.com/gota33/go-config-server.git`

func TestGit(t *testing.T) {
	ctx := context.TODO()

	g := NewGit(repo)

	fs, err := g.Provide(ctx, "master")
	if !assert.NoError(t, err) {
		return
	}

	if fd, err := fs.Open("go.mod"); assert.NoError(t, err) {
		assert.NotNil(t, fd)
	}
}

func TestB(t *testing.T) {
	store := memory.NewStorage()
	auth := &http.BasicAuth{}

	repo0, err := git.Clone(store, nil,
		&git.CloneOptions{URL: repo, Auth: auth})
	assert.NoError(t, err)

	err = repo0.Fetch(&git.FetchOptions{
		RefSpecs: []config.RefSpec{"refs/*:refs/*"},
		Auth:     auth,
		Force:    true,
	})
	assert.NoError(t, err)

	ShowBranches(t, repo0)

	fs0 := memfs.New()
	repo, err := git.Open(store, fs0)
	assert.NoError(t, err)

	wt, err := repo.Worktree()
	assert.NoError(t, err)

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: "refs/heads/master",
		Force:  true,
	})
	assert.NoError(t, err)

	fs1 := memfs.New()
	repo, err = git.Open(store, fs1)
	assert.NoError(t, err)

	wt, err = repo.Worktree()
	assert.NoError(t, err)

	err = wt.Checkout(&git.CheckoutOptions{
		Branch: "refs/heads/example",
		Force:  true,
	})
	assert.NoError(t, err)

	Cat(t, fs0, "README.md")
	Cat(t, fs1, "arith.jsonnet")
}

func Cat(t *testing.T, fs billy.Filesystem, name string) {
	t.Logf("Cat %q", name)
	fd, err := fs.Open(name)
	assert.NoError(t, err)

	data, err := ioutil.ReadAll(fd)
	assert.NoError(t, err)

	t.Log(string(data))
}

func ShowBranches(t *testing.T, repo *git.Repository) {
	iter, err := repo.Branches()
	assert.NoError(t, err)

	err = iter.ForEach(func(reference *plumbing.Reference) error {
		t.Logf("ref: %s", reference.Name())
		return nil
	})
	assert.NoError(t, err)
}
