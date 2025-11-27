package newmedia

type NewMediaSearchItem struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	InfoHash string `json:"info_hash"`
	Leechers string `json:"leechers"`
	Seeders  string `json:"seeders"`
	NumFiles string `json:"num_files"`
	Size     string `json:"size"`
	Username string `json:"username"`
	Added    string `json:"added"`
	Status   string `json:"status"`
	Category string `json:"category"`
	Imdb     string `json:"imdb"`
}

type TransmissionRPCPayload struct {
	Tag       string `json:"tag,omitempty"`
	Method    string `json:"method"`
	Arguments any    `json:"arguments,omitempty"`
}

type TransmissionRPCResponse struct {
	Tag       string `json:"tag,omitempty"`
	Result    string `json:"result"`
	Arguments any    `json:"arguments,omitempty"`
}

type TransmissionTorrentAddArguments struct {
	Filename string `json:"filename"`
}

type TransmissionHeader struct {
	XTransmissionSessionId string `json:"X-Transmission-Session-Id"`
}
