package config

import (
	"github.com/xdg-go/scram"
)

// XDGSCRAMClient implements sarama.SCRAMClient interface
type XDGSCRAMClient struct {
	*scram.Client
	*scram.ClientConversation
	scram.HashGeneratorFcn
}

func (x *XDGSCRAMClient) Begin(userName, password, authzID string) error {
	client, err := x.HashGeneratorFcn.NewClient(userName, password, authzID)
	if err != nil {
		return err
	}
	x.Client = client

	conversation := x.Client.NewConversation()
	x.ClientConversation = conversation

	return nil
}

func (x *XDGSCRAMClient) Step(challenge string) (string, error) {
	return x.ClientConversation.Step(challenge)
}

func (x *XDGSCRAMClient) Done() bool {
	return x.ClientConversation.Done()
}

// Hash generators for different mechanisms
var (
	SHA256 scram.HashGeneratorFcn = scram.SHA256
	SHA512 scram.HashGeneratorFcn = scram.SHA512
)
