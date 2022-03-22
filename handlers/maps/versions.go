package maps

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"net/http"
	"strings"
)

func (h *Handler) MapVersionsGetHandler(c *gin.Context) {
	minigame := c.Param("minigame")
	format := c.Param("format")
	mapName := c.Param("map")

	contents, err := h.services.GetMiniGameFormatMapVersionsList(minigame, format, mapName)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrS3GettingMiniGameFormatMapVersionsList, err)
		return
	}

	if len(contents) == 0 {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrS3EmptyMiniGameFormatMapVersionsList)
		return
	}

	var versionsList []string
	for _, key := range contents {
		foldersList := *key.Key

		folders := strings.Split(strings.TrimSuffix(foldersList, "/"), "/")
		for index, folder := range folders {
			if index == 3 {
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

	c.JSON(http.StatusOK, models.MiniGameFormatMapVersionsResponse{
		Success:  true,
		Versions: versionsList,
	})
}
