package main

import (
	"bytes"
	"errors"
	"io"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type ScanUploadConfig struct {
	Root                 string
	BodyLimit, FileLimit int64
	NextID               func() (string, error)
}
type scanUploader struct {
	cfg     ScanUploadConfig
	publish sync.Mutex
}

func scanReply(w http.ResponseWriter, r *http.Request, status int, text string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	if r.Method != "HEAD" {
		_, _ = io.WriteString(w, text)
	}
}
func scanReadStatus(err error) int {
	var limit *http.MaxBytesError
	if errors.As(err, &limit) {
		return 413
	}
	var disk *os.PathError
	if errors.As(err, &disk) {
		return 500
	}
	return 400
}
func scanFilename(name string) bool {
	if len(name) == 0 || len(name) > 80 || strings.Contains(name, "..") {
		return false
	}
	for i, c := range name {
		alnum := c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9'
		if !alnum && (i == 0 || c != '.' && c != '_' && c != '-') {
			return false
		}
	}
	return true
}
func scanToken(token string) bool {
	if len(token) != 16 {
		return false
	}
	for _, c := range token {
		if !(c >= '0' && c <= '9' || c >= 'a' && c <= 'f') {
			return false
		}
	}
	return true
}
func NewScanUpload(cfg ScanUploadConfig) (http.Handler, error) {
	info, err := os.Stat(cfg.Root)
	if err != nil {
		return nil, err
	}
	if !info.IsDir() || cfg.NextID == nil || cfg.BodyLimit < 1 || cfg.BodyLimit > 2<<20 || cfg.FileLimit < 8 || cfg.FileLimit > 1<<20 {
		return nil, errors.New("invalid scan upload configuration")
	}
	return &scanUploader{cfg: cfg}, nil
}
func (u *scanUploader) publishScan(stage string) (string, int) {
	u.publish.Lock()
	defer u.publish.Unlock()
	token, err := u.cfg.NextID()
	if err != nil || !scanToken(token) {
		return "", 500
	}
	name := token + ".png"
	target := filepath.Join(u.cfg.Root, name)
	if _, err = os.Lstat(target); err == nil {
		return "", 409
	} else if !errors.Is(err, os.ErrNotExist) {
		return "", 500
	}
	if err = os.Rename(stage, target); err != nil {
		return "", 500
	}
	return name, 201
}
func (u *scanUploader) receive(w http.ResponseWriter, r *http.Request) (string, int) {
	r.Body = http.MaxBytesReader(w, r.Body, u.cfg.BodyLimit)
	defer r.Body.Close()
	reader, err := r.MultipartReader()
	if err != nil {
		return "", 400
	}
	part, err := reader.NextPart()
	if err != nil {
		return "", scanReadStatus(err)
	}
	defer part.Close()
	_, params, err := mime.ParseMediaType(part.Header.Get("Content-Disposition"))
	if err != nil || params["name"] != "scan" || !scanFilename(params["filename"]) {
		return "", 400
	}
	stage, err := os.CreateTemp(u.cfg.Root, ".scan-*")
	if err != nil {
		return "", 500
	}
	// Cleanup belongs to every exit, including failed parsing and publication.
	defer os.Remove(stage.Name())
	defer stage.Close()
	n, err := io.Copy(stage, io.LimitReader(part, u.cfg.FileLimit+1))
	if n > u.cfg.FileLimit {
		return "", 413
	}
	if err != nil {
		return "", scanReadStatus(err)
	}
	if _, err = reader.NextPart(); err != io.EOF {
		if err != nil {
			return "", scanReadStatus(err)
		}
		return "", 400
	}
	if _, err = io.Copy(io.Discard, r.Body); err != nil {
		return "", scanReadStatus(err)
	}
	if _, err = stage.Seek(0, io.SeekStart); err != nil {
		return "", 500
	}
	signature := make([]byte, 8)
	if _, err = io.ReadFull(stage, signature); err != nil || !bytes.Equal(signature, []byte{137, 80, 78, 71, 13, 10, 26, 10}) {
		return "", 415
	}
	if err = stage.Close(); err != nil {
		return "", 500
	}
	if r.Context().Err() != nil {
		return "", 400
	}
	return u.publishScan(stage.Name())
}
func (u *scanUploader) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/cargo-scan" {
		scanReply(w, r, 404, "unknown upload\n")
		return
	}
	if r.Method != "POST" {
		w.Header().Set("Allow", "POST")
		scanReply(w, r, 405, "use POST\n")
		return
	}
	if r.Context().Err() != nil {
		scanReply(w, r, 400, "invalid scan\n")
		return
	}
	name, status := u.receive(w, r)
	if status == 201 {
		scanReply(w, r, 201, name+"\n")
		return
	}
	messages := map[int]string{400: "invalid scan\n", 409: "scan exists\n", 413: "scan too large\n", 415: "PNG required\n", 500: "scan unavailable\n"}
	scanReply(w, r, status, messages[status])
}
func main() {}
