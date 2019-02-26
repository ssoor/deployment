package main

import (
	"io/ioutil"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	yaml "gopkg.in/yaml.v2"
)

func main() {
	New(false).Execute()
}

func makeTarget(name string, targets map[string]FileTarget) FileTarget {
	t := targets[name]

	for _, name := range t.Imports {
		targets[name] = makeTarget(name, targets)

		mergeStruct(&t, targets[name])
	}

	return t
}

// CommandRun is
func CommandRun(cmd *cobra.Command, conf *viper.Viper, args []string) {
	name := args[1]
	targetName := args[0]
	deploy := DeployFile{}

	context, err := ioutil.ReadFile(cmd.Flag("config").Value.String())
	if nil != err {
		log.Fatal(err)
	}

	if err := yaml.Unmarshal(context, &deploy); nil != err {
		log.Fatal(err)
	}

	envMap := map[string]string{"PORT": "95270", "GAMEID": "1234"}

	target := buildTarget(targetName, envMap, deploy.Targets)

	if err := buildAssets(target.Assets); nil != err {
		log.Fatal(err)
	}

	for _, runtime := range target.Runtime {
		switch runtime.Type {
		case "systemd":
			if def, err := buildSystemd(name, target.Command, runtime); nil != err {
				log.Fatal(err)
			} else {
				log.Info(def)
			}
		}
	}

	context, err = yaml.Marshal(target)

	log.Info(string(context), err)
}
