package action

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"webhook/src/hook"
	"webhook/src/model"
)

const (
	BashPath string = "#!/bin/bash\n"
)

type Shell struct {
	Hook         *hook.Hook
	Action       *hook.ShellAction
	LogModel     *model.ActionLogClient
	WebhookModel *model.WebhookClient
}

func NewShellAction(action *hook.ShellAction, hook *hook.Hook, webhook *model.WebhookClient, log *model.ActionLogClient) *Shell {
	return &Shell{
		Action:       action,
		Hook:         hook,
		LogModel:     log,
		WebhookModel: webhook,
	}
}

func (s *Shell) writeScriptFile(path string, content string) bool {
	_content := content
	if !strings.HasPrefix(content, "#!") {
		_content = BashPath + content
	}
	err := ioutil.WriteFile(path, []byte(_content), 0755)
	if err != nil {
		s.LogModel.AddErrorLog(fmt.Sprintf("Could not write script file. %s", err.Error()))
		return false
	}
	s.LogModel.AddLog(fmt.Sprintf("Script file created. [%s]", path))
	return true
}

func (s *Shell) removeScriptFile(path string) bool {
	err := os.Remove(path)
	if err != nil {
		s.LogModel.AddWarnLog(fmt.Sprintf("[%s] error removing file %s [%s]", s.Hook.ID, path, err.Error()))
		return false
	}
	s.LogModel.AddLog(fmt.Sprintf("Script file removed. [%s]", path))
	return true
}

func (s *Shell) tryAddChmodX(path string) {

}

func (s *Shell) Exec(envs []string) {
	start := time.Now().UnixMilli()

	if _, err := os.Stat(s.Action.WorkingDirectory); err != nil {
		if os.IsNotExist(err) {
			s.LogModel.AddErrorLog(fmt.Sprintf("Working directory not exists. %s", s.Action.WorkingDirectory))
			return
		}
	}
	lookpath := filepath.Join(s.Action.WorkingDirectory, fmt.Sprintf(".%s.sh", s.Hook.ID))
	ok := s.writeScriptFile(lookpath, s.Action.Scripts)
	if !ok {
		return
	}
	s.tryAddChmodX(lookpath)

	cmdPath, err := exec.LookPath(lookpath)
	if err != nil {
		s.LogModel.AddErrorLog(fmt.Sprintf("Could not find cmd path. %s", err.Error()))
		return
	}
	cmd := exec.Command(cmdPath)
	cmd.Dir = s.Action.WorkingDirectory
	cmd.Env = append(os.Environ(), envs...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		s.LogModel.AddErrorLog(fmt.Sprintf("[%s] error occurred %s", s.Hook.Name, err.Error()))
		return
	}
	s.removeScriptFile(lookpath)

	end := time.Now().UnixMilli()
	s.LogModel.AddLog("Exec successfully. took: " + strconv.FormatInt(end-start, 10) + "ms")
	s.LogModel.AddDebugLog("Output: " + string(out))
}
