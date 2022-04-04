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
		foldersList := strings.ReplaceAll(*key.Key, "maps/", "")
		if foldersList == "" {
			continue
		}

		folders := strings.Split(strings.TrimSuffix(foldersList, "/"), "/")
		for index, folder := range folders {
			if index == 0 {
				minigameName := folder

				var isCanAdd = true
				for _, minigame := range minigamesList {
					if minigame.Name != minigameName {
						continue
					}
					isCanAdd = false
					break
				}

				if isCanAdd || len(minigamesList) == 0 {
					minigamesList = append(minigamesList, models.MiniGames{
						Name:    minigameName,
						Formats: nil,
					})
				}
			} else if index == 1 {
				minigameName := folders[0]
				formatName := folder

				for minigameIndex, minigame := range minigamesList {
					if minigame.Name != minigameName {
						continue
					}

					var isCanAdd = true
					for _, format := range minigame.Formats {
						if format.Format != formatName {
							continue
						}
						isCanAdd = false
						break
					}

					if isCanAdd || len(minigame.Formats) == 0 {
						minigame.Formats = append(minigame.Formats, models.Format{
							Format: formatName,
							Maps:   nil,
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

					for formatIndex, format := range minigame.Formats {
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
							minigamesList[minigameIndex].Formats[formatIndex] = format
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

					for formatIndex, format := range minigame.Formats {
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
								minigamesList[minigameIndex].Formats[formatIndex].Maps[mapIndex] = minigameMap
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
