package main

import (
	"encoding/base64"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type Link struct {
	UUID      string   `yaml:"uuid"`
	VlessLinks []string `yaml:"vless_links"`
}

type LinkConfig struct {
	Links []Link `yaml:"links"`
}

var linkConfig LinkConfig

// 读取并解析 link.yml 文件
func loadLinks(filename string) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(data, &linkConfig)
	if err != nil {
		return err
	}
	return nil
}

// 根据 UUID 获取 VLESS 链接数组
func getVlessLinksByUUID(uuid string) ([]string, error) {
	for _, link := range linkConfig.Links {
		if link.UUID == uuid {
			return link.VlessLinks, nil
		}
	}
	return nil, fmt.Errorf("no links found for UUID: %s", uuid)
}

func linkHandler(w http.ResponseWriter, r *http.Request) {
	uuid := r.URL.Path[len("/sub/"):]
	if uuid == "" {
		http.Error(w, "UUID is required", http.StatusBadRequest)
		return
	}

	vlessLinks, err := getVlessLinksByUUID(uuid)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	// 将所有 VLESS 链接用换行符拼接成一个字符串
	linksContent := strings.Join(vlessLinks, "\n")

	// 对拼接后的内容进行 Base64 编码
	encodedContent := base64.StdEncoding.EncodeToString([]byte(linksContent))

	// 返回编码后的内容
	fmt.Fprintln(w, encodedContent)
}

func main() {
	// 从 link.yml 加载链接配置
	err := loadLinks("link.yml")
	if err != nil {
		log.Fatalf("Failed to load link.yml: %v", err)
	}

	// 设置路由
	http.HandleFunc("/sub/", linkHandler)

	// 启动服务器
	port := "8080"
	fmt.Printf("Server is running on port %s\n", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

