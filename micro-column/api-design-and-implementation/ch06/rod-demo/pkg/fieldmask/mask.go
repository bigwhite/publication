package fieldmask

import (
	"encoding/json"
	"strings"
)

// Prune 根据掩码列表裁剪数据
// data: 原始数据 (通常是 struct 或 struct 切片)
// fields: 字段列表，如 ["id", "name", "info.email"]
func Prune(data interface{}, fields []string) (interface{}, error) {
	// 如果没有指定 fields，直接返回原数据
	if len(fields) == 0 {
		return data, nil
	}

	// 1. 将 struct 序列化再反序列化为 map[string]interface{}
	// 这是一个"偷懒"但稳健的做法，利用 encoding/json 处理 tag 映射
	// 生产环境高性能场景可以考虑使用 reflect 直接遍历，或使用 mapstructure
	var temp interface{}
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(bytes, &temp); err != nil {
		return nil, err
	}

	// 2. 构建掩码树
	maskTree := buildMaskTree(fields)

	// 3. 递归裁剪
	return pruneRecursive(temp, maskTree), nil
}

// 掩码树结构：map[string]interface{}
// key 是字段名，value 是子掩码树（如果是叶子节点则为 nil）
func buildMaskTree(fields []string) map[string]interface{} {
	tree := make(map[string]interface{})
	for _, field := range fields {
		parts := strings.Split(field, ".")
		current := tree
		for i, part := range parts {
			part = strings.TrimSpace(part)
			if i == len(parts)-1 {
				// 叶子节点
				if _, exists := current[part]; !exists {
					current[part] = nil
				}
			} else {
				// 中间节点
				if _, exists := current[part]; !exists {
					current[part] = make(map[string]interface{})
				}
				// 如果已经是 nil (意味着之前作为叶子节点出现过，如 "a", "a.b")
				// 则覆盖为 map，以支持更深层级
				if current[part] == nil {
					current[part] = make(map[string]interface{})
				}
				current = current[part].(map[string]interface{})
			}
		}
	}
	return tree
}

func pruneRecursive(current interface{}, mask map[string]interface{}) interface{} {
	if current == nil {
		return nil
	}

	// 处理切片/数组
	if slice, ok := current.([]interface{}); ok {
		newSlice := make([]interface{}, len(slice))
		for i, item := range slice {
			newSlice[i] = pruneRecursive(item, mask)
		}
		return newSlice
	}

	// 处理对象 (Map)
	if obj, ok := current.(map[string]interface{}); ok {
		newObj := make(map[string]interface{})
		for k, v := range obj {
			// 检查该字段是否在掩码中
			if subMask, exists := mask[k]; exists {
				if subMask == nil {
					// 叶子节点，保留该字段的全部值
					newObj[k] = v
				} else {
					// 还有子掩码，递归处理
					newObj[k] = pruneRecursive(v, subMask.(map[string]interface{}))
				}
			}
		}
		return newObj
	}

	// 基本类型，直接返回
	return current
}
