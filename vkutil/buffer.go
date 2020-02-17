package vkutil

import (
	"code.witches.io/go/vulkan"
	"fmt"
	"reflect"
	"unsafe"
)

type Buffer struct {
	Handle vulkan.Buffer
	Memory Memory
	Device vulkan.Device
}

var allocator MemoryAllocator

func NewBuffer(device vulkan.Device, allocator MemoryAllocator, usage vulkan.BufferUsageFlags, data interface{}) (*Buffer, error) {
	rv := reflect.Indirect(reflect.ValueOf(data))
	var size uintptr
	var pointer uintptr

	switch rv.Kind() {
	case reflect.Slice:
		size = rv.Type().Elem().Size() * uintptr(rv.Len())
		pointer = rv.Pointer()
	case reflect.Struct:
		size = rv.Type().Size()
		pointer = rv.UnsafeAddr()
	default:
		return nil, fmt.Errorf("expected slice")
	}

	b := &Buffer{
		Device: device,
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
		vulkan.DestroyBuffer(b.Device, b.Handle, nil)
	}

	if b.Memory.Memory != 0 {
		vulkan.FreeMemory(b.Device, b.Memory.Memory, nil)
	}
}
