/*
 * Copyright (c) 2018.
 * andy-zhangtao <ztao8607@gmail.com>
 */

package model

//Write by zhangtao<ztao8607@gmail.com> . In 2018/2/24.

type TagEventMsg struct {
	Kind   string `json:"kind"`
	GitURL string `json:"git_url"`
	Tag    string `json:"tag"`
	Branch string `json:"branch"`
}

type ErrEventMsg struct {
	Name string `json:"name"`
	Err  string `json:"err"`
}
