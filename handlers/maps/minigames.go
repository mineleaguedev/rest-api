package maps

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"net/http"
	"strings"
)

func (h *Handler) MiniGameMapsGetHandler(c *gin.Context) {
	minigame := c.Param("minigame")

	contents, err := h.services.GetMiniGameMapsList(minigame)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrS3GettingMiniGameMapsList, err)
		return
	}

	if len(contents) == 0 {
		h.services.HandleErr(c, http.StatusBadRequest, errors.ErrS3EmptyMiniGameMapsList)
		return
	}

	var formatsList []models.Format
	for _, key := range contents {
		foldersList := *key.Key

		folders := strings.Split(strings.TrimSuffix(foldersList, "/"), "/")
		for index, folder := range folders {
			if index == 1 {
				formatName := folder

				var isCanAdd = true
				for _, format := range formatsList {
					if format.Format != formatName {
						continue
					}
					isCanAdd = false
					break
				}

				if isCanAdd || len(formatsList) == 0 {
					formatsList = append(formatsList, models.Format{
						Format: formatName,
						Maps:   nil,
					})
				}
			} else if index == 2 {
				formatName := folders[1]
				mapName := folder

				for formatIndex, format := range formatsList {
					if format.Format != formatName {
						continue
					}

					var isCanAdd = true
					for _, minigameMap := range format.Maps {
						if minigameMap.Name != mapName {
							continue
						}
						isCanAdd = false
						break
					}

					if isCanAdd || len(format.Maps) == 0 {
						format.Maps = append(format.Maps, models.Map{
							Name:     mapName,
							Versions: nil,
						})
						formatsList[formatIndex] = format
					}
					break
				}
			} else if index == 3 {
				formatName := folders[1]
				mapName := folders[2]
				version := folder

				for formatIndex, format := range formatsList {
					if format.Format != formatName {
						continue
					}

					for mapIndex, minigameMap := range format.Maps {
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
							formatsList[formatIndex].Maps[mapIndex].Versions = minigameMap.Versions
						}
						break
					}
					break
				}
			}
		}
	}

	c.JSON(http.StatusOK, models.MiniGameMapsResponse{
		Success: true,
		Formats: formatsList,
	})
}
