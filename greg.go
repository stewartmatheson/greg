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
    "github.com/fatih/color"
)

type dockerComposeApplication struct {
    Image       string   `yaml:"image,omitempty"`
    Tty         bool     `yaml:"tty,omitempty"`
    Build       string   `yaml:"build,omitempty"`
    Ports       []string `yaml:"ports,omitempty"`
    Volumes     []string `yaml:"volumes,omitempty"`
    Environment []string `yaml:"environment,omitempty"`
    DependsOn   []string `yaml:"depends_on,omitempty"`
    Entrypoint  string   `yaml:"entrypoint,omitempty"`
}

type configurationType string

const (
    // Build ... Comment
    Build configurationType = "B"

    // Image ... Comment
    Image configurationType = "I"
)

func (d dockerComposeApplication) configuredFor() (configurationType, error) {
    if len(d.Image) > 0 {
        return Image, nil
    }

    if len(d.Build) > 0 {
        return Build, nil
    }

    return "", errors.New("Invalid Configuration")
}

type dockerCompose struct {
    Version  string `yaml:"version"`
    Services map[string]dockerComposeApplication
}

func readDockerCompose() dockerCompose {
    dockerComposeFile, err := ioutil.ReadFile(".docker-compose.yml")

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
    Build string `yaml:"build,omitempty"`
    Image string `yaml:"image,omitempty"`
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

    ioutil.WriteFile("./docker-compose.yml", out, 500)
}

func captureInput(question string) string {
    reader := bufio.NewReader(os.Stdin)
    fmt.Print(question)
    text, _ := reader.ReadString('\n')
    return strings.TrimSuffix(text, "\n")
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

func resolveConfigurationType(textAnswer string) (configurationType, error) {
    if textAnswer == "B" {
        return Build, nil
    }

    if textAnswer == "I" {
        return Image, nil
    }

    // Can't retrun nil here so I guess we return anything as the first param
    return Build, errors.New("Unknown configuration type")
}

func main() {
    configuration := readConfigurationLayer()
    dockerCompose := readDockerCompose()

    for applicationName, configurationApplication := range configuration.Applications {
        question := "Should " + applicationName + " point to image or context"
        answer := askQuestion(question, []string{"I", "B"})

        // Here I have a whole function to get the configuration type from the user answer
        // There might be a better way to do this. I'm not sure. I have a feeling using
        // strings in an enum is not the most idomatic go
        answeredConfigurationType, err := resolveConfigurationType(answer)

        if err != nil {
            fmt.Println(err)
            // For now let's just give up
            os.Exit(1)
        }

        // In go we can't assign directly to a struct within a map as they are not directly
        // addressable. We work around with a tempory assignment
        currentDockerComposeApplication := dockerCompose.Services[applicationName]
        if answeredConfigurationType == Build {
            currentDockerComposeApplication.Build = configurationApplication.Build
            currentDockerComposeApplication.Image = ""
            fmt.Println("Using Build in " + applicationName)
        } else {
            currentDockerComposeApplication.Image = configurationApplication.Image
            currentDockerComposeApplication.Build = ""
            fmt.Println("Using Image in " + applicationName)
        }
        dockerCompose.Services[applicationName] = currentDockerComposeApplication
    }

    writeNewDockerCompose(dockerCompose)
}
