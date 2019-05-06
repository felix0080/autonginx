package model


type Context struct {
	Stop chan struct{}
	Done chan struct{}
}

