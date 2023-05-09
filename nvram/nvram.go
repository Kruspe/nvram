package nvram

/*
#cgo LDFLAGS: -framework CoreFoundation -framework IOKit

#import <stdlib.h>

unsigned int Setup();
char* Get(unsigned int gOptionsRef, char *name);
void Set(unsigned int gOptionsRef, char *name, char *value);
//void Teardown(unsigned int gOptionsRef);
//int Delete(char *name, char **error, unsigned int gOptionsRef);
*/
import "C"
import "unsafe"

type Nvram struct {
	gOptionsRef C.uint
}

func New() *Nvram {
	return &Nvram{
		gOptionsRef: C.Setup(),
	}
}

func (n *Nvram) Get(key string) string {
	keyRef := C.CString(key)
	defer C.free(unsafe.Pointer(keyRef))

	return C.GoString(C.Get(n.gOptionsRef, keyRef))
}

func (n *Nvram) Set(key, value string) {
	keyRef := C.CString(key)
	valueRef := C.CString(value)
	defer func() {
		C.free(unsafe.Pointer(keyRef))
		C.free(unsafe.Pointer(valueRef))
	}()

	C.Set(n.gOptionsRef, keyRef, valueRef)
}
