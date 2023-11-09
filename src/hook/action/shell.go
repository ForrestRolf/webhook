package action

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"webhook/src/hook"
	"webhook/src/model"
)

const (
	BashPath string = "#!/bin/bash\n"
)

type Shell struct {
	Hook         *hook.Hook
	Action       *hook.ShellAction
	LogModel     *model.LogClient
	WebhookModel *model.WebhookClient
}

func NewShellAction(action *hook.ShellAction, hook *hook.Hook, log *model.LogClient, webhook *model.WebhookClient) *Shell {
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
		s.LogModel.AddErrorLog(s.Hook, fmt.Sprintf("[Shell] Could not write script file. %s", err))
		return false
	}
	s.LogModel.AddLog(s.Hook, fmt.Sprintf("[Shell] Script file created. [%s]", path))
	return true
}

func (s *Shell) removeScriptFile(path string) bool {
	err := os.Remove(path)
	if err != nil {
		s.LogModel.AddWarnLog(s.Hook, fmt.Sprintf("[Shell] [%s] error removing file %s [%s]", s.Hook.ID, path, err))
		return false
	}
	s.LogModel.AddLog(s.Hook, fmt.Sprintf("[Shell] Script file removed. [%s]", path))
	return true
}

func (s *Shell) tryAddChmodX(path string) {

}

func (s *Shell) Exec(envs []string) {
	if _, err := os.Stat(s.Action.WorkingDirectory); err != nil {
		if os.IsNotExist(err) {
			s.LogModel.AddErrorLog(s.Hook, fmt.Sprintf("[Shell] Working directory not exists. %s", s.Action.WorkingDirectory))
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
		s.LogModel.AddErrorLog(s.Hook, fmt.Sprintf("[Shell] Could not find cmd path. %w", err))
		return
	}
	cmd := exec.Command(cmdPath)
	cmd.Dir = s.Action.WorkingDirectory
	cmd.Env = append(os.Environ(), envs...)

	out, err := cmd.CombinedOutput()
	if err != nil {
		s.LogModel.AddErrorLog(s.Hook, fmt.Sprintf("[Shell] [%s] error occurred %s", s.Hook.Name, err))
		return
	}
	s.LogModel.AddLog(s.Hook, "[Shell] Exec successfully")
	s.LogModel.AddDebugLog(s.Hook, "[Shell] Output: "+string(out))
	s.removeScriptFile(lookpath)
}
