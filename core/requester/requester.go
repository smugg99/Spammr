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

	"smuggr.xyz/wkurwiaczo-inator/common/logger"

	"github.com/google/brotli/go/cbrotli"
)

var Logger = logger.NewCustomLogger("requ")
var RequestTemplates map[string]RequestTemplate

func Initialize() {
	Logger.Log(logger.MsgInitializing)

	if requestTemplates, err := LoadRequestTemplates(); err != nil {
		Logger.Errorf("Failed to load request templates: %v", err.Error())
	} else {
		RequestTemplates = requestTemplates
	}

	templateName := "alians.oze.pl"
	if template, ok := RequestTemplates[templateName]; ok {
		if _, err := SendRequest(&template); err != nil {
			Logger.Errorf("Failed to send request template: %v", err.Error())
		}
	} else {
		Logger.Log(logger.ErrResourceNotFound.Format(templateName, logger.ResourceRequestTemplate))
	}

	Logger.Log(logger.MsgInitialized)
}

func LoadRequestTemplates() (map[string]RequestTemplate, error) {
	templates := make(map[string]RequestTemplate)
	directory := os.Getenv("REQUEST_TEMPLATES_DIRECTORY")

	files, err := os.ReadDir(directory)
	if err != nil {
		Logger.Errorf("Failed to read directory: %v", err.Error())
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
				Logger.Error(err.Error())
				return nil, err
			}

			var template RequestTemplate
			err = json.Unmarshal(data, &template)
			if err != nil {
				Logger.Errorf("Failed to unmarshal JSON: %v", err.Error())
				return nil, err
			}

			filename := strings.TrimSuffix(file.Name(), filepath.Ext(file.Name()))
			templates[filename] = template

			Logger.Info(logger.MsgResourceLoaded.Format(file.Name(), logger.ResourceRequestTemplate))
		}
	}

	return templates, nil
}

func GenerateRequest(template *RequestTemplate, has map[Want]string) (*RequestTemplate, error) {
	for key, value := range template.Headers {
		template.Headers[key] = ReplacePlaceholders(value, template.Wants, has)
	}

	template.Body = ReplacePlaceholders(template.Body, template.Wants, has)

	Logger.Debug("Generated Request Body:", template.Body)

	return template, nil
}

func SendRequest(template *RequestTemplate) (*http.Response, error) {
	client := &http.Client{}

	req, err := http.NewRequest(template.Method, template.URL, strings.NewReader(template.Body))
	if err != nil {
		Logger.Errorf("Failed to create request: %v", err.Error())
		return nil, err
	}

	for key, value := range template.Headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		Logger.Errorf("Failed to send request: %v", err.Error())
		return nil, err
	}

	LogResponse(resp)

	return resp, nil
}

func LogResponse(resp *http.Response) {
    defer resp.Body.Close()

    var reader io.Reader

    switch resp.Header.Get("Content-Encoding") {
    case "br":
        reader = cbrotli.NewReader(resp.Body)
        defer reader.(*cbrotli.Reader).Close()
    case "gzip":
        gzipReader, err := gzip.NewReader(resp.Body)
        if err != nil {
            Logger.Errorf("Error creating gzip reader: %v", err)
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
        Logger.Errorf("Error reading response body: %v", err)
        return
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
        Logger.Errorf("Error marshalling JSON: %v", err)
        return
    }

    Logger.Debugf("Response Details: %s", jsonResponse)
}