package vkutil

import (
	"code.witches.io/go/stellwerk/cmd/stellwerk-vulkan/vulkan"
	"fmt"
	"reflect"
	"unsafe"
)

type Buffer struct {
	Handle vulkan.Buffer
	Memory Memory

	device vulkan.Device
}

var allocator MemoryAllocator

func NewBuffer(device vulkan.Device, allocator MemoryAllocator, usage vulkan.BufferUsageFlags, data interface{}) (*Buffer, error) {
	rv := reflect.ValueOf(data)
	if rv.Kind() != reflect.Slice {
		return nil, fmt.Errorf("expected slice")
	}

	size := rv.Type().Elem().Size() * uintptr(rv.Len())
	pointer := rv.Pointer()

	b := &Buffer{
		device: device,
	}
	cleanup := b.Close
	defer func() {
		cleanup()
	}()

	buffer, err := vulkan.CreateBuffer(device, vulkan.BufferCreateInfo{
		Type:        vulkan.StructureTypeBufferCreateInfo,
		Size:        vulkan.DeviceSize(size),
		Usage:       usage,
		SharingMode: vulkan.SharingModeExclusive,
	}, nil)
	if err != nil {
		return nil, err
	}
	b.Handle = buffer

	requirements := vulkan.GetBufferMemoryRequirements(device, buffer)
	reqIndex, reqType, ok := allocator.Allocate(requirements)
	if !ok {
		return nil, fmt.Errorf("memory requirements can't be satisfied")
	}
	memory, err := vulkan.AllocateMemory(device, vulkan.MemoryAllocateInfo{
		Type:            vulkan.StructureTypeMemoryAllocateInfo,
		AllocationSize:  requirements.Size,
		MemoryTypeIndex: reqIndex,
	}, nil)
	if err != nil {
		return nil, err
	}
	b.Memory = Memory{
		Type:   reqType,
		Memory: memory,
		Device: device,
	}

	err = vulkan.BindBufferMemory(device, buffer, memory, 0)
	if err != nil {
		return nil, err
	}

	address, err := b.Memory.Map(0, vulkan.WholeSize, 0)

	vulkan.Memcpy(unsafe.Pointer(address), unsafe.Pointer(pointer), size)

	err = b.Memory.Unmap()
	if err != nil {
		return nil, err
	}

	cleanup = func() {}

	return b, nil
}

func (b *Buffer) Close() {
	if b == nil {
		return
	}

	if b.Handle != 0 {
		vulkan.DestroyBuffer(b.device, b.Handle, nil)
	}

	if b.Memory.Memory != 0 {
		vulkan.FreeMemory(b.device, b.Memory.Memory, nil)
	}
}
