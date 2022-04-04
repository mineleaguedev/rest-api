package maps

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"net/http"
	"strings"
)

func (h *Handler) MiniGameFormatMapsGetHandler(c *gin.Context) {
	minigame := c.Param("minigame")
	format := c.Param("format")

	contents, err := h.services.GetMiniGameFormatMapsList(minigame, format)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrS3GettingMiniGameFormatMapsList, err)
		return
	}

	if len(contents) == 0 {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrS3EmptyMiniGameFormatMapsList)
		return
	}

	var mapsList []models.Map
	for _, key := range contents {
		foldersList := strings.ReplaceAll(*key.Key, "maps/", "")
		if foldersList == "" {
			continue
		}

		folders := strings.Split(strings.TrimSuffix(foldersList, "/"), "/")
		for index, folder := range folders {
			if index == 2 {
				mapName := folder

				var isCanAdd = true
				for _, minigameMap := range mapsList {
					if minigameMap.Name != mapName {
						continue
					}
					isCanAdd = false
					break
				}

				if isCanAdd || len(mapsList) == 0 {
					mapsList = append(mapsList, models.Map{
						Name:     mapName,
						Versions: nil,
					})
				}
			} else if index == 3 {
				mapName := folders[2]
				version := folder

				for mapIndex, minigameMap := range mapsList {
					if minigameMap.Name != mapName {
						continue
					}

					var isCanAdd = true
					for _, ver := range minigameMap.Versions {
						if ver != version {
							continue
						}
						isCanAdd = false
						break
					}

					if isCanAdd || len(minigameMap.Versions) == 0 {
						minigameMap.Versions = append(minigameMap.Versions, version)
						mapsList[mapIndex].Versions = minigameMap.Versions
					}
					break
				}
			}
		}
	}

	c.JSON(http.StatusOK, models.MiniGameFormatMapsResponse{
		Success: true,
		Maps:    mapsList,
	})
}
