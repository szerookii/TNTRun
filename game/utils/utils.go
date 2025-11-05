package utils

import "github.com/df-mc/dragonfly/server/world"

// RemoveIndex removes T from s at the specified index.
func RemoveIndex[T any](s []T, index int) []T {
	return append(s[:index], s[index+1:]...)
}

// EntityHandleToEntity safely converts a world.EntityHandle to its underlying entity of type T.
func EntityHandleToEntity[T any](handle *world.EntityHandle) *T {
	if handle == nil {
		return nil
	}
	var entity *T
	handle.ExecWorld(func(tx *world.Tx, e world.Entity) {
		if casted, ok := e.(T); ok {
			entity = &casted
		}
	})
	return entity
}
