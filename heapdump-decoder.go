//https://github.com/golang/go/wiki/heapdump15-through-heapdump17
package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
	"os"
)

const (
	DT_EOF = iota
	DT_Object
	DT_OtherRoot
	DT_Type
	DT_Goroutine
	DT_StackFrame
	DT_DumpParams
	DT_RegisteredFinalizer
	DT_Itab
	DT_OSThread
	DT_MemStats
	DT_QueuedFinalizer
	DT_DataSegment
	DT_BSSSegment
	DT_DeferRecord
	DT_PanicRecord
	DT_AllocFreeProfileRecord
	DT_AllocStackTraceSample
)

var types = map[uint64]string{
	0:  "EOF",
	1:  "object",
	2:  "otherroot",
	3:  "type",
	4:  "goroutine",
	5:  "stack frame",
	6:  "dump params",
	7:  "registered finalizer",
	8:  "itab",
	9:  "OS thread",
	10: "mem stats",
	11: "queued finalizer",
	12: "data segment",
	13: "bss segment",
	14: "defer record",
	15: "panic record",
	16: "alloc/free profile record",
	17: "alloc stack trace sample",
}

var statusmap = map[uint64]string{
	0: "idle",
	1: "runnable",
	3: "syscall",
	4: "waiting",
}

func readString(r *bufio.Reader) (string, error) {
	l, err := binary.ReadUvarint(r)
	if err != nil {
		return "", err
	}

	b := make([]byte, l)
	_, err = io.ReadFull(r, b)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func readPrintStr(r *bufio.Reader, title string) {
	str, err := readString(r)
	if err != nil {
		panic(err)
	}
	fmt.Println(title, ":", str)
}

func readPrintStrLen(r *bufio.Reader, title string) {
	str, err := readString(r)
	if err != nil {
		panic(err)
	}
	fmt.Println(title, ":", len(str))
}

func readPrintAddr(r *bufio.Reader, title string) {
	addr, err := binary.ReadUvarint(r)
	if err != nil {
		panic(err)
	}

	hex_addr := fmt.Sprintf("0x%x", addr)
	fmt.Println(title, ":", hex_addr)
}

func readPrintInt(r *bufio.Reader, title string) uint64 {
	u64, err := binary.ReadUvarint(r)
	if err != nil {
		panic(err)
	}

	fmt.Println(title, ":", u64)
	return u64
}

func readFieldList(r *bufio.Reader) {
	for {
		field_kind, err := binary.ReadUvarint(r)
		if err != nil {
			break
		}
		if field_kind == 0 {
			break
		}

		field_offset, err := binary.ReadUvarint(r)
		if err != nil {
			break
		}

		fmt.Println("field_kind: ", field_kind, "field_offset: ", field_offset)
	}
}

func main() {
	file, err := os.Open("metadump")
	if err != nil {
		panic(err)
	}

	defer file.Close()

	r := bufio.NewReader(file)
	ver, _, err := r.ReadLine()
	fmt.Println(string(ver))

	for {
		i, err := binary.ReadUvarint(r)
		if err != nil {
			fmt.Println(err)
			break
		}

		t := types[i]
		if t == "" {
			fmt.Println("Unknown type: ", i)
			break
		}

		fmt.Println("\nType: ", t)

		switch i {
		case DT_EOF:
			break

		case DT_Object:
			readPrintAddr(r, "addr")
			readPrintStrLen(r, "content of object(len)")
			readFieldList(r)

		case DT_Type:
			readPrintAddr(r, "addr")
			readPrintInt(r, "size")
			readPrintStr(r, "name")
			readPrintInt(r, "type2ptr")

		case DT_Goroutine:
			readPrintAddr(r, "addr")
			readPrintAddr(r, "stacktop")
			readPrintInt(r, "goid")
			readPrintAddr(r, "create location(rip)")
			status := readPrintInt(r, "status")
			fmt.Println("status: ", statusmap[status])
			readPrintInt(r, "create by sys")
			readPrintInt(r, "background goroutine")
			readPrintInt(r, "last start waiting(ns)")
			readPrintStr(r, "wait reason")
			readPrintAddr(r, "context pointer of currently running frame")
			readPrintAddr(r, "M")
			readPrintInt(r, "top defer record")
			readPrintInt(r, "top panic record")

		case DT_StackFrame:
			readPrintAddr(r, "stack pointer")
			readPrintInt(r, "depth")
			readPrintAddr(r, "child stack pointer")
			readPrintStrLen(r, "content of stack frame(len)")
			readPrintAddr(r, "entry pc")
			readPrintAddr(r, "current pc")
			readPrintAddr(r, "continuation pc")
			readPrintStr(r, "function name")
			readFieldList(r)

		case DT_DumpParams:
			readPrintInt(r, "bigendian")
			readPrintInt(r, "pointer sz")
			readPrintAddr(r, "start addr")
			readPrintAddr(r, "end addr")
			readPrintStr(r, "param")
			readPrintStr(r, "env")
			readPrintInt(r, "CPU")

		case DT_RegisteredFinalizer:
			readPrintAddr(r, "object addr")
			readPrintAddr(r, "funval of the finalizer")
			readPrintAddr(r, "pc of finalizer")
			readPrintInt(r, "type of finalizer argument")
			readPrintInt(r, "type of object")

		case DT_Itab:
			readPrintAddr(r, "addr")
			readPrintAddr(r, "contained")

		case DT_OSThread:
			readPrintAddr(r, "os thread descriptor")
			readPrintInt(r, "go internal id of thread")
			readPrintInt(r, "os's id for thread")

		case DT_MemStats:
			readPrintInt(r, "Alloc")
			readPrintInt(r, "TotalAlloc")
			readPrintInt(r, "Sys")
			readPrintInt(r, "Lookups")
			readPrintInt(r, "Mallocs")
			readPrintInt(r, "Frees")
			readPrintInt(r, "HeapAlloc")
			readPrintInt(r, "HeapSys")
			readPrintInt(r, "HeapIdle")
			readPrintInt(r, "HeapInuse")
			readPrintInt(r, "HeapReleased")
			readPrintInt(r, "HeapObjects")
			readPrintInt(r, "StackInuse")
			readPrintInt(r, "StackSys")
			readPrintInt(r, "MSpanInuse")
			readPrintInt(r, "MSpanSys")
			readPrintInt(r, "MCacheInuse")
			readPrintInt(r, "MCacheSys")
			readPrintInt(r, "BuckHashSys")
			readPrintInt(r, "GCSys")
			readPrintInt(r, "OtherSys")
			readPrintInt(r, "NextGC")
			readPrintInt(r, "LastGC")
			readPrintInt(r, "PauseTotalNs")
			for i := 0; i < 256; i++ {
				readPrintInt(r, "PauseNs")
			}
			readPrintInt(r, "NumGC")

		case DT_DataSegment:
			readPrintAddr(r, "start addr")
			readPrintStrLen(r, "contents of the data segment")
			readFieldList(r)

		case DT_BSSSegment:
			readPrintAddr(r, "start addr")
			readPrintStrLen(r, "contents of the data segment")
			readFieldList(r)

		case DT_AllocFreeProfileRecord:
			readPrintInt(r, "record idnetifier")
			readPrintInt(r, "size of object")
			frames := readPrintInt(r, "Number of stack frames")
			for i := 0; i < int(frames); i++ {
				readPrintStr(r, "function name")
				readPrintStr(r, "file name")
				readPrintInt(r, "line number")
			}
			readPrintInt(r, "number of allocations")
			readPrintInt(r, "number of frees")

		case DT_AllocStackTraceSample:
			readPrintAddr(r, "object addr")
			readPrintAddr(r, "alloc/free profile record identifier")

		default:
			panic("Unknown type")
		}
	}
}
