package vkutil

import (
	"fmt"
	"reflect"
	"unsafe"

	"code.witches.io/go/vulkan"
)

type TransferBuffer struct {
	HandleHost   vulkan.Buffer
	HandleDevice vulkan.Buffer
	MemoryHost   Memory
	MemoryDevice Memory
	Device       vulkan.Device

	fromFamily uint32
	toFamily   uint32
	size       vulkan.DeviceSize

	free           func()
	commandPool    vulkan.CommandPool
	commandBuffers []vulkan.CommandBuffer
}

func NewTransferBuffer(device vulkan.Device, allocatorFrom, allocatorTo MemoryAllocator, fromFamily, toFamily uint32, usage vulkan.BufferUsageFlags, data any) (*TransferBuffer, error) {
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

	b := &TransferBuffer{
		Device:     device,
		fromFamily: fromFamily,
		toFamily:   toFamily,
		size:       vulkan.DeviceSize(size),
	}
	cleanup := b.Close
	defer func() {
		cleanup()
	}()

	buffer, err := vulkan.CreateBuffer(device, vulkan.BufferCreateInfo{
		Type:               vulkan.StructureTypeBufferCreateInfo,
		Size:               vulkan.DeviceSize(size),
		Usage:              vulkan.BufferUsageTransferSrcBit,
		SharingMode:        vulkan.SharingModeConcurrent,
		QueueFamilyIndices: []uint32{fromFamily, toFamily},
	}, nil)
	if err != nil {
		return nil, err
	}
	b.HandleHost = buffer

	requirements := vulkan.GetBufferMemoryRequirements(device, buffer)
	reqIndex, reqType, ok := allocatorFrom.Allocate(requirements)
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
	b.MemoryHost = Memory{
		Type:   reqType,
		Memory: memory,
		Device: device,
	}

	err = vulkan.BindBufferMemory(device, buffer, memory, 0)
	if err != nil {
		return nil, err
	}

	address, err := b.MemoryHost.Map(0, vulkan.WholeSize, 0)

	vulkan.Memcpy(unsafe.Pointer(address), unsafe.Pointer(pointer), size)

	err = b.MemoryHost.Unmap()
	if err != nil {
		return nil, err
	}

	buffer, err = vulkan.CreateBuffer(device, vulkan.BufferCreateInfo{
		Type:               vulkan.StructureTypeBufferCreateInfo,
		Size:               vulkan.DeviceSize(size),
		Usage:              usage | vulkan.BufferUsageTransferDstBit,
		SharingMode:        vulkan.SharingModeConcurrent,
		QueueFamilyIndices: []uint32{fromFamily, toFamily},
	}, nil)
	if err != nil {
		return nil, err
	}
	b.HandleDevice = buffer

	reqIndex, reqType, ok = allocatorTo.Allocate(requirements)
	if !ok {
		return nil, fmt.Errorf("memory requirements can't be satisfied")
	}
	memory, err = vulkan.AllocateMemory(device, vulkan.MemoryAllocateInfo{
		Type:            vulkan.StructureTypeMemoryAllocateInfo,
		AllocationSize:  requirements.Size,
		MemoryTypeIndex: reqIndex,
	}, nil)
	if err != nil {
		return nil, err
	}
	b.MemoryDevice = Memory{
		Type:   reqType,
		Memory: memory,
		Device: device,
	}

	err = vulkan.BindBufferMemory(device, buffer, memory, 0)
	if err != nil {
		return nil, err
	}

	cleanup = func() {}

	return b, nil
}

type Transferer interface {
	Transfer(queue vulkan.Queue, semaphore vulkan.Semaphore, stageMask vulkan.PipelineStageFlags2, fence vulkan.Fence) error
}

func (b *TransferBuffer) Transfer(queue vulkan.Queue, semaphore vulkan.Semaphore, stageMask vulkan.PipelineStageFlags2, fence vulkan.Fence) error {
	commandPool, err := vulkan.CreateCommandPool(b.Device, vulkan.CommandPoolCreateInfo{
		Type:             vulkan.StructureTypeCommandPoolCreateInfo,
		Flags:            vulkan.CommandPoolCreateTransient,
		QueueFamilyIndex: b.fromFamily,
	}, nil)
	if err != nil {
		return err
	}
	b.commandPool = commandPool

	commandBuffers, err := vulkan.AllocateCommandBuffers(b.Device, vulkan.CommandBufferAllocateInfo{
		Type:               vulkan.StructureTypeCommandBufferAllocateInfo,
		CommandPool:        commandPool,
		Level:              vulkan.CommandBufferLevelPrimary,
		CommandBufferCount: 1,
	})
	if err != nil {
		return err
	}
	b.commandBuffers = commandBuffers

	err = vulkan.BeginCommandBuffer(commandBuffers[0], vulkan.CommandBufferBeginInfo{
		Type:            vulkan.StructureTypeCommandBufferBeginInfo,
		Flags:           vulkan.CommandBufferUsageOneTimeSubmitBit,
		InheritanceInfo: nil,
	})
	if err != nil {
		return err
	}

	vulkan.CmdCopyBuffer2(commandBuffers[0], vulkan.CopyBufferInfo2{
		Type:      vulkan.StructureTypeCopyBufferInfo2,
		SrcBuffer: b.HandleHost,
		DstBuffer: b.HandleDevice,
		Regions: []vulkan.BufferCopy2{
			{
				Type:      vulkan.StructureTypeBufferCopy2,
				SrcOffset: 0,
				DstOffset: 0,
				Size:      b.size,
			},
		},
	})

	err = vulkan.EndCommandBuffer(commandBuffers[0])
	if err != nil {
		return err
	}

	info := vulkan.SubmitInfo2{
		Type: vulkan.StructureTypeSubmitInfo2,
		CommandBufferInfos: []vulkan.CommandBufferSubmitInfo{
			{
				Type:          vulkan.StructureTypeCommandBufferSubmitInfo,
				CommandBuffer: commandBuffers[0],
			},
		},
	}

	if semaphore != vulkan.NullHandle {
		info.SignalSemaphoreInfos = []vulkan.SemaphoreSubmitInfo{
			{
				Type:      vulkan.StructureTypeSemaphoreSubmitInfo,
				Semaphore: semaphore,
				StageMask: stageMask,
			},
		}
	}

	free, err := vulkan.QueueSubmit2(queue, []vulkan.SubmitInfo2{info}, fence)
	if err != nil {
		free()
		return err
	}
	b.free = func() {
		b.free = nil
		free()
	}

	return nil
}

func (b *TransferBuffer) Buffer() vulkan.Buffer {
	return b.HandleDevice
}

func (b *TransferBuffer) Close() {
	if b == nil {
		return
	}

	if b.free != nil {
		b.free()
	}

	if b.commandBuffers != nil {
		vulkan.FreeCommandBuffers(b.Device, b.commandPool, b.commandBuffers)
	}

	if b.commandPool != 0 {
		vulkan.DestroyCommandPool(b.Device, b.commandPool, nil)
	}

	if b.HandleDevice != 0 {
		vulkan.DestroyBuffer(b.Device, b.HandleDevice, nil)
	}

	if b.MemoryDevice.Memory != 0 {
		vulkan.FreeMemory(b.Device, b.MemoryDevice.Memory, nil)
	}

	if b.HandleHost != 0 {
		vulkan.DestroyBuffer(b.Device, b.HandleHost, nil)
	}

	if b.MemoryHost.Memory != 0 {
		vulkan.FreeMemory(b.Device, b.MemoryHost.Memory, nil)
	}
}
