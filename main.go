package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/mal-as/music-export/pkg/apple_music"

	"github.com/joho/godotenv"
)

func main() {
	var name string
	flag.StringVar(&name, "name", "", "name of the track")
	flag.Parse()

	if name == "" {
		log.Fatal("no track name was specified")
	}

	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Load: %s", err)
	}

	ctx := context.Background()
	client := &http.Client{}
	cfg := apple_music.Config{
		AuthToken:      os.Getenv("APPLE_MUSIC_AUTH_TOKEN"),
		MediaUserToken: os.Getenv("APPLE_MUSIC_MEDIA_USER_TOKEN"),
		Origin:         apple_music.Origin,
		BaseURL:        apple_music.BaseURL,
	}

	tp := apple_music.NewTracksProvider(cfg, client)

	resp, err := tp.FindByName(ctx, name)
	if err != nil {
		log.Fatalf("FindByName: %s", err)
	}

	if len(resp.Results.Songs.Data) == 0 {
		log.Fatal("no song was find")
	}

	trackID := resp.Results.Songs.Data[0].ID

	fmt.Println("track ID: ", trackID)

	track, err := tp.GetByID(ctx, trackID)
	if err != nil {
		log.Fatalf("GetTrack: %s", err)
	}

	if len(track.Data) == 0 {
		log.Fatal("no song was got")
	}

	fmt.Printf("%s - %s\n", track.Data[0].Attributes.ArtistName, track.Data[0].Attributes.Name)

	if err = tp.AddToLibrary(ctx, trackID); err != nil {
		if err != nil {
			log.Fatalf("AddToLibrary: %s", err)
		}
	}

	fmt.Println("success")
}
