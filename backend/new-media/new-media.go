package newmedia

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
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

	return n.AddTransmissionDownload(magnetUrl)
}

func (n *NewMedia) AddTransmissionDownload(magnetUrl string) error {
	tsid := n.getTransmissionSessionId()

	if len(tsid) == 0 {
		return errors.New("Unable to fetch transmission session id from given rpc url")
	}

	response, err := n.SendRPCToTransmission(&TransmissionRPCPayload{
		Tag:    "77",
		Method: "torrent-add",
		Arguments: TransmissionTorrentAddArguments{
			Filename: magnetUrl,
		},
	}, tsid)

	if err != nil {
		return err
	}

	var rpcResponse TransmissionRPCResponse
	err = json.NewDecoder(response.Body).Decode(&rpcResponse)

	if err != nil {
		return err
	}

	if rpcResponse.Result != "success" {
		return errors.New("RPC call failed. Result was not success")
	}

	return nil
}

func (n *NewMedia) getTransmissionSessionId() string {
	res, err := n.SendRPCToTransmission(nil, "")
	if err != nil {
		fmt.Printf("Error while trying to fetch transmission session id: %s", err.Error())
		return ""
	}

	if res.StatusCode == 409 {
		return res.Header.Get("X-Transmission-Session-Id")
	}

	return ""
}

func (n *NewMedia) SendRPCToTransmission(rpcPayload *TransmissionRPCPayload, transmissionSessionId string) (*http.Response, error) {
	rpcURL := os.Getenv("TRANSMISSION_RPC_URL")

	if len(rpcURL) > 0 {
		response, err := n.h.PostJSON(rpcURL, rpcPayload, &http.Header{
			"X-Transmission-Session-Id": []string{transmissionSessionId},
		})

		if err != nil {
			return nil, err
		}

		// NOTE: response.StatusCode is not guaranteed to be 200
		return response, nil

	} else {
		return nil, errors.New("Cannot find environment variable TRANSMISSION_RPC_URL")
	}
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
