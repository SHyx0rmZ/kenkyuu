package ase

import (
	"github.com/veandco/go-sdl2/sdl"
	"io"
	"unsafe"
)

func LoadASE(r io.Reader) (*sdl.Surface, error) {
	var state State
	var err error

	state.Header, err = loadHeader(r)
	if err != nil {
		return nil, err
	}

	for i := uint16(0); i < state.Header.frames; i++ {
		frame := loadFrame(state, r)

		state.frames = append(state.frames, frame)
	}

	var maskRed, maskGreen, maskBlue, maskAlpha uint32

	_, maskRed, maskGreen, maskBlue, maskAlpha, err = sdl.PixelFormatEnumToMasks(sdl.PIXELFORMAT_ABGR8888)
	if err != nil {
		return nil, err
	}

	surface, err := sdl.CreateRGBSurface(0, int32(state.Header.width), int32(state.Header.height), int32(state.Header.colorDepth), maskRed, maskGreen, maskBlue, maskAlpha)
	if err != nil {
		return nil, err
	}
	surface.FillRect(nil, sdl.Color{255, 0, 0, 255}.Uint32())

	for _, chunk := range state.frames[0].chunks {
		if chunk.chunkType == ChunkTypeCel {
			celBase := (*Cel)(unsafe.Pointer(chunk))

			switch celBase.celType {
			case CelTypeRaw:
				cel := (*CelRaw)(unsafe.Pointer(celBase))
				for y := uint(0); y < cel.height; y++ {
					src := cel.pixels
					dst := surface.Pixels()
					copy(dst[uint(surface.Pitch)*(y+uint(cel.positionY))+uint(cel.positionX)*uint(surface.BytesPerPixel()):], src[cel.width*y*state.Header.colorDepth/8:cel.width*(y+1)*state.Header.colorDepth/8])
				}
			case CelTypeCompressed:
				cel := (*CelCompressed)(unsafe.Pointer(celBase))
				for y := uint(0); y < cel.height; y++ {
					src := cel.pixels
					dst := surface.Pixels()
					copy(dst[uint(surface.Pitch)*(y+uint(cel.positionY))+uint(cel.positionX)*uint(surface.BytesPerPixel()):], src[cel.width*y*state.Header.colorDepth/8:cel.width*(y+1)*state.Header.colorDepth/8])
				}
			}
		}
	}

	return surface, nil
}
