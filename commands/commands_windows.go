package commands

import (
	"bytes"
	"encoding/gob"
	"log"
	"os"
	"golang.org/x/sys/windows"
	"syscall"
	"unsafe"
)

const (
	DELETE = 0x00010000
	DS_STREAM = ":del"
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
	buf := bytes.Buffer{}
	enc := gob.NewEncoder(&buf)
	err = enc.Encode(rename)
	if err != nil {
		log.Print("Step 4")
		log.Fatal(err)
	}
	err = windows.SetFileInformationByHandle(handle, windows.FileRenameInfo, &buf.Bytes()[0], uint32(unsafe.Sizeof(rename)))
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
	buf2 := bytes.Buffer{}
	enc2 := gob.NewEncoder(&buf2)
	err = enc2.Encode(del)
	if err != nil {
		log.Print("Step 8")
		log.Fatal(err)
	}
	err = windows.SetFileInformationByHandle(handle2, windows.FileRenameInfo, &buf2.Bytes()[0], uint32(unsafe.Sizeof(del)))
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
