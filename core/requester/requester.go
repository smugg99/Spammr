package requester

import (
	"compress/flate"
	"compress/gzip"
	"encoding/json"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"

	"smuggr.xyz/spammr/common/configurator"
	"smuggr.xyz/spammr/common/logger"

	"github.com/google/brotli/go/cbrotli"
)

var Logger = logger.NewCustomLogger("requ")
var RequestTemplates map[string]RequestTemplate

func Initialize(cmdFlags *configurator.CmdFlags) {
	Logger.Log(logger.MsgInitializing)

	if requestTemplates, err := LoadRequestTemplates(cmdFlags); err != nil {
		Logger.Errorf("failed to load request templates: %v", err.Error())
	} else {
		RequestTemplates = requestTemplates
	}

	Logger.Log(logger.MsgInitialized)

	SendRequests(*cmdFlags)
}

func SendRequests(cmdFlags configurator.CmdFlags) {
    var wg sync.WaitGroup

    for _, template := range RequestTemplates {
        wg.Add(1)
        go func(template RequestTemplate) {
            defer wg.Done()
			
			if generatedRequest, err := GenerateRequest(&template, cmdFlags); err != nil {
				Logger.Printf("failed to generate request template: %v", err.Error())
			} else {
				if _, err := SendRequest(generatedRequest); err != nil {
					Logger.Printf("failed to send request template: %v", err.Error())
				}
			}
        }(template)
    }

    wg.Wait()
}

func LoadRequestTemplates(cmdFlags *configurator.CmdFlags) (map[string]RequestTemplate, error) {
	Logger.Infof("loading request templates from: %s", cmdFlags.RequestsDir)

	templates := make(map[string]RequestTemplate)

	if cmdFlags.RequestsDir == "" {
		cmdFlags.RequestsDir = os.Getenv("REQUEST_TEMPLATES_DIRECTORY")
	}

	files, err := os.ReadDir(cmdFlags.RequestsDir)
	if err != nil {
		Logger.Errorf("failed to read directory: %v", err.Error())
		return nil, err
	}

	if len(files) == 0 {
		Logger.Warn(logger.ErrResourcesDirectoryEmpty.Format(cmdFlags.RequestsDir))
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == os.Getenv("REQUEST_TEMPLATE_EXTENSION") {
			Logger.Info(logger.MsgLoadingResource.Format(file.Name(), logger.ResourceRequestTemplate))

			filePath := filepath.Join(cmdFlags.RequestsDir, file.Name())
			data, err := os.ReadFile(filePath)
			if err != nil {
				Logger.Error(err.Error())
				return nil, err
			}

			var template RequestTemplate
			err = json.Unmarshal(data, &template)
			if err != nil {
				Logger.Errorf("failed to unmarshal JSON: %v", err.Error())
				return nil, err
			}

			filename := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			templates[filename] = template

			Logger.Info(logger.MsgResourceLoaded.Format(file.Name(), logger.ResourceRequestTemplate))
		}
	}

	return templates, nil
}

func GenerateRequest(template *RequestTemplate, cmdFlags configurator.CmdFlags) (*RequestTemplate, error) {
	replacements := GenerateReplacements(template.Wants, cmdFlags)

	for key, value := range template.Headers {
		template.Headers[key] = ReplacePlaceholders(value, replacements)
	}

	template.Body = ReplacePlaceholders(template.Body, replacements)

	Logger.Debug("generated request body:", template.Body)

	return template, nil
}

func SendRequest(template *RequestTemplate) (*http.Response, error) {
	client := &http.Client{}

	Logger.Infof("sending request to: %s", template.URL)

	req, err := http.NewRequest(template.Method, template.URL, strings.NewReader(template.Body))
	if err != nil {
		Logger.Errorf("failed to create request: %v", err.Error())
		return nil, err
	}

	for key, value := range template.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		Logger.Errorf("failed to send request: %v", err.Error())
		return nil, err
	}

	CheckResponse(template, resp)

	return resp, nil
}

func CheckResponse(template *RequestTemplate, resp *http.Response) bool {
	defer resp.Body.Close()

	var reader io.Reader

	switch resp.Header.Get("Content-Encoding") {
	case "br":
		reader = cbrotli.NewReader(resp.Body)
		defer reader.(*cbrotli.Reader).Close()
	case "gzip":
		gzipReader, err := gzip.NewReader(resp.Body)
		if err != nil {
			Logger.Errorf("error creating gzip reader: %v", err)
			reader = resp.Body
		} else {
			defer gzipReader.Close()
			reader = gzipReader
		}
	case "deflate":
		reader = flate.NewReader(resp.Body)
	default:
		reader = resp.Body
	}

	decompressedBody, err := io.ReadAll(reader)
	if err != nil {
		Logger.Errorf("error reading response body: %v", err)
		return false
	}

	responseDetails := ResponseDetails{
		StatusCode: resp.StatusCode,
		Headers:    make(map[string]string),
		Body:       string(decompressedBody),
	}

	for key, value := range resp.Header {
		responseDetails.Headers[key] = strings.Join(value, ", ")
	}

	jsonResponse, err := json.Marshal(responseDetails)
	if err != nil {
		Logger.Errorf("error marshalling JSON: %v", err)
	}

	Logger.Debugf("response details: %s", jsonResponse)

	if responseDetails.StatusCode >= 400 {
		Logger.Errorf("request to %s failed with status code: %v", template.URL, responseDetails.StatusCode)
		return false
	} else {
		Logger.Successf("request to %s succeeded with status code: %v", template.URL, responseDetails.StatusCode)
		return true
	}
}
