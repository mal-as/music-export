package apple_music

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"path"

	"github.com/mal-as/music-export/pkg/http/utils"
)

const (
	findTrackPath string = "catalog/ru/search"
	getTrackPath  string = "catalog/ru/songs"
	addTrackPath  string = "me/library"
)

type (
	TracksProvider struct {
		httpClient *http.Client
		cfg        Config
	}

	FindTrackResponse struct {
		Results *FindTrackResults `json:"results"`
	}
	FindTrackResults struct {
		Songs *Songs `json:"songs"`
	}
	Songs struct {
		Data []*Data `json:"data"`
	}
	Data struct {
		ID string `json:"id"`
	}

	GetTrackResponse struct {
		Data []*TrackData `json:"data"`
	}
	TrackData struct {
		Attributes *TrackAttributes `json:"attributes"`
	}
	TrackAttributes struct {
		Name       string `json:"name"`
		ArtistName string `json:"artistName"`
	}
)

func NewTracksProvider(cfg Config, httpClient *http.Client) *TracksProvider {
	return &TracksProvider{cfg: cfg, httpClient: httpClient}
}

func (tp *TracksProvider) FindByName(ctx context.Context, name string) (*FindTrackResponse, error) {
	u, err := url.Parse(tp.cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("FindByName.Parse: %w", err)
	}

	u.Path = path.Join(u.Path, findTrackPath)

	q := u.Query()
	q.Add("limit", "1")
	q.Add("types", "songs")
	q.Add("term", name)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("FindByName.NewRequestWithContext: %w", err)
	}

	req.Header.Set("authorization", tp.cfg.AuthToken)
	req.Header.Set("origin", tp.cfg.Origin)
	req.URL.RawQuery = q.Encode()

	resp, err := tp.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("FindByName.Do: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if err = utils.CheckHTTPErr(resp); err != nil {
		return nil, fmt.Errorf("FindByName.Do: %s", err)
	}

	result, err := utils.ParseHTTPResponse[FindTrackResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("FindByName.parseHTTPResponse: %w", err)
	}

	return result, nil
}

func (tp *TracksProvider) GetByID(ctx context.Context, id string) (*GetTrackResponse, error) {
	u, err := url.Parse(tp.cfg.BaseURL)
	if err != nil {
		return nil, fmt.Errorf("GetByID.Parse: %w", err)
	}

	u.Path = path.Join(u.Path, getTrackPath, id)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("GetByID.NewRequestWithContext: %w", err)
	}

	req.Header.Set("authorization", tp.cfg.AuthToken)
	req.Header.Set("origin", tp.cfg.Origin)

	resp, err := tp.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("GetByID.Do: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if err = utils.CheckHTTPErr(resp); err != nil {
		return nil, fmt.Errorf("GetByID.Do: %s", err)
	}

	result, err := utils.ParseHTTPResponse[GetTrackResponse](resp)
	if err != nil {
		return nil, fmt.Errorf("GetByID.parseHTTPResponse: %w", err)
	}

	return result, nil
}

func (tp *TracksProvider) AddToLibrary(ctx context.Context, id string) error {
	u, err := url.Parse(tp.cfg.BaseURL)
	if err != nil {
		return fmt.Errorf("AddToLibrary.Parse: %w", err)
	}

	u.Path = path.Join(u.Path, addTrackPath)

	q := u.Query()
	q.Add("art[url]", "f")
	q.Add("ids[songs]", id)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u.String(), nil)
	if err != nil {
		return fmt.Errorf("AddToLibrary.NewRequestWithContext: %w", err)
	}

	req.Header.Set("authorization", tp.cfg.AuthToken)
	req.Header.Set("origin", tp.cfg.Origin)
	req.Header.Set("media-user-token", tp.cfg.MediaUserToken)
	req.URL.RawQuery = q.Encode()

	resp, err := tp.httpClient.Do(req)
	if err != nil {
		return fmt.Errorf("AddToLibrary.Do: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	if err = utils.CheckHTTPErr(resp); err != nil {
		return fmt.Errorf("AddToLibrary.Do: %s", err)
	}

	return nil
}
