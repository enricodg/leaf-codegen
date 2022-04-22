package handler

import (
	"leaf-codegen/helper"
	"leaf-codegen/helper/directory"
	"os"
	"path/filepath"
	"strings"
)

func (h handler) Init(project string) error {
	var splitProjectURL = strings.Split(project, "/")
	projectName := splitProjectURL[len(splitProjectURL)-1]

	mainPath := filepath.Join(directory.Cmd, projectName)
	if err := os.MkdirAll(mainPath, os.ModePerm); err != nil {
		return err
	}

	inboundPath := filepath.Join(directory.Internal, directory.Inbound)
	healthPath := filepath.Join(inboundPath, directory.Http, directory.Health)
	if err := os.MkdirAll(healthPath, os.ModePerm); err != nil {
		return err
	}
	outboundPath := filepath.Join(directory.Internal, directory.Outbound)
	if err := os.MkdirAll(outboundPath, os.ModePerm); err != nil {
		return err
	}
	useCasesPath := filepath.Join(directory.Internal, directory.UseCases)
	if err := os.MkdirAll(useCasesPath, os.ModePerm); err != nil {
		return err
	}

	configPath := filepath.Join(directory.Pkg, directory.Config)
	if err := os.MkdirAll(configPath, os.ModePerm); err != nil {
		return err
	}

	diPath := filepath.Join(directory.Pkg, directory.Di)
	if err := os.MkdirAll(diPath, os.ModePerm); err != nil {
		return err
	}

	resourcePath := filepath.Join(directory.Pkg, directory.Resource)
	injectionPath := filepath.Join(resourcePath, directory.Injection)
	if err := os.MkdirAll(injectionPath, os.ModePerm); err != nil {
		return err
	}

	if err := helper.InitializeProject(helper.InitializeProjectRequestDTO{
		ProjectURL:   project,
		ProjectName:  projectName,
		MainPath:     mainPath,
		InboundPath:  inboundPath,
		ConfigPath:   configPath,
		ResourcePath: resourcePath,
		DiPath:       diPath,
		OutboundPath: outboundPath,
		UseCasesPath: useCasesPath,
	}); err != nil {
		h.log.StandardLogger().Errorf("[%s] error initializing project: %+v", project, err.Error())
		return err
	}

	h.log.StandardLogger().Infof("[%s] finish initializing project", project)
	return nil
}
