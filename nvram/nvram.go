package nvram

/*
#cgo LDFLAGS: -framework CoreFoundation -framework IOKit

#import <stdlib.h>

unsigned int Setup();
char* Get(unsigned int gOptionsRef, char *name, char **err);
void Set(unsigned int gOptionsRef, char *name, char *value, char **err);
void Teardown(unsigned int gOptionsRef);
void Delete(unsigned int gOptionsRef, char *name, char **error);
*/
import "C"
import (
	"errors"
	"unsafe"
)

type Nvram struct {
	gOptionsRef C.uint
}

func New() *Nvram {
	return &Nvram{
		gOptionsRef: C.Setup(),
	}
}

func (n *Nvram) Teardown() {
	C.Teardown(n.gOptionsRef)
}

func (n *Nvram) Get(key string) (string, error) {
	keyRef := C.CString(key)
	defer C.free(unsafe.Pointer(keyRef))

	var errRef *C.char
	result := C.Get(n.gOptionsRef, keyRef, &errRef)
	if errRef != nil {
		err := errors.New(C.GoString(errRef))
		C.free(unsafe.Pointer(errRef))
		return "", err
	}
	return C.GoString(result), nil
}

func (n *Nvram) Set(key, value string) error {
	keyRef := C.CString(key)
	valueRef := C.CString(value)
	defer func() {
		C.free(unsafe.Pointer(keyRef))
		C.free(unsafe.Pointer(valueRef))
	}()

	var errRef *C.char
	C.Set(n.gOptionsRef, keyRef, valueRef, &errRef)
	if errRef != nil {
		err := errors.New(C.GoString(errRef))
		C.free(unsafe.Pointer(errRef))
		return err
	}
	return nil
}

func (n *Nvram) Delete(key string) error {
	keyRef := C.CString(key)
	defer C.free(unsafe.Pointer(keyRef))

	var errRef *C.char
	C.Delete(n.gOptionsRef, keyRef, &errRef)
	if errRef != nil {
		err := errors.New(C.GoString(errRef))
		C.free(unsafe.Pointer(errRef))
		return err
	}
	return nil
}
