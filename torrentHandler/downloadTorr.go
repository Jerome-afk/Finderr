package torrentDownload

import (
	"encoding/json"
	"fmt"
	
	"net/http"
	"os"
	"path/filepath"

	"finderr/torrentHandler/bencodeTorrent"
	"finderr/torrentHandler/peers"
	"finderr/torrentHandler/torrentfile"
	"finderr/torrentHandler/websocket"
)

type MediaFlags struct {
	Movies  string
	Anime   string
	TVShows string
	Music   string
}

var mediaFlags = MediaFlags{
	Movies:  "-mo",
	Anime:   "-an",
	TVShows: "-tv",
	Music:   "-mu",
}

func GetDestinationPath(mediaFlag string) string {
	baseDir := os.Getenv("DOWNLOAD_BASE_DIR")
	if baseDir == "" {
		baseDir = "downloads"
	}

	switch mediaFlag {
	case mediaFlags.Movies:
		return filepath.Join(baseDir, "movies")
	case mediaFlags.Anime:
		return filepath.Join(baseDir, "anime")
	case mediaFlags.TVShows:
		return filepath.Join(baseDir, "tv-shows")
	case mediaFlags.Music:
		return filepath.Join(baseDir, "music")
	default:
		return filepath.Join(baseDir, "others")
	}
}

func DownloadHandler(w http.ResponseWriter, r *http.Request, hub *websocket.Hub) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	mediaFlag := r.URL.Query().Get("media-flag")
	torrentFile := r.URL.Query().Get("torrent-file")

	if torrentFile == "" {
		http.Error(w, "Torrent file path is required", http.StatusBadRequest)
		return
	}

	dstPath := GetDestinationPath(mediaFlag)
	if err := os.MkdirAll(dstPath, 0755); err != nil {
		http.Error(w, "Failed to create destination directory", http.StatusInternalServerError)
		return
	}

	bct, err := bencodeTorrent.Decode(torrentFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to decode torrent: %v", err), http.StatusInternalServerError)
		return
	}

	p2p, err := peers.Peers(bct)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to get peers: %v", err), http.StatusInternalServerError)
		return
	}

	tf := torrentfile.Torrent{
		Peers:       p2p.Peers,
		PeerID:      [20]byte(p2p.PeerId),
		InfoHash:    bct.InfoHash,
		Name:        bct.Info.Name,
		PieceLength: bct.Info.PieceLength,
		Length:      bct.Info.Length,
		PieceHashes: bct.PieceHashes,
	}

	// Create progress channel and connect it to WebSocket
	tf.ProgressChan = make(chan float64, 1)
	go func() {
		for progress := range tf.ProgressChan {
			status := fmt.Sprintf("Downloading... %.1f%%", progress)
			hub.BroadcastProgress(progress, status)
		}
	}()

	file, err := tf.Download(dstPath)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to download: %v", err), http.StatusInternalServerError)
		return
	}
	defer file.Close()

	response := map[string]string{
		"status":  "success",
		"message": "Download completed successfully",
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
}