package torrentfile

import (
	"bytes"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"finderr/torrentHandler/bitfield"
	handShake "finderr/torrentHandler/handshake"
	"finderr/torrentHandler/message"
	"finderr/torrentHandler/peers"
	"finderr/torrentHandler/worker"
)

type resultPiece struct {
	pieceIndex int
	pieceHash  [20]byte
	length     int
	buf        []byte
}

type Torrent struct {
	Peers        []peers.Peer
	PeerID       [20]byte
	InfoHash     [20]byte
	PieceHashes  [][20]byte
	PieceLength  int
	Length       int
	Name         string
	ProgressChan chan float64
}

func (t *Torrent) performHandShake(conn net.Conn) error {
	hs := &handShake.HandShake{
		PeerID:   t.PeerID,
		InfoHash: t.InfoHash,
		Protocol: "BitTorrent protocol",
	}

	bufHs := hs.Serialize()
	_, err := conn.Write(bufHs)
	if err != nil {
		return err
	}

	readHs, err := handShake.Read(conn)
	if err != nil {
		return err
	}

	if !bytes.Equal(hs.InfoHash[:], readHs.InfoHash[:]) {
		return fmt.Errorf("wrong file")
	}
	return nil
}

func (t *Torrent) initDownloadWorker(peer peers.Peer, resultsQueue chan<- resultPiece, workQueue chan *worker.WorkPiece) error {

	//connect to peer
	conn, err := net.DialTimeout("tcp", peer.String(), 10*time.Second)
	if err != nil {
		return fmt.Errorf("failed to connect to peer %w", err)
	}

	//set read and write connection timeout
	conn.SetDeadline(time.Now().Add(10 * time.Minute))
	defer func() {
		conn.SetDeadline(time.Time{})
		conn.Close()
	}()

	// handshake peer
	err = t.performHandShake(conn)
	if err != nil {
		return fmt.Errorf("failed to complete handshake with peer %w", err)
	}

	//receive bitfield
	bf, err := bitfield.Read(conn)
	if err != nil {
		return fmt.Errorf("failed to receive bitfield from peer %w", err)
	}

	//worker creation
	wkr := &worker.Worker{
		Conn: conn,
		Bf:   bf,
	}

	err = wkr.Interested()
	if err != nil {
		return fmt.Errorf("failed to send interested message %w", err)
	}

	err = wkr.SendUnchoke()
	if err != nil {
		return fmt.Errorf("failed to send unchoke message %w", err)
	}

	msg, err := message.ReadMessage(wkr.Conn)
	if err != nil {
		log.Print(err)
		return err
	}

	if msg == nil {
		return nil
	}

	if msg.ID != message.MsgUnchoke {
		return err
	}

	for wp := range workQueue {
		if !wkr.Bf.HasPiece(wp.PieceIndex) {
			//Try another workpiece
			workQueue <- wp
			continue
		}

		buf, err := wkr.Download(wp.PieceIndex, wp.Length)

		if err != nil {
			//Try another workpiece
			workQueue <- wp
			continue
		}

		if !wkr.CheckPieceIntegrity(buf, wp.PieceHash) {
			//Try another workpiece
			workQueue <- wp
			continue
		}
		//send to piece to  resultPiece chan
		rp := resultPiece{
			pieceIndex: wp.PieceIndex,
			buf:        buf,
			length:     wp.Length,
			pieceHash:  wp.PieceHash,
		}
		resultsQueue <- rp
	}
	return nil
}

func (t *Torrent) Download(dstPath string) (*os.File, error) {
	activePeersCount := len(t.Peers)
	donePieces := 0
	totalPieces := len(t.PieceHashes)

	activePeers := make(chan peers.Peer, activePeersCount)
	errors := make(chan error, activePeersCount)
	workQueue := make(chan *worker.WorkPiece, totalPieces)
	resultsQueue := make(chan resultPiece, totalPieces)
	downloadCompleted := make(chan bool)

	// Initialize progress channel if not already set
	if t.ProgressChan == nil {
		t.ProgressChan = make(chan float64, 1)
	}

	// Function to update progress
	updateProgress := func() {
		progress := float64(donePieces) / float64(totalPieces) * 100
		t.ProgressChan <- progress
	}

	// Open the file for writing (create if it doesn't exist, truncate if it does)
	file, err := os.Create(dstPath)
	if err != nil {
		return nil, fmt.Errorf("error opening file: %w", err)

	}
	//seek "pre-allocates" before we actually start to write in it
	_, err = file.Seek((int64(t.Length))-1, 0)
	if err != nil {
		return nil, fmt.Errorf("error creating sparse file: %w", err)
	}

	//TODO: Don't load the entire file into memory first before writing to disk
	//fileBuf := make([]byte, t.Length)

	//creating queue of workPieces
	for index, pieceHash := range t.PieceHashes {
		workQueue <- &worker.WorkPiece{PieceIndex: index, PieceHash: pieceHash, Length: t.PieceLength}
	}

	for _, activePeer := range t.Peers {
		go func(peer peers.Peer) {
			//returns when there's nothing more to download in other words workQueue is empty or error
			err := t.initDownloadWorker(peer, resultsQueue, workQueue)
			if err != nil {
				errors <- err
			}

		}(activePeer)
	}

	go func() {
		<-downloadCompleted
		//download is completed close all channels
		log.Println("Download completed, closing all channels")
		close(resultsQueue)
		close(workQueue)
		close(activePeers)
		close(t.ProgressChan)
	}()

	go func() {
		for range resultsQueue {
			donePieces++
			updateProgress()
		}
	}()

	go func() {
		for error := range errors {
			//handle error from connection related issues with peers
			log.Println(error.Error())

		}

	}()

	for resultPiece := range resultsQueue {
		begin := resultPiece.pieceIndex * resultPiece.length
		_, err := file.WriteAt(resultPiece.buf, int64(begin))
		if err != nil {
			return nil, fmt.Errorf("error creating sparse file: %w", err)
		}

		donePieces++

		if donePieces == len(t.PieceHashes) {
			downloadCompleted <- true
		}

		percentComplete := float64(donePieces) / float64(len(t.PieceHashes)) * 100
		log.Printf("%.2f%% complete\n", percentComplete)
	}
	return file, nil

}
