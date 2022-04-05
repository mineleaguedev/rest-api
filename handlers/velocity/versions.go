package velocity

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"net/http"
	"strings"
)

func (h *Handler) VelocityVersionsGetHandler(c *gin.Context) {
	contents, err := h.services.GetVelocityVersionList()
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrS3GettingVelocityVersionsList, err)
		return
	}

	if len(contents) == 0 {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrS3EmptyVelocityVersionsList)
		return
	}

	var versionsList []string
	for _, key := range contents {
		foldersList := strings.ReplaceAll(*key.Key, "velocity/", "")
		if foldersList == "" {
			continue
		}

		folders := strings.Split(strings.TrimSuffix(foldersList, "/"), "/")
		for index, folder := range folders {
			if index == 0 {
				version := folder

				var isCanAdd = true
				for _, ver := range versionsList {
					if ver != version {
						continue
					}
					isCanAdd = false
					break
				}

				if isCanAdd || len(versionsList) == 0 {
					versionsList = append(versionsList, version)
				}
			}
		}
	}

	c.JSON(http.StatusOK, models.VelocityResponse{
		Success:  true,
		Versions: versionsList,
	})
}
