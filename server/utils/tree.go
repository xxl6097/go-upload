package utils

import (
	"fmt"
	"os"
	"path/filepath"
)

// Node 表示目录树中的一个节点
type Node struct {
	Title    string `json:"title"`
	Path     string `json:"path"`
	Spread   bool   `json:"spread"`
	Children []Node `json:"children,omitempty"`
}

// buildTree 是递归构建目录树的函数
func buildTree(path string) (Node, error) {
	root := Node{Title: filepath.Base(path)}
	info, err := os.Stat(path)
	if err != nil {
		return root, err
	}
	root.Path = path
	root.Spread = true

	if !info.IsDir() {
		return root, nil
	}

	entries, err := os.ReadDir(path)
	if err != nil {
		return root, err
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		childPath := filepath.Join(path, entry.Name())
		childNode, err := buildTree(childPath)
		if err != nil {
			return root, err
		}
		root.Children = append(root.Children, childNode)
	}

	return root, nil
}

func GetDirJsonTree(rootDir string) *Node {
	//rootDir := "./files" // 可以将根目录修改为你想要遍历的目录
	tree, err := buildTree(rootDir)
	if err != nil {
		fmt.Printf("Error building tree: %v\n", err)
		return nil
	}
	return &tree

	//jsonData, err := json.MarshalIndent(tree, "", "  ")
	//if err != nil {
	//	fmt.Printf("Error marshalling to JSON: %v\n", err)
	//	return nil
	//}
	//
	//return jsonData
}
