package helper

import (
	"leaf-codegen/helper/directory"
	"leaf-codegen/helper/templates"
	"os"
	"path/filepath"
	"text/template"
)

type (
	InitializeProjectRequestDTO struct {
		ProjectURL   string
		ProjectName  string
		MainPath     string
		InboundPath  string
		ConfigPath   string
		ResourcePath string
		DiPath       string
		OutboundPath string
		UseCasesPath string
	}
)

func InitializeProject(request InitializeProjectRequestDTO) error {
	if err := createFileByTemplate(".gitignore", templates.GitIgnoreTemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate("go.mod", templates.GoModTemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate(".env.example", templates.EnvTemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate("generateMock.sh", templates.GenerateMockTemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate(filepath.Join(request.ConfigPath, "configApp.go"), templates.ConfigAppTemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate(filepath.Join(request.ConfigPath, "configSentry.go"), templates.ConfigSentryTemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate(filepath.Join(request.ConfigPath, "configNewRelic.go"), templates.ConfigNewRelicTemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate(filepath.Join(request.ConfigPath, "di.go"), templates.ConfigDITemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate(filepath.Join(request.ResourcePath, directory.Injection, "logger.go"), templates.LoggerTemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate(filepath.Join(request.ResourcePath, directory.Injection, "tracer.go"), templates.TracerTemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate(filepath.Join(request.ResourcePath, directory.Injection, "translator.go"), templates.TranslatorTemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate(filepath.Join(request.ResourcePath, directory.Injection, "validator.go"), templates.ValidatorTemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate(filepath.Join(request.ResourcePath, "resource.go"), templates.ResourceTemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate(filepath.Join(request.ResourcePath, "di.go"), templates.ResourceDiTemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate(filepath.Join(request.DiPath, "di.go"), templates.DITemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate(filepath.Join(request.InboundPath, directory.Http, directory.Health, "controller.go"), templates.HealthControllerTemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate(filepath.Join(request.InboundPath, directory.Http, directory.Health, "check.go"), templates.HealthCheckTemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate(filepath.Join(request.InboundPath, directory.Http, directory.Health, "routes.go"), templates.HealthRoutesTemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate(filepath.Join(request.InboundPath, directory.Http, "routes.go"), templates.HttpRoutesTemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate(filepath.Join(request.InboundPath, "di.go"), templates.InboundDITemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate(filepath.Join(request.OutboundPath, "di.go"), templates.OutboundDITemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate(filepath.Join(request.UseCasesPath, "di.go"), templates.UseCasesDITemplate, request); err != nil {
		return err
	}
	if err := createFileByTemplate(filepath.Join(request.MainPath, "main.go"), templates.MainTemplate, request); err != nil {
		return err
	}

	return nil
}

func createFileByTemplate(outputPath string, fileTemplate string, request interface{}) error {
	t, err := template.New(outputPath).Parse(fileTemplate)
	if err != nil {
		return err
	}

	file, err := os.Create(outputPath)
	if err != nil {
		return err
	}
	defer file.Close()

	err = t.Execute(file, request)
	if err != nil {
		return err
	}

	return nil
}
