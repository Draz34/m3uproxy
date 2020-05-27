package db

/*
   "num": "9058",
   "name": "beIN Sports 1",
   "stream_type": "live",
   "stream_id": "9058",
   "stream_icon": "",
   "epg_channel_id": "NULL",
   "added": "1537364165",
   "category_id": "3",
   "custom_sid": "",
   "tv_archive": 0,
   "direct_source": "",
   "tv_archive_duration": 0
*/
/*
type Live struct {
	Num               string `json:"num"`
	Name              string `json:"name"`
	StreamType        string `json:"stream_type"`
	StreamID          string `json:"stream_id"`
	StreamIcon        string `json:"stream_icon"`
	EpgChannelID      string `json:"epg_channel_id"`
	Added             string `json:"added"`
	CategoryID        string `json:"category_id"`
	CustomSid         string `json:"custom_sid"`
	TvArchive         int    `json:"tv_archive"`
	DirectSource      string `json:"direct_source"`
	TvArchiveDuration int    `json:"tv_archive_duration"`
}
*/
type Live struct {
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
