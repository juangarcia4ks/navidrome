package scanner

import (
	"context"

	"github.com/navidrome/navidrome/core"

	"github.com/navidrome/navidrome/conf"
	"github.com/navidrome/navidrome/consts"
	"github.com/navidrome/navidrome/model"
	"github.com/navidrome/navidrome/tests"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("playlistImporter", func() {
	var ds model.DataStore
	var ps *playlistImporter
	var pls core.Playlists
	ctx := context.Background()

	BeforeEach(func() {
		ds = &tests.MockDataStore{
			MockedMediaFile: &mockedMediaFile{},
			MockedPlaylist:  &mockedPlaylist{},
		}
		pls = core.NewPlaylists(ds)
	})

	Describe("processPlaylists", func() {
		Context("Default PlaylistsPath", func() {
			BeforeEach(func() {
				conf.Server.PlaylistsPath = consts.DefaultPlaylistsPath
			})
			It("finds and import playlists at the top level", func() {
				ps = newPlaylistImporter(ds, pls, "tests/fixtures/playlists/subfolder1")
				Expect(ps.processPlaylists(ctx, "tests/fixtures/playlists/subfolder1")).To(Equal(int64(1)))
			})

			It("finds and import playlists at any subfolder level", func() {
				ps = newPlaylistImporter(ds, pls, "tests")
				Expect(ps.processPlaylists(ctx, "tests/fixtures/playlists/subfolder1")).To(Equal(int64(1)))
			})
		})

		It("ignores playlists not in the PlaylistsPath", func() {
			conf.Server.PlaylistsPath = "subfolder1"
			ps = newPlaylistImporter(ds, pls, "tests/fixtures/playlists")

			Expect(ps.processPlaylists(ctx, "tests/fixtures/playlists/subfolder1")).To(Equal(int64(1)))
			Expect(ps.processPlaylists(ctx, "tests/fixtures/playlists/subfolder2")).To(Equal(int64(0)))
		})

		It("only imports playlists from the root of MusicFolder if PlaylistsPath is '.'", func() {
			conf.Server.PlaylistsPath = "."
			ps = newPlaylistImporter(ds, pls, "tests/fixtures/playlists")

			Expect(ps.processPlaylists(ctx, "tests/fixtures/playlists")).To(Equal(int64(3)))
			Expect(ps.processPlaylists(ctx, "tests/fixtures/playlists/subfolder1")).To(Equal(int64(0)))
		})

	})
})

type mockedMediaFile struct {
	model.MediaFileRepository
}

func (r *mockedMediaFile) FindByPath(s string) (*model.MediaFile, error) {
	return &model.MediaFile{
		ID:   "123",
		Path: s,
	}, nil
}

type mockedPlaylist struct {
	model.PlaylistRepository
}

func (r *mockedPlaylist) FindByPath(path string) (*model.Playlist, error) {
	return nil, model.ErrNotFound
}

func (r *mockedPlaylist) Put(pls *model.Playlist) error {
	return nil
}
