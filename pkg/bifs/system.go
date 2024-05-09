package bifs

import (
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/johnkerl/miller/pkg/mlrval"
	"github.com/johnkerl/miller/pkg/platform"
	"github.com/johnkerl/miller/pkg/version"
)

func BIF_version() *mlrval.Mlrval {
	return mlrval.FromString(version.STRING)
}

func BIF_os() *mlrval.Mlrval {
	return mlrval.FromString(runtime.GOOS)
}

func BIF_hostname() *mlrval.Mlrval {
	hostname, err := os.Hostname()
	if err != nil {
		return mlrval.FromErrorString("could not retrieve system hostname")
	} else {
		return mlrval.FromString(hostname)
	}
}

func BIF_system(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsStringOrVoid() {
		return mlrval.FromNotStringError("system", input1)
	}
	commandString := input1.AcquireStringValue()

	shellRunArray := platform.GetShellRunArray(commandString)

	outputBytes, err := exec.Command(shellRunArray[0], shellRunArray[1:]...).Output()
	if err != nil {
		return mlrval.FromError(err)
	}
	outputString := strings.TrimRight(string(outputBytes), "\n")

	return mlrval.FromString(outputString)
}

func BIF_exec(mlrvals []*mlrval.Mlrval) *mlrval.Mlrval {

	if len(mlrvals) == 0 {
		return mlrval.FromErrorString("exec: zero-length input given")
	}

	cmd := exec.Command(mlrvals[0].String())
	combinedOutput := false

	args := []string{mlrvals[0].String()}
	if len(mlrvals) > 1 {
		for _, val := range mlrvals[1].GetArray()[0:] {
			args = append(args, val.String())
		}
	}
	cmd.Args = args

	if len(mlrvals) > 2 {

		for pe := mlrvals[2].AcquireMapValue().Head; pe != nil; pe = pe.Next {
			if pe.Key == "env" {
				env := []string{}
				for _, val := range pe.Value.GetArray()[0:] {
					env = append(env, val.String())
				}
				cmd.Env = env
			}
			if pe.Key == "dir" {
				cmd.Dir = pe.Value.String()
			}
			if pe.Key == "combined_output" {
				combinedOutput = pe.Value.AcquireBoolValue()
			}

			if pe.Key == "stdin_string" {
				cmd.Stdin = strings.NewReader(pe.Value.String())
			}
		}

	}

	outputBytes := []byte(nil)
	err := error(nil)

	if combinedOutput {
		outputBytes, err = cmd.CombinedOutput()
	} else {
		outputBytes, err = cmd.Output()
	}

	if err != nil {
		return mlrval.FromError(err)
	}

	outputString := strings.TrimRight(string(outputBytes), "\n")
	return mlrval.FromString(outputString)
}

func BIF_stat(input1 *mlrval.Mlrval) *mlrval.Mlrval {
	if !input1.IsStringOrVoid() {
		return mlrval.FromNotStringError("system", input1)
	}
	path := input1.AcquireStringValue()

	fileInfo, err := os.Stat(path)

	if err != nil {
		return mlrval.FromError(err)
	}

	output := mlrval.NewMlrmap()
	output.PutReference("name", mlrval.FromString(fileInfo.Name()))
	output.PutReference("size", mlrval.FromInt(fileInfo.Size()))
	output.PutReference("mode", mlrval.FromIntShowingOctal(int64(fileInfo.Mode())))
	output.PutReference("modtime", mlrval.FromInt(fileInfo.ModTime().UTC().Unix()))
	output.PutReference("isdir", mlrval.FromBool(fileInfo.IsDir()))

	return mlrval.FromMap(output)
}
