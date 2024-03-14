package webcontroller

import (
	"fmt"
	"html/template"
	"net/http"

	"fornaxian.tech/log"
	"fornaxian.tech/pixeldrain_api_client/pixelapi"
)

func (wc *WebController) adminGlobalsForm(td *TemplateData, r *http.Request) (f Form) {
	if !td.Authenticated || !td.User.IsAdmin {
		return Form{Title: ";-)"}
	}

	f = Form{
		Name:        "admin_globals",
		Title:       "Pixeldrain global configuration",
		PreFormHTML: template.HTML("<p>Careful! The slightest typing error could bring the whole website down</p>"),
		SubmitLabel: "Submit",
	}

	globals, err := td.PixelAPI.AdminGetGlobals()
	if err != nil {
		f.SubmitMessages = []template.HTML{template.HTML(err.Error())}
		return f
	}
	var globalsMap = make(map[string]string)
	for _, v := range globals {
		f.Fields = append(f.Fields, Field{
			Name:         v.Key,
			DefaultValue: v.Value,
			Label:        v.Key,
			Type: func() FieldType {
				switch v.Key {
				case
					"email_address_change_body",
					"email_password_reset_body",
					"email_register_user_body":
					return FieldTypeTextarea
				case
					"api_ratelimit_limit",
					"api_ratelimit_rate",
					"cron_interval_seconds",
					"file_inactive_expiry_days",
					"max_file_size",
					"pixelstore_min_redundancy":
					return FieldTypeNumber
				default:
					return FieldTypeText
				}
			}(),
		})
		globalsMap[v.Key] = v.Value
	}

	if f.ReadInput(r) {
		var successfulUpdates = 0
		for k, v := range f.Fields {
			if v.EnteredValue == globalsMap[v.Name] {
				continue // Change changes, no need to update
			}

			// Value changed, try to update global setting
			if err = td.PixelAPI.AdminSetGlobals(v.Name, v.EnteredValue); err != nil {
				if apiErr, ok := err.(pixelapi.Error); ok {
					f.SubmitMessages = append(f.SubmitMessages, template.HTML(apiErr.Message))
				} else {
					log.Error("%s", err)
					f.SubmitMessages = append(f.SubmitMessages, template.HTML(
						fmt.Sprintf("Failed to set '%s': %s", v.Name, err),
					))
					return f
				}
			} else {
				f.Fields[k].DefaultValue = v.EnteredValue
				successfulUpdates++
			}

		}
		if len(f.SubmitMessages) == 0 {
			// Request was a success
			f.SubmitSuccess = true
			f.SubmitMessages = []template.HTML{template.HTML(
				fmt.Sprintf("Success! %d values updated", successfulUpdates),
			)}
		}
	}
	return f
}
