package newmedia

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os/exec"
	"strings"

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
	var sb strings.Builder
	sb.WriteString("magnet:")
	sb.WriteString(fmt.Sprintf("?xt=urn:btih:%s", item.InfoHash))
	sb.WriteString(fmt.Sprintf("&dn=%s", url.QueryEscape(item.Name)))
	allTrackers := n.getTrackers()
	for _, tracker := range allTrackers {
		sb.WriteString(fmt.Sprintf("&tr=%s", url.QueryEscape(tracker)))
	}

	return sb.String()
}

func (n *NewMedia) StartDownloadNewMedia(item NewMediaSearchItem) error {
	magnetUrl := n.GetMagneticURL(item)

	log.Printf("Adding url to download queue: %s\n", magnetUrl)
	trCmd := exec.Command("transmission-remote", "-a", magnetUrl)
	if err := trCmd.Run(); err != nil {
		return err
	}

	return nil
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
