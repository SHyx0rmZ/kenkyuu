package vkutil

import (
	"code.witches.io/go/stellwerk/cmd/stellwerk-vulkan/vulkan"
)

type MemoryAllocator interface {
	Allocate(requirements vulkan.MemoryRequirements) (uint32, vulkan.MemoryType, bool)
}

type GPUMemoryAllocator vulkan.PhysicalDeviceMemoryProperties

func (a GPUMemoryAllocator) Allocate(requirements vulkan.MemoryRequirements) (uint32, vulkan.MemoryType, bool) {
	for i, memoryType := range a.MemoryTypes[:a.MemoryTypeCount] {
		if memoryType.PropertyFlags&vulkan.MemoryPropertyDeviceLocalBit == 0 {
			continue
		}
		if requirements.MemoryTypeBits&(1<<uint(i)) != (1 << uint(i)) {
			continue
		}
		return uint32(i), memoryType, true
	}
	return 0, vulkan.MemoryType{}, false
}

type CPUMemoryAllocator vulkan.PhysicalDeviceMemoryProperties

func (a CPUMemoryAllocator) Allocate(requirements vulkan.MemoryRequirements) (uint32, vulkan.MemoryType, bool) {
	for i, memoryType := range a.MemoryTypes[:a.MemoryTypeCount] {
		if memoryType.PropertyFlags&vulkan.MemoryPropertyHostVisibleBit == 0 {
			continue
		}
		if requirements.MemoryTypeBits&(1<<uint(i)) != (1 << uint(i)) {
			continue
		}
		return uint32(i), memoryType, true
	}
	return 0, vulkan.MemoryType{}, false
}
