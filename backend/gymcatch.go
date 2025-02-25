package backend

import "sync"

type Gymcatch struct {
	mu            sync.RWMutex
	publicSession string
	csrfToken     string
}

func NewGymcatch() *Gymcatch {
	return &Gymcatch{
		publicSession: "",
		csrfToken:     "",
	}
}

func (g *Gymcatch) LoadCreds() error {
	if g.publicSession != "" && g.csrfToken != "" {
		return nil
	}

	g.mu.Lock()
	defer g.mu.Unlock()

	sId, csrfToken, err := GetCreds()
	if err != nil {
		return err
	}

	g.publicSession = sId
	g.csrfToken = csrfToken

	return nil
}

func (g *Gymcatch) RemoveCreds() {
	g.mu.Lock()
	defer g.mu.Unlock()

	g.publicSession = ""
	g.csrfToken = ""
}
