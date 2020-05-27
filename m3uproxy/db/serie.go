package db

/*
	"num": 1,
    "name": "Game of Thrones",
    "series_id": 1,
    "cover": "https:\/\/image.tmdb.org\/t\/p\/w600_and_h900_bestv2\/u3bZgnGQ9T01sWNhyveQz0wH0Hl.jpg",
    "plot": "Seven noble families fight for control of the mythical land of Westeros. Friction between the houses leads to full-scale war. All while a very ancient evil awakens in the farthest north. Amidst the war, a neglected military order of misfits, the Night's Watch, is all that stands between the realms of men and icy horrors beyond.",
    "cast": "Emilia Clarke, Kit Harington, Peter Dinklage, Lena Headey, Nikolaj Coster-Waldau",
    "director": "",
    "genre": "Sci-Fi & Fantasy, Drama",
    "releaseDate": "2011-04-17",
    "last_modified": "1560597573",
    "rating": "8",
    "rating_5based": 4,
    "backdrop_path": [
      "https:\/\/image.tmdb.org\/t\/p\/w1280\/suopoADq0k8YZr4dQXcU6pToj6s.jpg"
    ],
    "youtube_trailer": "bjqEWgDVPe0",
    "episode_run_time": "60",
    "category_id": "28"
*/
/*
type Serie struct {
	Num            int      `json:"num"`
	Name           string   `json:"name"`
	SeriesID       int      `json:"series_id"`
	Cover          string   `json:"cover"`
	Plot           string   `json:"plot"`
	Cast           string   `json:"cast"`
	Director       string   `json:"director"`
	Genre          string   `json:"genre"`
	ReleaseDate    string   `json:"releaseDate"`
	LastModified   string   `json:"last_modified"`
	Rating         string   `json:"rating"`
	Rating5based   float32  `json:"rating_5based"`
	BackdropPath   []string `json:"backdrop_path"`
	YoutubeTrailer string   `json:"youtube_trailer"`
	EpisodeRunTime string   `json:"episode_run_time"`
	CategoryID     string   `json:"category_id"`
}
*/
type Serie struct {
	Num            int      `json:"num"`
	Name           string   `json:"name"`
	SeriesID       int      `json:"series_id"`
	Cover          string   `json:"cover"`
	Plot           string   `json:"plot"`
	Cast           string   `json:"cast"`
	Director       string   `json:"director"`
	Genre          string   `json:"genre"`
	ReleaseDate    string   `json:"releaseDate"`
	LastModified   string   `json:"last_modified"`
	Rating         string   `json:"rating"`
	Rating5Based   float32  `json:"rating_5based"`
	BackdropPath   []string `json:"backdrop_path"`
	YoutubeTrailer string   `json:"youtube_trailer"`
	EpisodeRunTime string   `json:"episode_run_time"`
	CategoryID     string   `json:"category_id"`
}
