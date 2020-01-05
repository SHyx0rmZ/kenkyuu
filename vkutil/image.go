package vkutil

import (
	sdl "code.witches.io/go/sdl2"
	"code.witches.io/go/stellwerk/cmd/stellwerk-vulkan/vulkan"
	"fmt"
	"image"
	"image/color"
	"image/jpeg"
	"image/png"
	"os"
	"path/filepath"
	"unsafe"
)

type Image struct {
	Image        image.Image
	Handle       vulkan.Image
	View         vulkan.ImageView
	Memory       vulkan.DeviceMemory
	Buffer       vulkan.Buffer
	BufferMemory vulkan.DeviceMemory

	vulkan.Extent2D

	Device vulkan.Device
}

type __image struct {
	addr   uintptr
	bounds image.Rectangle
	model  color.Model
}

func (i *__image) ColorModel() color.Model {
	return i.model
}

func (i *__image) Bounds() image.Rectangle {
	return i.bounds
}

func (i *__image) At(x, y int) color.Color {
	//const bytes = 8
	//canon := i.Bounds().Canon()
	//stride := canon.Dx() * bytes
	return color.RGBA{}
}

func (i *__image) Set(x, y int, c color.Color) {

}

func (i *__image) pixel(x, y int) *uint32 {
	return nil
}

func OpenImage(device vulkan.Device, allocator MemoryAllocator, path string, text *sdl.Surface) (*Image, error) {
	imageFile, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer func() {
		err := imageFile.Close()
		if err != nil {
			panic(err)
		}
	}()

	var j image.Image
	switch filepath.Ext(path) {
	case ".jpeg", ".jpg":
		j, err = jpeg.Decode(imageFile)
	case ".png":
		j, err = png.Decode(imageFile)
	default:
		panic("not implemented: " + filepath.Ext(path))
	}
	if err != nil {
		return nil, err
	}

	return LoadImage(device, allocator, j, text)
}

func LoadImage(device vulkan.Device, allocator MemoryAllocator, j image.Image, text *sdl.Surface) (*Image, error) {
	i := &Image{
		Image: j,
		Extent2D: vulkan.Extent2D{
			Width:  uint32(j.Bounds().Dx()),
			Height: uint32(j.Bounds().Dy()),
		},
		Device: device,
	}
	cleanup := i.Close
	defer func() {
		cleanup()
	}()

	imageBuffer, err := vulkan.CreateBuffer(device, vulkan.BufferCreateInfo{
		Type:        vulkan.StructureTypeBufferCreateInfo,
		Size:        vulkan.DeviceSize(4 * j.Bounds().Dx() * j.Bounds().Dy()),
		Usage:       vulkan.BufferUsageTransferSrcBit,
		SharingMode: vulkan.SharingModeExclusive,
	}, nil)
	if err != nil {
		return nil, err
	}
	i.Buffer = imageBuffer

	requirements := vulkan.GetBufferMemoryRequirements(device, imageBuffer)
	reqIndex, _, ok := allocator.Allocate(requirements)
	if !ok {
		return nil, fmt.Errorf("memory requirements can't be satisfied")
	}

	imageBufferMemory, err := vulkan.AllocateMemory(device, vulkan.MemoryAllocateInfo{
		Type:            vulkan.StructureTypeMemoryAllocateInfo,
		AllocationSize:  requirements.Size,
		MemoryTypeIndex: reqIndex,
	}, nil)
	if err != nil {
		return nil, err
	}
	i.BufferMemory = imageBufferMemory

	err = vulkan.BindBufferMemory(device, imageBuffer, imageBufferMemory, 0)
	if err != nil {
		return nil, err
	}
	fmt.Println("buffer:", imageBuffer, "requirements:", requirements, "size:", requirements.Size)
	addr, err := vulkan.MapMemory(device, imageBufferMemory, 0, vulkan.WholeSize, 0)
	if err != nil {
		return nil, err
	}
	err = vulkan.InvalidateMappedMemoryRanges(device, []vulkan.MappedMemoryRange{
		{
			Type:   vulkan.StructureTypeMappedMemoryRange,
			Memory: imageBufferMemory,
			Offset: 0,
			Size:   vulkan.WholeSize,
		},
	})
	if err != nil {
		return nil, err
	}

	fmt.Println("mapped memory at", unsafe.Pointer(addr))

	fmt.Println(unsafe.Pointer(addr), j.Bounds().Dx(), j.Bounds().Canon().Dx(), j.Bounds().Canon().Min, j.Bounds().Canon().Max, unsafe.Sizeof(uint16(0)))

	src := func(data uintptr, i image.Image, width int) {
		const bytes = 4
		canon := i.Bounds().Canon()
		stride := width * bytes
		for y := canon.Min.Y; y < canon.Max.Y; y++ {
			for x := canon.Min.X; x < canon.Max.X; x++ {
				r, g, b, a := i.At(x, y).RGBA()
				*(*uint8)(unsafe.Pointer(addr + uintptr(y*stride+x*bytes+0))) = uint8(r)
				*(*uint8)(unsafe.Pointer(addr + uintptr(y*stride+x*bytes+1))) = uint8(g)
				*(*uint8)(unsafe.Pointer(addr + uintptr(y*stride+x*bytes+2))) = uint8(b)
				*(*uint8)(unsafe.Pointer(addr + uintptr(y*stride+x*bytes+3))) = uint8(a)
			}
		}
	}
	over := func(data uintptr, i image.Image, width int) {
		const bytes = 4
		const m = (1 << 8) - 1
		canon := i.Bounds().Canon()
		stride := width * bytes
		for y := canon.Min.Y; y < canon.Max.Y; y++ {
			for x := canon.Min.X; x < canon.Max.X; x++ {
				r, g, b, a := i.At(x, y).RGBA()
				rp := (*uint8)(unsafe.Pointer(addr + uintptr(y*stride+x*bytes+0)))
				gp := (*uint8)(unsafe.Pointer(addr + uintptr(y*stride+x*bytes+1)))
				bp := (*uint8)(unsafe.Pointer(addr + uintptr(y*stride+x*bytes+2)))
				ap := (*uint8)(unsafe.Pointer(addr + uintptr(y*stride+x*bytes+3)))
				*rp = uint8(uint32(*rp)*(m-a)/m + r)
				*gp = uint8(uint32(*gp)*(m-a)/m + g)
				*bp = uint8(uint32(*bp)*(m-a)/m + b)
				*ap = uint8(uint32(*ap)*(m-a)/m + a)
			}
		}
	}
	over = over

	{
		//err = text.Lock()
		//if err != nil {
		//	panic(err)
		//}
		//defer text.Unlock()
		//
		src(addr, j, j.Bounds().Dx())
		if text != nil {
			over(addr, text, j.Bounds().Dx())
		}
	}

	err = vulkan.FlushMappedMemoryRanges(device, []vulkan.MappedMemoryRange{
		{
			Type:   vulkan.StructureTypeMappedMemoryRange,
			Memory: imageBufferMemory,
			Offset: 0,
			Size:   vulkan.WholeSize,
		},
	})
	if err != nil {
		return nil, err
	}
	vulkan.UnmapMemory(device, imageBufferMemory)

	imageImage, err := vulkan.CreateImage(device, vulkan.ImageCreateInfo{
		Type:      vulkan.StructureTypeImageCreateInfo,
		ImageType: vulkan.ImageType2D,
		//Format:    vulkan.FormatR16G16B16A16UInt,
		Format: vulkan.FormatR8G8B8A8UInt,
		Extent: vulkan.Extent3D{
			Width:  uint32(j.Bounds().Dx()),
			Height: uint32(j.Bounds().Dy()),
			Depth:  1,
		},
		MipLevels:     1,
		ArrayLayers:   1,
		Samples:       vulkan.SampleCount1Bit,
		Tiling:        vulkan.ImageTilingOptimal,
		Usage:         vulkan.ImageUsageTransferDstBit | vulkan.ImageUsageSampledBit,
		SharingMode:   vulkan.SharingModeExclusive,
		InitialLayout: vulkan.ImageLayoutUndefined,
	}, nil)
	if err != nil {
		return nil, err
	}
	i.Handle = imageImage

	requirements = vulkan.GetImageMemoryRequirements(device, imageImage)
	reqIndex, _, ok = allocator.Allocate(requirements)
	if !ok {
		return nil, fmt.Errorf("memory requirements can't be satisfied")
	}

	imageImageMemory, err := vulkan.AllocateMemory(device, vulkan.MemoryAllocateInfo{
		Type:            vulkan.StructureTypeMemoryAllocateInfo,
		AllocationSize:  requirements.Size,
		MemoryTypeIndex: reqIndex,
	}, nil)
	if err != nil {
		return nil, err
	}
	i.Memory = imageImageMemory

	err = vulkan.BindImageMemory(device, imageImage, imageImageMemory, 0)
	if err != nil {
		return nil, err
	}

	imageImageView, err := vulkan.CreateImageView(device, vulkan.ImageViewCreateInfo{
		Type:     vulkan.StructureTypeImageViewCreateInfo,
		Image:    imageImage,
		ViewType: vulkan.ImageViewType2D,
		//Format:     vulkan.FormatR16G16B16A16UInt,
		Format:     vulkan.FormatR8G8B8A8UInt,
		Components: vulkan.ComponentMapping{},
		SubresourceRange: vulkan.ImageSubresourceRange{
			AspectMask:     vulkan.ImageAspectColorBit,
			BaseMIPLevel:   0,
			LevelCount:     1,
			BaseArrayLayer: 0,
			LayerCount:     1,
		},
	}, nil)
	if err != nil {
		return nil, err
	}
	i.View = imageImageView

	cleanup = func() {}

	return i, nil
}

func (i *Image) Close() {
	if i == nil {
		return
	}

	if i.View != 0 {
		vulkan.DestroyImageView(i.Device, i.View, nil)
	}

	if i.Handle != 0 {
		vulkan.DestroyImage(i.Device, i.Handle, nil)
	}

	if i.Memory != 0 && i.Memory != i.BufferMemory {
		vulkan.FreeMemory(i.Device, i.Memory, nil)
	}

	if i.Buffer != 0 {
		vulkan.DestroyBuffer(i.Device, i.Buffer, nil)
	}

	if i.BufferMemory != 0 {
		vulkan.FreeMemory(i.Device, i.BufferMemory, nil)
	}
}
