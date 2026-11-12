package main

import (
	"bytes"
	"context"
	"encoding/binary"
	"errors"
	"io"
	"net"
	"net/http"
	"net/url"
	"time"
)

type SurveyUplink struct {
	Client     *http.Client
	URL, Token string
}

var ErrSurveyFrame = errors.New("invalid survey frame")
var ErrSurveyUpstream = errors.New("uplink rejected survey")

func RelaySurvey(conn net.Conn, uplink SurveyUplink, deadline time.Time) error {
	defer conn.Close()
	target, err := url.Parse(uplink.URL)
	if err != nil || target.Scheme != "https" || target.Host == "" || target.User != nil {
		return ErrSurveyUpstream
	}
	if err = conn.SetDeadline(deadline); err != nil {
		return err
	}
	var header [2]byte
	if _, err = io.ReadFull(conn, header[:]); err != nil {
		return err
	}
	size := int(binary.BigEndian.Uint16(header[:]))
	if size < 1 || size > 64 {
		return ErrSurveyFrame
	}
	name := make([]byte, size)
	if _, err = io.ReadFull(conn, name); err != nil {
		return err
	}
	query := target.Query()
	query.Set("station", string(name))
	target.RawQuery = query.Encode()
	ctx, cancel := context.WithDeadline(context.Background(), deadline)
	defer cancel()
	for attempt := 0; attempt < 2; attempt++ {
		req, err := http.NewRequestWithContext(ctx, "GET", target.String(), nil)
		if err != nil {
			return err
		}
		req.Header.Set("Authorization", "Bearer "+uplink.Token)
		resp, err := uplink.Client.Do(req)
		if err != nil {
			return err
		}
		if resp.StatusCode == 503 && attempt == 0 {
			resp.Body.Close()
			continue
		}
		if resp.StatusCode != 200 {
			resp.Body.Close()
			return ErrSurveyUpstream
		}
		data, err := io.ReadAll(io.LimitReader(resp.Body, 1025))
		resp.Body.Close()
		if err != nil {
			return err
		}
		if len(data) > 1024 {
			return ErrSurveyFrame
		}
		binary.BigEndian.PutUint16(header[:], uint16(len(data)))
		frame := append(header[:], data...)
		_, err = io.Copy(conn, bytes.NewReader(frame))
		return err
	}
	return ErrSurveyUpstream
}
func main() {}
