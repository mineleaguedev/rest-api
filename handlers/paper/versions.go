package paper

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"net/http"
	"strings"
)

func (h *Handler) PaperVersionsGetHandler(c *gin.Context) {
	contents, err := h.services.GetPaperVersionList()
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrS3GettingPaperVersionsList, err)
		return
	}

	if len(contents) == 0 {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrS3EmptyPaperVersionsList)
		return
	}

	var versionsList []string
	for _, key := range contents {
		foldersList := strings.ReplaceAll(*key.Key, "paper/", "")
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

	c.JSON(http.StatusOK, models.PaperResponse{
		Success:  true,
		Versions: versionsList,
	})
}
