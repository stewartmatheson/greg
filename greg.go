package main

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v2"
)

type dockerComposeApplication struct {
	Image       string   `yaml:"image,omitempty"`
	Context     string   `yaml:"context,omitempty"`
	Ports       []string `yaml:"ports,omitempty"`
	Volumes     []string `yaml:"volumes,omitempty"`
	Environment []string `yaml:"environment,omitempty"`
	DependsOn   []string `yaml:"depends_on,omitempty"`
	Entrypoint  string   `yaml:"entrypoint,omitempty"`
}

type configurationType string

const (
	// Context ... Comment
	Context configurationType = "C"

	// Image ... Comment
	Image configurationType = "I"
)

func (d dockerComposeApplication) configuredFor() (configurationType, error) {
	if len(d.Image) > 0 {
		return Image, nil
	}

	if len(d.Context) > 0 {
		return Context, nil
	}

	return "", errors.New("Invalid Configuration")
}

type dockerCompose struct {
	Version  string `yaml:"version"`
	Services map[string]dockerComposeApplication
}

func readDockerCompose() dockerCompose {
	dockerComposeFile, err := ioutil.ReadFile("docker-compose.yml")

	if err != nil {
		log.Printf("Can not load applications file #%v ", err)
	}

	var dockerCompose dockerCompose
	err = yaml.Unmarshal(dockerComposeFile, &dockerCompose)

	if err != nil {
		log.Printf("Can not unmarshal docker compose YAML file #%v ", err)
	}

	return dockerCompose
}

type configurationLayerApplication struct {
	Context string `yaml:"context,omitempty"`
	Image   string `yaml:"image,omitempty"`
}

type configurationLayer struct {
	Applications map[string]configurationLayerApplication
}

func readConfigurationLayer() configurationLayer {
	configurationLayerFile, err := ioutil.ReadFile("apps.yml")

	if err != nil {
		log.Printf("Can not unmarshal configuration layer YAML file #%v ", err)
	}

	var configurationLayer configurationLayer
	err = yaml.Unmarshal(configurationLayerFile, &configurationLayer)

	if err != nil {
		log.Printf("Can not unmarshal docker compose YAML file #%v ", err)
	}

	return configurationLayer
}

func writeNewDockerCompose(dockerCompose dockerCompose) {
	out, err := yaml.Marshal(dockerCompose)

	if err != nil {
		log.Printf("Can not marshal docker compose configuration #%v ", err)
	}

	ioutil.WriteFile("./docker-compose-out.yml", out, 500)
}

func captureInput(question string) string {
	reader := bufio.NewReader(os.Stdin)
	fmt.Print(question)
	text, _ := reader.ReadString('\n')
	return text
}

func askQuestion(question string, options []string) string {
	var buffer bytes.Buffer
	defaultOption := options[0]
	buffer.WriteString(question)
	buffer.WriteString(" (")
	buffer.WriteString(strings.Join(options, ","))
	buffer.WriteString(") ")
	buffer.WriteString("[")
	buffer.WriteString(defaultOption)
	buffer.WriteString("] ? ")

	fmt.Println("")
	response := captureInput(buffer.String())

	if len(response) == 0 {
		return defaultOption
	}

	return response
}

func main() {
	configuration := readConfigurationLayer()
	dockerCompose := readDockerCompose()

	for applicationName, configurationApplication := range configuration.Applications {
		question := "Should " + applicationName + " point to image or context"
		answer := askQuestion(question, []string{"I", "C"})

		// In go we can't assign directly to a struct within a map as they are not directly
		// addressable. We work around with a tempory assignment
		currentDockerComposeApplication := dockerCompose.Services[applicationName]
		if configurationType(answer) == Context {
			currentDockerComposeApplication.Context = configurationApplication.Context
			fmt.Println("Using Context in " + applicationName)
		} else {
			currentDockerComposeApplication.Image = configurationApplication.Image
			fmt.Println("Using Image in " + applicationName)
		}
		dockerCompose.Services[applicationName] = currentDockerComposeApplication
	}

	writeNewDockerCompose(dockerCompose)
}
