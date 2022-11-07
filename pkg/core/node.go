package core

import (
	"github.com/machinefi/w3bstream/pkg/core/access"
	"github.com/machinefi/w3bstream/pkg/core/storage"
	"github.com/machinefi/w3bstream/pkg/core/types"
	"github.com/machinefi/w3bstream/pkg/core/web3"
)

type (
	// Processor defines an interface which processes an input message with given authentication
	Processor interface {
		Process(access.Authentication, storage.KVStore, string, uint8, []byte) error
	}
	Config struct {
	}

	// Node defines a w3bstream node
	Node struct {
		ac        access.Control
		client    web3.Client
		rawDB     storage.KVStore
		processor Processor
	}
)

// NewNode creates a new w3bstream node
func NewNode(cfg Config, ac access.Control, client web3.Client, rawDB storage.KVStore) (*Node, error) {
	// TODO: create processor according to cfg, with node
	return &Node{
		ac:        ac,
		client:    client,
		rawDB:     rawDB,
		processor: nil,
	}, nil
}

// Put puts message into db after processing
//
//	Single node mode, without consensus module
func (node *Node) Put(msg types.Message) error {
	auth, err := node.ac.Check(msg.Sender, msg.Nonce, msg.Hash(), msg.Authentication)
	if err != nil {
		return err
	}
	// TODO: replace wrap node.db into a database with access control and some instantiated interface
	db := node.rawDB
	if err := node.processor.Process(auth, db, msg.Sender, msg.Type, msg.Data); err != nil {
		return err
	}

	return node.commit(msg, db)
}

func (node *Node) commit(msg types.Message, db storage.KVStore) error {
	return nil
}

// AccessControl returns the instance of access control
func (node *Node) AccessControl() access.Control {
	return node.ac
}

// Web3Client returns the web3 client
func (node *Node) Web3Client() web3.Client {
	return node.client
}
