package malloc

//#cgo LDFLAGS:
//#include <stdio.h>
//#include <stdlib.h>
//#include <string.h>
import "C"
import "unsafe"

func allocateMemory(size uint64) uintptr {
	//aa := C.CString("aa")
	return uintptr(C.malloc(C.ulong(size)))
}

func freeMemory(p uintptr) {
	C.free(unsafe.Pointer(p))
}

//Объект управления выделяемой памяти
type memSpace struct {
	pointer uintptr //
	count   uint64  //
}

//Создание нового объекта управления памятью
func newMemSpace(size uint64) *memSpace {
	ms := new(memSpace)
	ms.pointer = allocateMemory(size)
	return ms
}

//Очистка выделенной памяти
func (ms *memSpace) free() {
	freeMemory(ms.pointer)
}
