package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"sync"

	"github.com/joho/godotenv"

	"github.com/mal-as/music-export/pkg/apple_music"
	"github.com/mal-as/music-export/pkg/yandex_music"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Load: %s", err)
	}
	ctx := context.Background()
	client := &http.Client{}
	amCFG := apple_music.Config{
		AuthToken:      os.Getenv("APPLE_MUSIC_AUTH_TOKEN"),
		MediaUserToken: os.Getenv("APPLE_MUSIC_MEDIA_USER_TOKEN"),
		Origin:         apple_music.Origin,
		BaseURL:        apple_music.BaseURL,
	}
	am := apple_music.NewTracksProvider(amCFG, client)

	ymCFG := yandex_music.Config{
		UserID:     os.Getenv("YANDEX_MUSIC_USER_ID"),
		PlaylistID: os.Getenv("YANDEX_MUSIC_PLAYLIST_ID"),
		AuthToken:  os.Getenv("YANDEX_MUSIC_AUTH_TOKEN"),
	}
	ym := yandex_music.NewPlayListProvider(client, ymCFG)

	tracks, err := ym.GetTracks(ctx)
	if err != nil {
		log.Fatalf("GetTracks: %s", err)
	}

	if len(tracks) == 0 {
		log.Fatal("no tracks were found")
	}

	fmt.Printf("found %d tracks\n", len(tracks))

	sema := make(chan struct{}, 10)
	var wg sync.WaitGroup

	for _, track := range tracks {
		sema <- struct{}{}
		wg.Add(1)

		go func(track *yandex_music.Track) {
			defer wg.Done()
			defer func() { <-sema }()

			addTrackToAppleMusic(ctx, am, fmt.Sprintf("%s %s", track.Track.Artists[0].Name, track.Track.Title))
		}(track)
	}

	wg.Wait()
}

func addTrackToAppleMusic(ctx context.Context, am *apple_music.TracksProvider, tackName string) {
	resp, err := am.FindByName(ctx, tackName)
	if err != nil {
		fmt.Printf("addTrackToAppleMusic.FindByName: %s\n", err)
		return
	}

	if len(resp.Results.Songs.Data) == 0 {
		fmt.Printf("addTrackToAppleMusic: song %s was not found\n", tackName)
		return
	}

	trackID := resp.Results.Songs.Data[0].ID

	fmt.Printf("addTrackToAppleMusic: song '%s' - track ID: %s\n", tackName, trackID)

	if err = am.AddToLibrary(ctx, trackID); err != nil {
		if err != nil {
			fmt.Printf("addTrackToAppleMusic.AddToLibrary: %s\n", err)
			return
		}
	}

	fmt.Printf("addTrackToAppleMusic: successfully added '%s'\n", tackName)
}
