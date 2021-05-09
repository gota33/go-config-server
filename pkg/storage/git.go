package storage

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"sync"
	"time"

	"github.com/go-git/go-billy/v5"
	"github.com/go-git/go-billy/v5/memfs"
	. "github.com/go-git/go-git/v5"
	"github.com/go-git/go-git/v5/config"
	"github.com/go-git/go-git/v5/plumbing"
	"github.com/go-git/go-git/v5/plumbing/storer"
	"github.com/go-git/go-git/v5/plumbing/transport"
	"github.com/go-git/go-git/v5/storage"
	"github.com/go-git/go-git/v5/storage/memory"
)

const (
	DefaultFetchTTL = 10 * time.Second
)

type info struct {
	remote plumbing.Hash
	local  plumbing.Hash
	fs     internalFs
}

func (i info) SameHash() bool {
	return i.local == i.remote
}

type Git struct {
	URL      string
	Auth     transport.AuthMethod
	FetchTTL time.Duration

	store    storage.Storer
	remote   *Remote
	lock     sync.Locker
	syncTime time.Time
	infos    map[plumbing.ReferenceName]info
}

func NewGit(URL string) *Git {
	store := memory.NewStorage()
	remote := NewRemote(store, &config.RemoteConfig{
		Name:  "origin",
		URLs:  []string{URL},
		Fetch: []config.RefSpec{"refs/heads/*:refs/heads/*"},
	})

	return &Git{
		URL:      URL,
		FetchTTL: DefaultFetchTTL,
		lock:     &sync.Mutex{},
		store:    store,
		remote:   remote,
		infos:    make(map[plumbing.ReferenceName]info),
	}
}

func (g *Git) Provide(ctx context.Context, namespace string) (fs ReadonlyFs, err error) {
	if !g.skipFetch() {
		if err = g.fetch(ctx); err != nil {
			return
		}
	}

	branch := g.branchRef(namespace)
	if _, ok := g.infos[branch]; !ok {
		err = fmt.Errorf("namespace not found: %q", namespace)
		return
	}

	if !g.skipCheckout(branch) {
		if fs, err = g.checkout(branch); err != nil {
			return
		}
	}
	return g.infos[branch].fs, nil
}

func (g *Git) fetch(ctx context.Context) (err error) {
	g.lock.Lock()
	defer g.lock.Unlock()

	if err = g.remote.FetchContext(ctx, &FetchOptions{
		RefSpecs: g.remote.Config().Fetch,
		Auth:     g.Auth,
		Force:    true,
		Depth:    1,
		Progress: os.Stdout,
	}); errors.Is(err, NoErrAlreadyUpToDate) ||
		errors.Is(err, transport.ErrEmptyUploadPackRequest) {
		err = nil
	} else if err != nil {
		return
	}

	ref := plumbing.NewSymbolicReference(plumbing.HEAD, plumbing.Master)
	if err = g.store.SetReference(ref); err != nil {
		return
	}
	if err = g.updateRefs(); err != nil {
		return
	}

	g.syncTime = time.Now()
	log.Println("Fetch: OK")
	return
}

func (g *Git) updateRefs() (err error) {
	var (
		iter storer.ReferenceIter
		refs = make(map[plumbing.ReferenceName]plumbing.Hash)
		each = func(ref *plumbing.Reference) (_ error) { refs[ref.Name()] = ref.Hash(); return }
	)
	if iter, err = g.store.IterReferences(); err != nil {
		return
	}
	if err = iter.ForEach(each); err != nil {
		return
	}

	next := make(map[plumbing.ReferenceName]info, len(refs))
	for name, hash := range refs {
		old := g.infos[name]
		next[name] = info{
			remote: hash,
			local:  old.local,
			fs:     old.fs,
		}
	}
	g.infos = next
	return
}

func (g *Git) checkout(branch plumbing.ReferenceName) (fs internalFs, err error) {
	g.lock.Lock()
	defer g.lock.Unlock()

	var (
		repo *Repository
		wt   *Worktree
		head *plumbing.Reference
	)
	fs.Filesystem = memfs.New()

	if repo, err = Open(g.store, fs.Filesystem); err != nil {
		return
	}
	if wt, err = repo.Worktree(); err != nil {
		return
	}

	if err = wt.Checkout(&CheckoutOptions{
		Branch: branch,
		Force:  true,
	}); err != nil {
		return
	}

	if head, err = repo.Head(); err != nil {
		return
	}

	old := g.infos[branch]
	g.infos[branch] = info{
		remote: old.remote,
		local:  head.Hash(),
		fs:     fs,
	}

	log.Println("Checkout: OK")
	return
}

func (g *Git) skipFetch() bool                                 { return time.Since(g.syncTime) < g.FetchTTL }
func (g *Git) skipCheckout(branch plumbing.ReferenceName) bool { return g.infos[branch].SameHash() }

func (g *Git) branchRef(namespace string) plumbing.ReferenceName {
	return plumbing.ReferenceName(fmt.Sprintf("refs/heads/%s", namespace))
}

type internalFs struct{ billy.Filesystem }

func (fs internalFs) Open(name string) (ReadonlyFile, error) { return fs.Filesystem.Open(name) }
func (fs internalFs) Close() (_ error)                       { return }
