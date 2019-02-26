package grmln

import "fmt"

type clusterError int

const (
	clusterErrorClusterClosed clusterError = iota
)

var clusterErrorStrings = map[clusterError]string{
	clusterErrorClusterClosed: "Cluster Closed",
}

func (e clusterError) Error() string {
	str, ok := clusterErrorStrings[e]
	if !ok {
		return fmt.Sprintf("invalid cluster error: %d", e)
	}

	return str
}

func (e clusterError) IsClusterClosed() bool {
	return e == clusterErrorClusterClosed
}

type clusterClosed interface {
	IsClusterClosed() bool
}

// IsClusterClosed returns whether or not the error is a cluster closed error
func IsClusterClosed(err error) bool {
	e, ok := err.(clusterClosed)
	return ok && e.IsClusterClosed()
}
