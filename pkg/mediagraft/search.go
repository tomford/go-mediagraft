package mediagraft

type SearchResult struct {
	Artists       []Artist
	Albums        []Album
	Tracks        []Track
	TrackVersions []TrackVersion
	Genres        []Genre
	RadioStations []RadioStation
	Playlists     []Playlist

	SearchResultInfo
}

type SearchResultInfo struct {
	SeatchResults []SearchResult
	IsSearchFuzzy bool
	DidYouMean    string
}

func (*Client) SimpleSearch(q string, types []string) (*SearchResult, error) {
	return nil, nil
}

func (*Client) SimpleSearchWithInfo(q string, types []string) (*SearchResult, error) {
	return nil, nil
}

func (*Client) FindMatch(title string, artistname string, types string) (*SearchResult, error) {
	return nil, nil
}
