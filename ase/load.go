package ase

import (
	"io"
	"unsafe"

	"code.witches.io/go/sdl2"
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

	surface, err := sdl.CreateRGBSurfaceWithFormat(0, int(state.Header.width), int(state.Header.height), int(state.Header.colorDepth), sdl.PixelFormatABGR8888)
	if err != nil {
		return nil, err
	}
	format := surface.Format()

	err = surface.FillRect(nil, sdl.MapRGBA(&format, 255, 0, 0, 255))
	if err != nil {
		return nil, err
	}

	for _, chunk := range state.frames[0].chunks {
		if chunk.chunkType == ChunkTypeCel {
			celBase := (*Cel)(unsafe.Pointer(chunk))

			switch celBase.celType {
			case CelTypeRaw:
				cel := (*CelRaw)(unsafe.Pointer(celBase))
				for y := uint(0); y < cel.height; y++ {
					src := cel.pixels
					dst := surface.Pixels()
					copy(dst[uint(surface.Pitch())*(y+uint(cel.positionY))+uint(cel.positionX)*uint(format.BytesPerPixel()):], src[cel.width*y*state.Header.colorDepth/8:cel.width*(y+1)*state.Header.colorDepth/8])
				}
			case CelTypeCompressed:
				cel := (*CelCompressed)(unsafe.Pointer(celBase))
				for y := uint(0); y < cel.height; y++ {
					src := cel.pixels
					dst := surface.Pixels()
					copy(dst[uint(surface.Pitch())*(y+uint(cel.positionY))+uint(cel.positionX)*uint(format.BytesPerPixel()):], src[cel.width*y*state.Header.colorDepth/8:cel.width*(y+1)*state.Header.colorDepth/8])
				}
			}
		}
	}

	return surface, nil
}
