package db

/*
	"num": "26286",
    "name": "The Kill Team",
    "stream_type": "movie",
    "stream_id": "26286",
    "stream_icon": "http:\/\/covers.suptv.net\/form\/covers\/The-Kill-Team_20200413010042.jpg",
    "rating": "6.6",
    "rating_5based": 3.3,
    "added": "1543790367",
    "category_id": "101",
    "container_extension": "mp4",
    "custom_sid": "",
    "direct_source": ""
*/
/*
type Movie struct {
	ID                 string  `json:"num"`
	Name               string  `json:"name"`
	StreamType         string  `json:"stream_type"`
	StreamID           string  `json:"stream_id"`
	StreamIcon         string  `json:"stream_icon"`
	Rating             string  `json:"rating"`
	Rating5based       float32 `json:"rating_5based"`
	Added              string  `json:"added"`
	CategoryID         string  `json:"category_id"`
	ContainerExtension string  `json:"container_extension"`
	CustomSid          string  `json:"custom_sid"`
	DirectSource       string  `json:"direct_source"`
}
*/
type Movie struct {
	Num                int     `json:"num"`
	Name               string  `json:"name"`
	StreamType         string  `json:"stream_type"`
	StreamID           int     `json:"stream_id"`
	StreamIcon         string  `json:"stream_icon"`
	Rating             string  `json:"rating"`
	Rating5Based       float64 `json:"rating_5based"`
	Added              string  `json:"added"`
	CategoryID         string  `json:"category_id"`
	ContainerExtension string  `json:"container_extension"`
	CustomSid          string  `json:"custom_sid"`
	DirectSource       string  `json:"direct_source"`
}

/*
func (movies []Movie) Sort() []Movie {

	return movies
}
*/
