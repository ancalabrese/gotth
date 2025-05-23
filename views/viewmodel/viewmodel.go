package viewmodel

// ViewModelData is a function that modifies the ViewModel struct
type ViewModelData[T any] func(vm *T)

// NewViewModel creates a new ViewModel of type T. It initialize the ViewModel fields with
// the ViewModelData
func NewViewModel[T any](data ...ViewModelData[T]) *T {
	// Create zeroed valued instance of ViewModel type
	vm := new(T)

	for _, d := range data {
		d(vm)
	}

	return vm
}
