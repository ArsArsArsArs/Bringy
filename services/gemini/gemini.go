package gemini

import (
	"context"
	"sync"

	"google.golang.org/genai"
)

type ClientManagerType struct {
	clients map[string]*genai.Client
	mutex   sync.RWMutex
}

func (cm *ClientManagerType) getClient(apiKey string) (*genai.Client, error) {
	cm.mutex.RLock()
	client, exists := cm.clients[apiKey]
	cm.mutex.RUnlock()

	if exists {
		return client, nil
	}

	cm.mutex.Lock()
	defer cm.mutex.Unlock()

	//Double-check to handle race conditions
	if client, exists = cm.clients[apiKey]; exists {
		return client, nil
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, &genai.ClientConfig{
		APIKey: apiKey,
	})
	if err != nil {
		return nil, err
	}

	_, err = client.Models.GenerateContent(ctx, "gemini-2.0-flash", genai.Text("test"), nil)
	if err != nil {
		return nil, err
	}

	cm.clients[apiKey] = client
	return client, nil
}

var ClientManager *ClientManagerType

func NewClientManager() {
	ClientManager = &ClientManagerType{
		clients: make(map[string]*genai.Client),
	}
}

func (cm *ClientManagerType) CheckAvailability(apiKey string) bool {
	if _, err := cm.getClient(apiKey); err != nil {
		return false
	}
	return true
}
