package leaf

import (
	"context"
	"io/ioutil"
	"os"
	"testing"

	"github.com/gopasspw/gopass/internal/backend"
	"github.com/gopasspw/gopass/pkg/ctxutil"

	"github.com/blang/semver"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGit(t *testing.T) {
	ctx := context.Background()

	tempdir, err := ioutil.TempDir("", "gopass-")
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(tempdir)
	}()

	s, err := createSubStore(tempdir)
	require.NoError(t, err)

	assert.NotNil(t, s.RCS())
	assert.Equal(t, "noop", s.RCS().Name())
	assert.NoError(t, s.RCS().InitConfig(ctx, "foo", "bar@baz.com"))
	assert.Equal(t, semver.Version{}, s.RCS().Version(ctx))
	assert.NoError(t, s.RCS().AddRemote(ctx, "foo", "bar"))
	assert.NoError(t, s.RCS().Pull(ctx, "origin", "master"))
	assert.NoError(t, s.RCS().Push(ctx, "origin", "master"))

	assert.NoError(t, s.GitInit(ctx))
	assert.NoError(t, s.GitInit(backend.WithRCSBackend(ctx, backend.Noop)))
	assert.Error(t, s.GitInit(backend.WithRCSBackend(ctx, -1)))

	ctx = ctxutil.WithUsername(ctx, "foo")
	ctx = ctxutil.WithEmail(ctx, "foo@baz.com")
	assert.NoError(t, s.GitInit(backend.WithRCSBackend(ctx, backend.GitCLI)))
}

func TestGitRevisions(t *testing.T) {
	ctx := context.Background()

	tempdir, err := ioutil.TempDir("", "gopass-")
	require.NoError(t, err)
	defer func() {
		_ = os.RemoveAll(tempdir)
	}()

	s, err := createSubStore(tempdir)
	require.NoError(t, err)

	assert.NotNil(t, s.RCS())
	assert.Equal(t, "noop", s.RCS().Name())
	assert.NoError(t, s.RCS().InitConfig(ctx, "foo", "bar@baz.com"))

	revs, err := s.ListRevisions(ctx, "foo")
	assert.NoError(t, err)
	assert.Equal(t, 1, len(revs))

	sec, err := s.GetRevision(ctx, "foo", "bar")
	require.NoError(t, err)
	assert.Equal(t, "foo", sec.Get("password"))
}
