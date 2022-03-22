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
						Name:    minigameName,
						Formats: nil,
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

					for _, format := range minigame.Formats {
						if format.Format == formatName {
							continue
						}
						isCanAdd = true
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

						var isCanAdd bool
						for _, minigameMap := range format.Maps {
							if minigameMap.Name == mapName {
								continue
							}
							isCanAdd = true
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

				var isCanAdd bool
				for _, format := range formatsList {
					if format.Format == formatName {
						continue
					}
					isCanAdd = true
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

					var isCanAdd bool
					for _, minigameMap := range format.Maps {
						if minigameMap.Name == mapName {
							continue
						}
						isCanAdd = true
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
		foldersList := *key.Key

		folders := strings.Split(strings.TrimSuffix(foldersList, "/"), "/")
		for index, folder := range folders {
			if index == 2 {
				mapName := folder

				var isCanAdd bool
				for _, minigameMap := range mapsList {
					if minigameMap.Name == mapName {
						continue
					}
					isCanAdd = true
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

	c.JSON(http.StatusOK, models.MiniGameFormatMapVersionsResponse{
		Success:  true,
		Versions: versionsList,
	})
}

func (h *Handler) MapWorldGetHandler(c *gin.Context) {
	minigame := c.Param("minigame")
	format := c.Param("format")
	mapName := c.Param("map")
	version := c.Param("version")

	filePath, fileName, err := h.services.DownloadMapWorld(minigame, format, mapName, version)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrS3DownloadingMapWorld, err)
		return
	}

	c.FileAttachment(*filePath, *fileName)
}

func (h *Handler) MapConfigGetHandler(c *gin.Context) {
	minigame := c.Param("minigame")
	format := c.Param("format")
	mapName := c.Param("map")
	version := c.Param("version")

	filePath, fileName, err := h.services.DownloadMapConfig(minigame, format, mapName, version)
	if err != nil {
		h.services.HandleInternalErr(c, errors.ErrS3DownloadingMapConfig, err)
		return
	}

	c.FileAttachment(*filePath, *fileName)
}
