package utils

import (
	"NotificationManagement/utils/errutil"
	"errors"
	"reflect"

	"gorm.io/gorm"
)

func SyncHasManyAssociation(db *gorm.DB, parent interface{}, assocName string, newItems interface{}) (interface{}, error) {
	// Use reflection to create a slice of the same type as newItems
	newVal := reflect.ValueOf(newItems)
	if newVal.Kind() != reflect.Ptr || newVal.Elem().Kind() != reflect.Slice {
		return nil, errutil.NewAppError(errutil.ErrGormInvalidSlicePointer, errors.New("newItems must be a pointer to a slice"))
	}

	sliceType := newVal.Elem().Type()
	existingItems := reflect.New(sliceType).Interface()

	assoc := db.Model(parent).Association(assocName)
	if assoc == nil {
		return nil, errutil.NewAppError(errutil.ErrGormAssociationNotFound, errors.New("association not found"))
	}
	if err := assoc.Find(existingItems); err != nil {
		return nil, errutil.NewAppError(errutil.ErrGormFindFailed, err)
	}

	// Build map of existing by ID
	existingMap := make(map[uint]reflect.Value)
	existingSlice := reflect.ValueOf(existingItems).Elem()
	for i := 0; i < existingSlice.Len(); i++ {
		item := existingSlice.Index(i)
		idField := item.FieldByName("ID")
		if !idField.IsValid() || idField.IsZero() {
			continue
		}
		id := uint(idField.Uint())
		existingMap[id] = item
	}

	// Prepare for update/create/delete
	newSlice := newVal.Elem()
	for i := 0; i < newSlice.Len(); i++ {
		item := newSlice.Index(i)
		idField := item.FieldByName("ID")
		var id uint
		if idField.IsValid() && !idField.IsZero() {
			id = uint(idField.Uint())
		}

		if id == 0 {
			// New item: create
			if err := db.Create(item.Addr().Interface()).Error; err != nil {
				return nil, errutil.NewAppError(errutil.ErrGormCreateFailed, err)
			}
		} else if existing, ok := existingMap[id]; ok {
			// Existing item: update
			if err := db.Model(existing.Addr().Interface()).Updates(item.Addr().Interface()).Error; err != nil {
				return nil, errutil.NewAppError(errutil.ErrGormUpdateFailed, err)
			}
			delete(existingMap, id)
		}
	}

	// Delete items not in new slice
	for _, item := range existingMap {
		if err := db.Delete(item.Addr().Interface()).Error; err != nil {
			return nil, errutil.NewAppError(errutil.ErrGormDeleteFailed, err)
		}
	}

	// Fetch the final association result
	finalItems := reflect.New(sliceType).Interface()
	if err := db.Model(parent).Association(assocName).Find(finalItems); err != nil {
		return nil, errutil.NewAppError(errutil.ErrGormFindFailed, err)
	}

	return finalItems, nil
}
