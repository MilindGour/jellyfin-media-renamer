package newmedia

import (
	"encoding/json"
	"fmt"
	"net/url"

	"github.com/MilindGour/jellyfin-media-renamer/network"
)

type NewMedia struct {
	h network.HttpResponseProvider
}

func NewNewMedia(h network.HttpResponseProvider) *NewMedia {
	return &NewMedia{
		h: h,
	}
}

func (n *NewMedia) SearchMedia(searchTerm string) []NewMediaSearchItem {
	searchUrl := fmt.Sprintf("https://apibay.org/q.php?q=%s", url.QueryEscape(searchTerm))

	response, err := n.h.GetResponse(searchUrl)
	if err != nil {
		return []NewMediaSearchItem{}
	}

	var out []NewMediaSearchItem
	if jsonErr := json.NewDecoder(response.Body).Decode(&out); jsonErr != nil {
		return []NewMediaSearchItem{}
	}

	return out
}

func (n *NewMedia) GetMagneticURL(item NewMediaSearchItem) string {
	q := url.Values{}
	q.Add("xt", fmt.Sprintf("urn:btih:%s", item.InfoHash))
	q.Add("dn", url.QueryEscape(item.Name))

	allTrackers := n.getTrackers()
	for _, tracker := range allTrackers {
		q.Add("tr", url.QueryEscape(tracker))
	}

	finalURL := fmt.Sprintf("magnet:?%s", q.Encode())
	return finalURL
}

func (n *NewMedia) getTrackers() []string {
	return []string{
		"udp://tracker.opentrackr.org:1337",
		"udp://open.stealth.si:80/announce",
		"udp://tracker.torrent.eu.org:451/announce",
		"udp://tracker.bittor.pw:1337/announce",
		"udp://public.popcorn-tracker.org:6969/announce",
		"udp://tracker.dler.org:6969/announce",
		"udp://exodus.desync.com:6969",
		"udp://open.demonii.com:1337/announce",
		"udp://glotorrents.pw:6969/announce",
		"udp://tracker.coppersurfer.tk:6969",
		"udp://torrent.gresille.org:80/announce",
		"udp://p4p.arenabg.com:1337",
		"udp://tracker.internetwarriors.net:1337",
	}
}
