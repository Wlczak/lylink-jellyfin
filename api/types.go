package api

type Api struct {
	Username    string
	AccessToken string
}

type AuthResponse struct {
	AccessToken string `json:"AccessToken"`
}

type SessionItem struct {
	PlayState      PlayState      `json:"PlayState"`
	NowPlayingItem NowPlayingItem `json:"NowPlayingItem"`
}

type PlayState struct {
	PositionTicks       int64  `json:"PositionTicks"`
	CanSeek             bool   `json:"CanSeek"`
	IsPaused            bool   `json:"IsPaused"`
	IsMuted             bool   `json:"IsMuted"`
	VolumeLevel         int    `json:"VolumeLevel"`
	AudioStreamIndex    int    `json:"AudioStreamIndex"`
	SubtitleStreamIndex int    `json:"SubtitleStreamIndex"`
	MediaSourceId       string `json:"MediaSourceId"`
	PlayMethod          string `json:"PlayMethod"`
	RepeatMode          string `json:"RepeatMode"`
	PlaybackOrder       string `json:"PlaybackOrder"`
}

type NowPlayingItem struct {
	RunTimeTicks int64 `json:"RunTimeTicks"`
}

type Media interface {
}

type MediaInfo struct {
	Type string `json:"Type"`
}

type EpisodeInfo struct {
	Id                string `json:"Id"`
	Name              string `json:"Name"`
	Type              string `json:"Type"`
	SeriesName        string `json:"SeriesName"`
	IndexNumber       int    `json:"IndexNumber"`
	ParentIndexNumber int    `json:"ParentIndexNumber"`
	ParentId          string `json:"ParentId"`
}

type SeasonInfo struct {
	Id       string `json:"Id"`
	ParentId string `json:"ParentId"`
}

type SeriesInfo struct {
	Id       string `json:"Id"`
	ParentId string `json:"ParentId"`
}

type GetEpisodeInfoWithParentsResponse struct {
	Id                string `json:"Id"`
	Name              string `json:"Name"`
	Type              string `json:"Type"`
	SeriesName        string `json:"SeriesName"`
	IndexNumber       int    `json:"IndexNumber"`
	ParentIndexNumber int    `json:"ParentIndexNumber"`
	SeasonId          string `json:"SeasonId"`
	SeriesId          string `json:"SeriesId"`
}

type GetTokenRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type GetPlaybackInfoRequest struct {
	AccessToken string `json:"token"`
}

type GetMediaInfoRequest struct {
	AccessToken string `json:"token"`
}
