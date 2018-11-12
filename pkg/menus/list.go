package menus

import (
	"bytes"

	"github.com/isard-vdi/isard-ipxe/pkg/client/list"
	"github.com/isard-vdi/isard-ipxe/pkg/client/mocks"
	"github.com/isard-vdi/isard-ipxe/pkg/config"
)

// GenerateList generates an iPXE menu with the VM list
func GenerateList(webRequest mocks.WebRequest, token string, username string) (string, error) {
	config := config.Config{}
	err := config.ReadConfig()
	if err != nil {
		buf := new(bytes.Buffer)

		t := parseTemplate("error.ipxe")
		t.Execute(buf, menuTemplateData{
			Err: "reading the configuration file",
		})

		return buf.String(), err
	}

	vms, err := list.Call(webRequest, token)
	if err != nil {
		if err.Error() == "HTTP Code: 403" {
			buf := new(bytes.Buffer)

			t := parseTemplate("login.ipxe")
			t.Execute(buf, menuTemplateData{
				BaseURL: config.BaseURL,
			})

			return buf.String(), err
		}

		buf := new(bytes.Buffer)

		t := parseTemplate("error.ipxe")
		t.Execute(buf, menuTemplateData{
			Err: "calling the API",
		})

		return buf.String(), err
	}

	buf := new(bytes.Buffer)

	t := parseTemplate("VMList.ipxe")
	t.Execute(buf, menuTemplateData{
		BaseURL:  config.BaseURL,
		Token:    token,
		Username: username,
		VMs:      vms.VMs,
	})

	return buf.String(), err
}
