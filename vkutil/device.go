package vkutil

import (
	"fmt"
	"sort"

	"code.witches.io/go/vulkan"
)

func Info(instance vulkan.Instance) error {
	ls, err := vulkan.EnumerateInstanceLayerProperties()
	if err == nil {
		for _, l := range ls {
			fmt.Println(l)
			e, err := vulkan.EnumerateInstanceExtensionProperties(l.LayerName.String())
			if err == nil {
				fmt.Println(e)
			}
		}
	}

	es, err := vulkan.EnumerateInstanceExtensionProperties("")
	if err == nil {
		sort.Slice(es, func(i, j int) bool {
			return es[i].ExtensionName.String() < es[j].ExtensionName.String()
		})
		for _, e := range es {
			fmt.Println("extension:", e)
		}
	}

	pgs, err := vulkan.EnumeratePhysicalDeviceGroups(instance)
	if err != nil {
		return err
	}

	for _, pg := range pgs {
		if pg.Type != vulkan.StructureTypePhysicalDeviceGroupProperties {
			continue
		}

		for _, p := range pg.PhysicalDevices {
			pp := vulkan.GetPhysicalDeviceProperties2(p)
			fmt.Println(pp)

			mp := vulkan.GetPhysicalDeviceMemoryProperties(p)
			fmt.Println(mp)

			qp := vulkan.GetPhysicalDeviceQueueFamilyProperties2(p)
			fmt.Println(qp)

			ls, err := vulkan.EnumerateDeviceLayerProperties(p)
			if err == nil {
				for _, l := range ls {
					fmt.Println(l)
					e, err := vulkan.EnumerateDeviceExtensionProperties(p, l.LayerName.String())
					if err == nil {
						fmt.Println(e)
					}
				}
			}

			es, err := vulkan.EnumerateDeviceExtensionProperties(p, "")
			if err == nil {
				sort.Slice(es, func(i, j int) bool {
					return es[i].ExtensionName.String() < es[j].ExtensionName.String()
				})
				for _, e := range es {
					fmt.Println("extension:", e)
				}
			}

			// todo vulkan.GetPhysicalDeviceFeatures2(p)
		}
	}
	return nil
}
