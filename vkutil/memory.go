package vkutil

import (
	"code.witches.io/go/stellwerk/cmd/stellwerk-vulkan/vulkan"
)

type Memory struct {
	Type   vulkan.MemoryType
	Memory vulkan.DeviceMemory
	Device vulkan.Device

	offset vulkan.DeviceSize
	size   vulkan.DeviceSize
}

func (m *Memory) Map(offset, size vulkan.DeviceSize, flags vulkan.MemoryMapFlags) (uintptr, error) {
	addr, err := vulkan.MapMemory(m.Device, m.Memory, offset, size, flags)
	if err != nil {
		return 0, err
	}

	if (m.Type.PropertyFlags & vulkan.MemoryPropertyHostCoherentBit) != vulkan.MemoryPropertyHostCoherentBit {
		err = vulkan.InvalidateMappedMemoryRanges(m.Device, []vulkan.MappedMemoryRange{
			{
				Type:   vulkan.StructureTypeMappedMemoryRange,
				Memory: m.Memory,
				Offset: offset,
				Size:   size,
			},
		})
		m.offset = offset
		m.size = size
	}
	return addr, err
}

func (m *Memory) Unmap() error {
	var err error
	if (m.Type.PropertyFlags & vulkan.MemoryPropertyHostCoherentBit) != vulkan.MemoryPropertyHostCoherentBit {
		err = vulkan.FlushMappedMemoryRanges(m.Device, []vulkan.MappedMemoryRange{
			{
				Type:   vulkan.StructureTypeMappedMemoryRange,
				Memory: m.Memory,
				Offset: m.offset,
				Size:   m.size,
			},
		})
	}
	vulkan.UnmapMemory(m.Device, m.Memory)
	return err
}