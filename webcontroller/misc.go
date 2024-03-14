package webcontroller

import (
	"encoding/base64"
	"fmt"
	"net/http"

	"fornaxian.tech/log"
	"github.com/julienschmidt/httprouter"
)

// ServeFileViewerDemo is a dummy API response that responds with info about a
// non-existent demo file. This is required by the a-ads ad network to allow for
// automatic checking of the presence of the ad unit on this page.
func (wc *WebController) serveShareXConfig(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	templateData := wc.newTemplateData(w, r)

	w.Header().Add("Content-Disposition", "attachment; filename=pixeldrain.com.sxcu")
	if templateData.Authenticated {
		sess, err := templateData.PixelAPI.PostUserSession("sharex")
		if err != nil {
			log.Error("Failed to create user session: %s", err)
			wc.templates.Get().ExecuteTemplate(w, "500", templateData)
			return
		}

		w.Write([]byte(fmt.Sprintf(
			`{
	"Version": "12.4.1",
	"DestinationType": "ImageUploader, TextUploader, FileUploader",
	"RequestMethod": "POST",
	"RequestURL": "https://pixeldrain.com/api/file",
	"Headers": {
		"Authorization": "Basic %s"
	},
	"Body": "MultipartFormData",
	"FileFormName": "file",
	"URL": "https://pixeldrain.com/u/$json:id$",
	"ThumbnailURL": "https://pixeldrain.com/api/file/$json:id$/thumbnail"
}
`,
			base64.StdEncoding.EncodeToString([]byte(
				templateData.User.Username+":"+sess.AuthKey.String(),
			)),
		)))
	} else {
		w.Write([]byte(
			`{
	"Version": "12.4.1",
	"DestinationType": "ImageUploader, TextUploader, FileUploader",
	"RequestMethod": "POST",
	"RequestURL": "https://pixeldrain.com/api/file",
	"Body": "MultipartFormData",
	"FileFormName": "file",
	"URL": "https://pixeldrain.com/u/$json:id$",
	"ThumbnailURL": "https://pixeldrain.com/api/file/$json:id$/thumbnail"
}
`,
		))
	}
}
