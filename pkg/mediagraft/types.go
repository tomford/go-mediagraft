package mediagraft

import "time"

// ImageSize is a string representation of hte size of an image available
// from the image store. This should be either "original", or "WxH"
type ImageSize string

type URL string

// Images is a set of image URLs keyed by size
type Images map[ImageSize]URL

type Artist struct {
	Id          int    `json:"artistId,string"`
	Name        string `json:"artistName"`
	DisplayName string `json:"artistDisplayName"`
	Images      Images
	Streamable  bool `json:"streamable,string"`

	IsAlikeTitleMatch  bool `json:"isAlikeTitleMatch,string"`
	IsAlikeArtistMatch bool `json:"isAlikeTitleMatch,string"`

	Description      string `json:"artistDescription"`
	URL              URL
	IsVarious        bool `json:"isVarious,string"`
	CommentaryArtist bool `json:"commentaryArtist,string"`

	Genres []Genre
	Albums []Album
	Tracks []Track
}

type Track struct {
	Id         int    `json:"trackId,string"`
	Title      string `json:"trackTitle"`
	Images     Images
	Streamable bool          `json:"streamable,string"`
	duration   int           `json:"duration,string"`
	Duration   time.Duration `json:"-"`

	Purchaseable   bool `json:"purchaseable,string"`
	Radioable      bool `json:"radioable,string"`
	CopyRight      string
	OwnerId        int `json:"ownerId,string"`
	PurchasPrice   string
	TrackVersionId int  `json:"trackVersionId,string"`
	TrackNumber    int  `json:"trackNumber,string"`
	DiscNumber     int  `json:"discNumber,string"`
	Explicit       bool `json:"explicit,string"`

	Genres []Genre
	Artist
	Album

	IsAlikeTitleMatch  bool `json:"isAlikeTitleMatch,string"`
	IsAlikeArtistMatch bool `json:"isAlikeTitleMatch,string"`
}

type TrackVersion struct {
	Id         int    `json:"trackId,string"`
	Title      string `json:"trackVersionTitle"`
	Images     Images
	Streamable bool          `json:"streamable,string"`
	duration   int           `json:"duration,string"`
	Duration   time.Duration `json:"-"`

	Artist

	IsAlikeTitleMatch  bool `json:"isAlikeTitleMatch,string"`
	IsAlikeArtistMatch bool `json:"isAlikeTitleMatch,string"`
}

type Album struct {
	Id         int    `json:"albumId,string"`
	Title      string `json:"albumTitle"`
	Images     Images
	Streamable bool `json:"streamable,string"`

	Artist
	Genres []Genre
	Tracks []Track

	Composer   string
	ComposerID int `json:"composerId,string"`

	IsAlikeTitleMatch  bool `json:"isAlikeTitleMatch,string"`
	IsAlikeArtistMatch bool `json:"isAlikeTitleMatch,string"`
}

type RadioStation struct {
	Id         int    `json:"stationId,string"`
	Name       string `json:"stationName"`
	Images     Images
	Promoted   bool `json:"promoted,string"`
	Streamable bool `json:"streamable,string"`

	IsAlikeTitleMatch  bool `json:"isAlikeTitleMatch,string"`
	IsAlikeArtistMatch bool `json:"isAlikeTitleMatch,string"`
}

type Genre struct {
	Id         int    `json:"genreId,string"`
	Name       string `json:"genreName"`
	Images     Images
	Streamable bool `json:"streamable,string"`

	IsAlikeTitleMatch  bool `json:"isAlikeTitleMatch,string"`
	IsAlikeArtistMatch bool `json:"isAlikeTitleMatch,string"`
}

type Playlist struct {
	Id          int    `json:"playlistId,string"`
	Name        string `json:"playlistName"`
	Description string `json:"playlistDescription"`
	Version     int    `json:"playlistVersion,string"`

	Images Images
	Tracks []Track

	User
}

type StreamUnique string
type Stream struct {
	Id       int          `json:"trackId,string"`
	Unique   StreamUnique `json:"streamUnique"`
	Location URL          `json:"streamLocation"`
	Format   string
}

type User struct {
	UserId   int    `json:"userId,string"`
	UserName string `json:"userName"`
}
