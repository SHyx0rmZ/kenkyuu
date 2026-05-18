package vkutil

import "code.witches.io/go/vulkan"

func CmdPipelineBarrier[Barrier vulkan.MemoryBarrier | vulkan.BufferMemoryBarrier | vulkan.ImageMemoryBarrier](
	commandBuffer vulkan.CommandBuffer,
	srcStageMask, dstStageMask vulkan.PipelineStageFlags,
	dependencyFlags vulkan.DependencyFlags,
	barriers ...Barrier,
) {
	switch barriers := any(barriers).(type) {
	case []vulkan.MemoryBarrier:
		vulkan.CmdPipelineBarrier(commandBuffer, srcStageMask, dstStageMask, dependencyFlags, barriers, nil, nil)
	case []vulkan.BufferMemoryBarrier:
		vulkan.CmdPipelineBarrier(commandBuffer, srcStageMask, dstStageMask, dependencyFlags, nil, barriers, nil)
	case []vulkan.ImageMemoryBarrier:
		vulkan.CmdPipelineBarrier(commandBuffer, srcStageMask, dstStageMask, dependencyFlags, nil, nil, barriers)
	}
}
