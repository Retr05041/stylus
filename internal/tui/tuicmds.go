package tui

import (
    "stylus/internal/api"
	tea "github.com/charmbracelet/bubbletea"
)

func LoginToApi(email string, pwd string) tea.Cmd {
    return func() tea.Msg {
        newSession, err := api.Login(email , pwd)                
        if err != nil {
            return errMsg{err}
        }
        return loginSuccessMsg{successfulSession: newSession}
    }
}
