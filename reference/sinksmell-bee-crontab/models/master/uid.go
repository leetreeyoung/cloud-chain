package master

import (
	"sync"

	"github.com/bwmarrin/snowflake"
)

var (
	once sync.Once
	node *snowflake.Node
)

// NewID generate global unique id
func NewID() string {
	once.Do(func() {
		var err error
		node, err = snowflake.NewNode(1)
		if err != nil {
			panic(err)
		}
	})

	return node.Generate().String()
}
