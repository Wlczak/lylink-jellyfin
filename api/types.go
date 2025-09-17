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
