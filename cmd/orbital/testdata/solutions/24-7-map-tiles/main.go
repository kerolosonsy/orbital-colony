package main

import (
	"bytes"
	"context"
	"fmt"
	"image"
	"image/draw"
	"image/jpeg"
	"image/png"
	"sync"
)

type MapTile struct {
	Name string
	Data []byte
}
type TileSpec struct{ Width, Height, Workers, MaxPixels int }
type TileReply struct {
	Name, Format  string
	Width, Height int
	Err           error
}

func tileDimensions(w, h int, spec TileSpec) (int, int) {
	if int64(spec.Width)*int64(h) < int64(spec.Height)*int64(w) {
		if w > spec.Width {
			return spec.Width, max(1, int(int64(h)*int64(spec.Width)/int64(w)))
		}
	} else if h > spec.Height {
		return max(1, int(int64(w)*int64(spec.Height)/int64(h))), spec.Height
	}
	return w, h
}
func resizeTile(ctx context.Context, encoded []byte, spec TileSpec) ([]byte, string, int, int, error) {
	if err := ctx.Err(); err != nil {
		return nil, "", 0, 0, err
	}
	cfg, format, err := image.DecodeConfig(bytes.NewReader(encoded))
	if err != nil {
		return nil, "", 0, 0, err
	}
	if (format != "png" && format != "jpeg") || cfg.Width <= 0 || cfg.Height <= 0 || cfg.Width > spec.MaxPixels/cfg.Height {
		return nil, "", 0, 0, fmt.Errorf("invalid tile dimensions or format")
	}
	src, _, err := image.Decode(bytes.NewReader(encoded))
	if err != nil {
		return nil, "", 0, 0, err
	}
	w, h := tileDimensions(cfg.Width, cfg.Height, spec)
	dst := image.NewNRGBA(image.Rect(0, 0, w, h))
	bounds := src.Bounds()
	for y := 0; y < h; y++ {
		if err := ctx.Err(); err != nil {
			return nil, "", 0, 0, err
		}
		for x := 0; x < w; x++ {
			point := image.Pt(bounds.Min.X+x*cfg.Width/w, bounds.Min.Y+y*cfg.Height/h)
			draw.Draw(dst, image.Rect(x, y, x+1, y+1), src, point, draw.Src)
		}
	}
	var out bytes.Buffer
	if format == "png" {
		err = png.Encode(&out, dst)
	} else {
		err = jpeg.Encode(&out, dst, &jpeg.Options{Quality: 90})
	}
	return out.Bytes(), format, w, h, err
}
func saveMapTile(ctx context.Context, tile MapTile, spec TileSpec, save func(context.Context, string, []byte) error) TileReply {
	payload, format, w, h, err := resizeTile(ctx, tile.Data, spec)
	if err == nil {
		err = ctx.Err()
	}
	if err == nil {
		err = save(ctx, tile.Name, payload)
	}
	if err != nil {
		return TileReply{Name: tile.Name, Err: err}
	}
	return TileReply{Name: tile.Name, Format: format, Width: w, Height: h}
}
func BuildMapTiles(ctx context.Context, tiles []MapTile, spec TileSpec, save func(context.Context, string, []byte) error) ([]TileReply, error) {
	if save == nil || spec.Workers <= 0 || spec.Width < 1 || spec.Width > 2048 || spec.Height < 1 || spec.Height > 2048 || spec.MaxPixels < 1 || spec.MaxPixels > 1048576 {
		return nil, fmt.Errorf("invalid tile settings")
	}
	seen := map[string]bool{}
	for _, tile := range tiles {
		if tile.Name == "" || seen[tile.Name] {
			return nil, fmt.Errorf("invalid tile name")
		}
		seen[tile.Name] = true
	}
	replies := make([]TileReply, len(tiles))
	for i, tile := range tiles {
		replies[i] = TileReply{Name: tile.Name}
	}
	started := make([]bool, len(tiles))
	jobs := make(chan int)
	var wg sync.WaitGroup
	for i := 0; i < min(spec.Workers, len(tiles)); i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for index := range jobs {
				if ctx.Err() != nil {
					return
				}
				started[index] = true
				replies[index] = saveMapTile(ctx, tiles[index], spec, save)
			}
		}()
	}
feed:
	for i := range tiles {
		if ctx.Err() != nil {
			break
		}
		select {
		case jobs <- i:
		case <-ctx.Done():
			break feed
		}
	}
	close(jobs)
	wg.Wait()
	for i := range replies {
		if !started[i] {
			replies[i].Err = ctx.Err()
		}
	}
	return replies, nil
}
func main() {}
