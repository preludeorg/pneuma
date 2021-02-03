package commands

import (
	"golang.org/x/sys/windows"
	"log"
	"os"
	"syscall"
	"unsafe"
)

const (
	DELETE = 0x00010000
	DS_STREAM = ":del"
	errnoERROR_IO_PENDING = 997
)

var (
	modkernel32 = syscall.NewLazyDLL("kernel32.dll")

	procSetFileInformationByHandle   = modkernel32.NewProc("SetFileInformationByHandle")
	errERROR_IO_PENDING error = syscall.Errno(errnoERROR_IO_PENDING)
	errERROR_EINVAL     error = syscall.EINVAL
)

type FILE_RENAME_INFO struct {
	DummyUnionName  uint32
	ReplaceIfExists int8
	RootDirectory   windows.Handle
	FileNameLength	uint32
	FileName		uint16
}

type FILE_DISPOSITION_INFO struct {
	DeleteFile int8
}

func Shutdown() {
	path, err := os.Executable()
	if err != nil {
		log.Print("Step 1")
		log.Fatal(err)
	}
	u16pathname, err := syscall.UTF16FromString(path)
	if err != nil {
		log.Print("Step 2")
		log.Fatal(err)
	}
	handle, err := windows.CreateFile(&u16pathname[0], syscall.GENERIC_READ | syscall.SYNCHRONIZE | DELETE , 0, nil, syscall.OPEN_EXISTING, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		log.Print("Step 3")
		log.Fatal(err)
	}
	u16rename, err := syscall.UTF16FromString(DS_STREAM)
	rename := &FILE_RENAME_INFO{
		FileNameLength: uint32(len(DS_STREAM)),
		FileName: u16rename[0],
	}
	err = SetFileInformationByHandle(handle, windows.FileRenameInfo, uintptr(unsafe.Pointer(&rename)), uint32(unsafe.Sizeof(rename)) + uint32(unsafe.Sizeof(DS_STREAM)))
	if err != nil {
		log.Print("Step 5")
		log.Fatal(err)
	}

	err = windows.CloseHandle(handle)
	if err != nil {
		log.Print("Step 6")
		log.Fatal(err)
	}

	handle2, err := windows.CreateFile(&u16pathname[0], syscall.GENERIC_READ | syscall.SYNCHRONIZE | DELETE , 0, nil, syscall.OPEN_EXISTING, syscall.FILE_ATTRIBUTE_NORMAL, 0)
	if err != nil {
		log.Print("Step 7")
		log.Fatal(err)
	}

	del := &FILE_DISPOSITION_INFO{
		DeleteFile: 0x00000001,
	}

	err = SetFileInformationByHandle(handle2, windows.FileRenameInfo, uintptr(unsafe.Pointer(&del)), uint32(unsafe.Sizeof(del)))
	if err != nil {
		log.Print("Step 9")
		log.Fatal(err)
	}

	err = windows.CloseHandle(handle2)
	if err != nil {
		log.Print("Step 10")
		log.Fatal(err)
	}

}

func SetFileInformationByHandle(handle windows.Handle, fileInformationClass uint32, buf uintptr, bufsize uint32) (err error) {
	r1, _, e1 := syscall.Syscall6(procSetFileInformationByHandle.Addr(), 4, uintptr(handle), uintptr(fileInformationClass), uintptr(buf), uintptr(bufsize), 0, 0)
	if r1 == 0 {
		err = errnoErr(e1)
	}
	return
}

func errnoErr(e syscall.Errno) error {
	switch e {
	case 0:
		return errERROR_EINVAL
	case errnoERROR_IO_PENDING:
		return errERROR_IO_PENDING
	}
	// TODO: add more here, after collecting data on the common
	// error values see on Windows. (perhaps when running
	// all.bat?)
	return e
}