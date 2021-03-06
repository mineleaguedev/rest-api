package plugins

import (
	"github.com/gin-gonic/gin"
	"github.com/mineleaguedev/rest-api/errors"
	"github.com/mineleaguedev/rest-api/models"
	"net/http"
	"strings"
)

func (h *Handler) PluginsGetHandler(c *gin.Context) {
	contents, err := h.services.GetPluginsList()
	if err != nil {
		h.services.HandleErr(c, http.StatusInternalServerError, errors.ErrS3GettingPluginsList)
		return
	}

	var pluginsList []models.Plugin
	for _, key := range contents {
		foldersList := strings.ReplaceAll(*key.Key, "plugins/", "")
		if foldersList == "" {
			continue
		}

		folders := strings.Split(strings.TrimSuffix(foldersList, "/"), "/")
		for index, folder := range folders {
			if index == 0 {
				pluginName := folder

				var isCanAdd = true
				for _, plugin := range pluginsList {
					if plugin.Name != pluginName {
						continue
					}
					isCanAdd = false
					break
				}

				if isCanAdd || len(pluginsList) == 0 {
					pluginsList = append(pluginsList, models.Plugin{
						Name:     pluginName,
						Versions: nil,
					})
				}
			} else if index == 1 {
				pluginName := folders[0]
				version := folder

				for pluginIndex, plugin := range pluginsList {
					if plugin.Name != pluginName {
						continue
					}

					var isCanAdd = true
					for _, ver := range plugin.Versions {
						if ver != version {
							continue
						}
						isCanAdd = false
						break
					}

					if isCanAdd || len(plugin.Versions) == 0 {
						plugin.Versions = append(plugin.Versions, version)
						pluginsList[pluginIndex] = plugin
					}
					break
				}
			}
		}
	}

	c.JSON(http.StatusOK, models.PluginsResponse{
		Success: true,
		Plugins: pluginsList,
	})
}
