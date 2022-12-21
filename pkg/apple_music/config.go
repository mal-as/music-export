package apple_music

const (
	BaseURL = "https://amp-api.music.apple.com/v1"
	Origin  = "https://music.apple.com"
)

type Config struct {
	AuthToken      string
	MediaUserToken string
	Origin         string
	BaseURL        string
}
