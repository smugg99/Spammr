package requester

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"

	"smuggr.xyz/wkurwiaczo-inator/common/logger"
)

var Logger = logger.NewCustomLogger("requ")

func Initialize() {
	Logger.Log(logger.MsgInitializing)

	LoadRequestTemplates()
}

func LoadRequestTemplates() (map[string]RequestTemplate, error) {
	templates := make(map[string]RequestTemplate)
	directory := os.Getenv("REQUEST_TEMPLATES_DIRECTORY")

	files, err := os.ReadDir(directory)
	if err != nil {
		Logger.Errorf(err.Error())
		return nil, err
	}

	if len(files) == 0 {
		Logger.Warn(logger.ErrResourcesDirectoryEmpty.Format(directory))
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == os.Getenv("REQUEST_TEMPLATE_EXTENSION") {
			Logger.Info(logger.MsgLoadingResource.Format(file.Name(), logger.ResourceRequestTemplate))
			
			filePath := filepath.Join(directory, file.Name())
			data, err := os.ReadFile(filePath)
			if err != nil {
				Logger.Errorf(err.Error())
				return nil, err
			}

			var template RequestTemplate
			err = json.Unmarshal(data, &template)
			if err != nil {
				Logger.Errorf(err.Error())
				return nil, err
			}

			filename := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			templates[filename] = template

			Logger.Info(logger.MsgResourceLoaded.Format(file.Name(), logger.ResourceRequestTemplate))
		}
	}

	return templates, nil
}
