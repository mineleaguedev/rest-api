package maps

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"net/http"
	"strings"
)

func (h *Handler) MapsGetHandler(c *gin.Context) {
	contents, err := h.services.GetMapsList()
	if err != nil {
		h.services.HandleErr(c, http.StatusInternalServerError, errors.ErrS3GettingMapsList)
		return
	}

	var minigamesList []models.MiniGames
	for _, key := range contents {
		foldersList := *key.Key

		folders := strings.Split(strings.TrimSuffix(foldersList, "/"), "/")
		for index, folder := range folders {
			if index == 0 {
				minigameName := folder

				var isCanAdd bool
				for _, minigame := range minigamesList {
					if minigame.Name == minigameName {
						continue
					}
					isCanAdd = true
					break
				}

				if isCanAdd || len(minigamesList) == 0 {
					minigamesList = append(minigamesList, models.MiniGames{
						Name:   minigameName,
						Format: nil,
					})
				}
			} else if index == 1 {
				minigameName := folders[0]
				formatName := folder

				var isCanAdd bool
				for minigameIndex, minigame := range minigamesList {
					if minigame.Name != minigameName {
						continue
					}

					for _, format := range minigame.Format {
						if format.Format == formatName {
							continue
						}
						isCanAdd = true
						break
					}

					if isCanAdd || len(minigame.Format) == 0 {
						minigame.Format = append(minigame.Format, models.Format{
							Format: formatName,
							Map:    nil,
						})
						minigamesList[minigameIndex] = minigame
					}
					break
				}
			} else if index == 2 {
				minigameName := folders[0]
				formatName := folders[1]
				mapName := folder

				for minigameIndex, minigame := range minigamesList {
					if minigame.Name != minigameName {
						continue
					}

					for formatIndex, format := range minigame.Format {
						if format.Format != formatName {
							continue
						}

						var isCanAdd bool
						for _, minigameMap := range format.Map {
							if minigameMap.Name == mapName {
								continue
							}
							isCanAdd = true
							break
						}

						if isCanAdd || len(format.Map) == 0 {
							format.Map = append(format.Map, models.Map{
								Name:     mapName,
								Versions: nil,
							})
							minigamesList[minigameIndex].Format[formatIndex] = format
						}
						break
					}
					break
				}
			} else if index == 3 {
				minigameName := folders[0]
				formatName := folders[1]
				mapName := folders[2]
				version := folder

				for minigameIndex, minigame := range minigamesList {
					if minigame.Name != minigameName {
						continue
					}

					for formatIndex, format := range minigame.Format {
						if format.Format != formatName {
							continue
						}

						for mapIndex, minigameMap := range format.Map {
							if minigameMap.Name != mapName {
								continue
							}

							var isCanAdd bool
							for _, ver := range minigameMap.Versions {
								if ver == version {
									continue
								}
								isCanAdd = true
								break
							}

							if isCanAdd || len(minigameMap.Versions) == 0 {
								minigameMap.Versions = append(minigameMap.Versions, version)
								minigamesList[minigameIndex].Format[formatIndex].Map[mapIndex] = minigameMap
							}
							break
						}
						break
					}
					break
				}
			}
		}
	}

	c.JSON(http.StatusOK, models.MapsResponse{
		Success:   true,
		MiniGames: minigamesList,
	})
}
