package plugins

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"net/http"
	"strings"
)

func (h *Handler) PluginVersionsGetHandler(c *gin.Context) {
	plugin := c.Param("name")

	contents, err := h.services.GetPluginVersionsList(plugin)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrS3GettingPluginVersionsList, err)
		return
	}

	if len(contents) == 0 {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrS3EmptyPluginVersionsList)
		return
	}

	var versionsList []string
	for _, key := range contents {
		foldersList := *key.Key

		folders := strings.Split(strings.TrimSuffix(foldersList, "/"), "/")
		for index, folder := range folders {
			if index == 1 {
				version := folder

				var isCanAdd bool
				for _, ver := range versionsList {
					if ver == version {
						continue
					}
					isCanAdd = true
					break
				}

				if isCanAdd || len(versionsList) == 0 {
					versionsList = append(versionsList, version)
				}
			}
		}
	}

	c.JSON(http.StatusOK, models.PluginResponse{
		Success:  true,
		Versions: versionsList,
	})
}
