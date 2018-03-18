/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package utils

import (
	"fmt"
	"strings"
)

//Write by zhangtao<ztao8607@gmail.com> . In 2018/3/8.

// parseName 通过git地址解析工程名称
func ParseName(url string) (name string) {
	gitName := strings.Split(url, "/")
	name = strings.Split(gitName[len(gitName)-1], ".")[0]
	fmt.Printf("GitAgent Will Clone [%s]\n", name)
	return
}

// ParsePath 通过git地址解析出clone后的路径
// 例如通过https://github.com/andy-zhangtao/humCICD.git提取 andy-zhangtao/humCICD
func ParsePath(url string) (path string) {
	if strings.HasPrefix(url, "https://") {
		path = url[len("https://") : len(url)-len(".git")]
	} else {
		path = url[len("http://") : len(url)-len(".git")]
	}

	return
}
