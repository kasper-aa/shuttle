package config

import (
	"io/ioutil"
	"log"
	"path"
	"path/filepath"
	"regexp"

	"fmt"
	"os"

	"bitbucket.org/LunarWay/shuttle/pkg/git"
	"gopkg.in/yaml.v2"
)

// ShuttlePlanScript is a ShuttlePlan sub-element
type ShuttlePlanScript struct {
	Description string              `yaml:"description"`
	Actions     []ShuttleAction     `yaml:"actions"`
	Args        []ShuttleScriptArgs `yaml:"args"`
}

type ShuttleScriptArgs struct {
	Name     string `yaml:"name"`
	Required bool   `yaml:"required"`
}

type ShuttleAction struct {
	Shell      string `yaml:"shell"`
	Dockerfile string `yaml:"dockerfile"`
}

// ShuttlePlanConfiguration is a ShuttlePlan sub-element
type ShuttlePlanConfiguration struct {
	Scripts map[string]ShuttlePlanScript `yaml:"scripts"`
}

// ShuttlePlan struct describes a plan
type ShuttlePlan struct {
	ProjectPath   string
	LocalPlanPath string
	Configuration ShuttlePlanConfiguration
}

// Load loads a plan from project path and shuttle config
func (p *ShuttlePlanConfiguration) Load(planPath string) *ShuttlePlanConfiguration {
	var configPath = path.Join(planPath, "plan.yaml")
	//log.Printf("configpath: %s", configPath)
	yamlFile, err := ioutil.ReadFile(configPath)
	if err != nil {
		log.Printf("yamlFile.Get err   #%v ", err)
	}
	err = yaml.Unmarshal(yamlFile, p)
	if err != nil {
		log.Fatalf("Unmarshal: %v", err)
	}

	return p
}

// FetchPlan so it exists locally and return path to that plan
func FetchPlan(plan string, projectPath string, localShuttleDirectoryPath string) string {
	switch {
	case git.IsGitPlan(plan):
		return git.GetGitPlan(plan, localShuttleDirectoryPath)
	case isMatching("^http://|^https://", plan):
		panic("plan not valid: http is not supported yet")
	case isFilePath(plan, true):
		return plan
	case isFilePath(plan, false):
		return path.Join(projectPath, plan)

	}
	panic("Unknown plan path '" + plan + "'")
}

func isFilePath(path string, matchOnlyAbs bool) bool {
	return filepath.IsAbs(path) == matchOnlyAbs
}

func isMatching(r string, content string) bool {
	match, err := regexp.MatchString(r, content)
	if err != nil {
		panic(err)
	}
	return match
}

// CheckIfError should be used to naively panics if an error is not nil.
func CheckIfError(err error) {
	if err == nil {
		return
	}

	fmt.Printf("\x1b[31;1m%s\x1b[0m\n", fmt.Sprintf("error: %s", err))
	os.Exit(1)
}