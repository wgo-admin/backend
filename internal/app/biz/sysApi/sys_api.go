package sysApi

import "github.com/wgo-admin/backend/internal/app/store"

type ISysApiBiz interface {
}

type sysApiBiz struct {
	ds store.IStore
}

var _ ISysApiBiz = (*sysApiBiz)(nil)

func NewBiz(ds store.IStore) *sysApiBiz {
	return &sysApiBiz{ds}
}
