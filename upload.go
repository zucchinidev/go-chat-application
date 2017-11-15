package main

import (
	"net/http"
	"io"
	"io/ioutil"
	"path"
)

func uploaderHandler(res http.ResponseWriter, req *http.Request) {
	userId := req.FormValue("userId")
	file, header, err := req.FormFile("avatarFile")
	if err != nil {
		io.WriteString(res, err.Error())
		return
	}

	data, err := ioutil.ReadAll(file)
	if err != nil {
		io.WriteString(res, err.Error())
		return
	}
	filename := path.Join("avatars", userId + path.Ext(header.Filename))
	err = ioutil.WriteFile(filename, data, 0777)
	if err != nil {
		io.WriteString(res, err.Error())
		return
	}
	io.WriteString(res, "Successful")
}
