package repository

import (
	"context"
	"fmt"
	"sync"

	"github.com/mmfshirokan/GoProject1/internal/model"
)

type MapRepositoryInterface[value *model.User | []*model.RefreshToken] interface {
	Set(ctx context.Context, key string, val value)
	Get(ctx context.Context, key string) (value, error)
	Remove(key string)
}

type MapRepository[value *model.User | []*model.RefreshToken] struct {
	source map[string]value
	mut    sync.RWMutex
}

func NewUserMap(source map[string]*model.User) MapRepositoryInterface[*model.User] { // changed: before map was created in func; now map is created outside the func
	return &MapRepository[*model.User]{
		source: source,
	}
}

func NewRftMap(source map[string][]*model.RefreshToken) MapRepositoryInterface[[]*model.RefreshToken] { // changed: before map was created in func; now map is created outside the func
	return &MapRepository[[]*model.RefreshToken]{
		source: source,
	}
}

func (mapRep *MapRepository[value]) Set(ctx context.Context, key string, val value) {
	mapRep.mut.Lock()

	mapRep.source[key] = val

	mapRep.mut.Unlock()
}

func (mapRep *MapRepository[value]) Get(ctx context.Context, key string) (value, error) {
	mapRep.mut.RLock()

	res, ok := mapRep.source[key]

	mapRep.mut.RUnlock()

	if !ok {
		return nil, fmt.Errorf("there is no object with %s key", key)
	}

	return res, nil
}

func (mapRep *MapRepository[value]) Remove(key string) {
	mapRep.mut.Lock()

	delete(mapRep.source, key)

	mapRep.mut.Unlock()
}
