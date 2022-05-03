package sqlstore

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/grafana/grafana/pkg/models"
)

func TestIntegrationPlaylistDataAccess(t *testing.T) {
	ss := InitTestDB(t)

	t.Run("Can create playlist", func(t *testing.T) {
		items := []models.PlaylistItemDTO{
			{Title: "graphite", Value: "graphite", Type: "dashboard_by_tag"},
			{Title: "Backend response times", Value: "3", Type: "dashboard_by_id"},
		}
		cmd := models.CreatePlaylistCommand{Name: "NYC office", Interval: "10m", OrgId: 1, Items: items}
		err := ss.CreatePlaylist(context.Background(), &cmd)
		require.NoError(t, err)

		uid := cmd.Result.Uid
		id := cmd.Result.Id
		fmt.Printf("uid: %s, id: %d\n", uid, id)

		t.Run("Can update playlist", func(t *testing.T) {
			items := []models.PlaylistItemDTO{
				{Title: "influxdb", Value: "influxdb", Type: "dashboard_by_tag"},
				{Title: "Backend response times", Value: "2", Type: "dashboard_by_id"},
			}
			query := models.UpdatePlaylistCommand{Name: "NYC office ", OrgId: 1, Uid: uid, Interval: "10s", Items: items}
			err = ss.UpdatePlaylist(context.Background(), &query)
			require.NoError(t, err)
		})

		t.Run("Can remove playlist", func(t *testing.T) {
			deleteQuery := models.DeletePlaylistCommand{Uid: uid, OrgId: 1}
			err = ss.DeletePlaylist(context.Background(), &deleteQuery)
			require.NoError(t, err)

			getQuery := models.GetPlaylistByUidQuery{Uid: uid}
			err = ss.GetPlaylist(context.Background(), &getQuery)
			require.NoError(t, err)
			require.Equal(t, uid, getQuery.Result.Uid, "playlist should've been removed")
		})
	})

	t.Run("Delete playlist that doesn't exist", func(t *testing.T) {
		deleteQuery := models.DeletePlaylistCommand{Uid: "654312", OrgId: 1}
		err := ss.DeletePlaylist(context.Background(), &deleteQuery)
		require.NoError(t, err)
	})

	t.Run("Delete playlist with invalid command yields error", func(t *testing.T) {
		testCases := []struct {
			desc string
			cmd  models.DeletePlaylistCommand
		}{
			{desc: "none", cmd: models.DeletePlaylistCommand{}},
			{desc: "no OrgId", cmd: models.DeletePlaylistCommand{Uid: "1"}},
			{desc: "no Uid", cmd: models.DeletePlaylistCommand{OrgId: 1}},
		}

		for _, tc := range testCases {
			t.Run(tc.desc, func(t *testing.T) {
				err := ss.DeletePlaylist(context.Background(), &tc.cmd)
				require.EqualError(t, err, models.ErrCommandValidationFailed.Error())
			})
		}
	})
}
