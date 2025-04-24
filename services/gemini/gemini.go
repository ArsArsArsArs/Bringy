package gemini

import (
	"sync"

	"google.golang.org/genai"
)

type ClientManager struct {
	Clients map[string]*genai.Client
	mutex   sync.RWMutex
}
