package yandex_music

import (
	"context"
	"fmt"
	"net/http"

	"github.com/mal-as/music-export/pkg/http/utils"
)

type (
	PlayListProvider struct {
		client *http.Client
		cfg    Config
	}

	PlayListResult struct {
		Tracks []*Track `json:"tracks"`
	}
	PlayListResponse struct {
		Result *PlayListResult `json:"result"`
	}
	Track struct {
		Track *TrackData `json:"track"`
	}
	TrackData struct {
		ID      string    `json:"id"`
		Title   string    `json:"title"`
		Artists []*Artist `json:"artists"`
	}
	Artist struct {
		Name string `json:"name"`
	}
)

func NewPlayListProvider(client *http.Client, cfg Config) *PlayListProvider {
	return &PlayListProvider{client: client, cfg: cfg}
}

func (p *PlayListProvider) GetTracks(ctx context.Context) ([]*Track, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, fmt.Sprintf(
		"https://api.music.yandex.net/users/%s/playlists/%s",
		p.cfg.UserID,
		p.cfg.PlaylistID),
		nil)
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", p.cfg.AuthToken)

	resp, err := p.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if err = utils.CheckHTTPErr(resp); err != nil {
		return nil, fmt.Errorf("GetTracks.Do: %s", err)
	}

	result, err := utils.ParseHTTPResponse[PlayListResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("GetTracks.parseHTTPResponse: %w", err)
	}

	return result.Result.Tracks, nil
}
