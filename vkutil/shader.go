package vkutil

import (
	"io"
	"os"

	"code.witches.io/go/vulkan"
)

type Shader struct {
	Handle vulkan.ShaderModule
	Flag   vulkan.ShaderStageFlagBits

	device vulkan.Device
}

func NewShader(device vulkan.Device, path string, flag vulkan.ShaderStageFlagBits) (*Shader, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return NewShaderFromReader(device, f, flag)
}

func NewShaderFromReader(device vulkan.Device, r io.Reader, flag vulkan.ShaderStageFlagBits) (*Shader, error) {
	code, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}

	handle, err := vulkan.CreateShaderModule(device, vulkan.ShaderModuleCreateInfo{
		Type: vulkan.StructureTypeShaderModuleCreateInfo,
		Code: code,
	}, nil)
	if err != nil {
		return nil, err
	}

	return &Shader{
		Handle: handle,
		Flag:   flag,
		device: device,
	}, nil
}

func (s Shader) Close() {
	vulkan.DestroyShaderModule(s.device, s.Handle, nil)
}

func (s Shader) Stage() vulkan.PipelineShaderStageCreateInfo {
	return vulkan.PipelineShaderStageCreateInfo{
		Type:               vulkan.StructureTypePipelineShaderStageCreateInfo,
		Stage:              s.Flag,
		Module:             s.Handle,
		Name:               "main",
		SpecializationInfo: vulkan.SpecializationInfo{},
	}
}
