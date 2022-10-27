package http

import (
	"encoding/base64"
	"io"
	"net/http"
	"phantom/global"
	"phantom/util"
	"strconv"
	"strings"
)

func HandlePFP(w http.ResponseWriter, r *http.Request) {
	split := strings.Split(r.URL.Path, "=")
	image := strings.Split(split[len(split)-1], ".")
	id, err := strconv.Atoi(image[0])

	if err != nil {
		io.WriteString(w, string("err1"))
		util.Error("Error converting id to an integer")
		return
	}

	acc, err2 := global.GetUserDataFromUserId(id)
	upl, err4 := global.GetUploadDataFromUserId(id)

	if !err4 {
		io.WriteString(w, string("err4"))
		util.Error("error getting user object from database (upload)")
		return
	}

	if !err2 {
		io.WriteString(w, string("err2"))
		util.Error("error getting user object from database (userdata)")
		return
	}

	res, err := base64.StdEncoding.DecodeString(upl.Avatar)
	if err != nil {
		io.WriteString(w, string("err3"))
		util.Error("Error decoding avatar")
		return
	}

	util.Debug("WebAPI", "Provided avatar for user %s", acc.Email)
	w.Write(res)
}
