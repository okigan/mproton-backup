package mproton

/*
#cgo darwin CFLAGS: -x objective-c
#cgo darwin CFLAGS: -Wobjc-property-no-attribute
#cgo darwin LDFLAGS: -framework Cocoa
#cgo darwin LDFLAGS: -framework Foundation
#cgo darwin LDFLAGS: -framework WebKit

//sudo apt-get -y install libgtk-3-dev libwebkit2gtk-4.0-dev libappindicator3-dev

#cgo linux openbsd freebsd CXXFLAGS: -DWEBVIEW_GTK -std=c++11
#cgo linux openbsd freebsd pkg-config: gtk+-3.0 webkit2gtk-4.0 appindicator3-0.1


// sudo apt-get install -y electric-fence
// #cgo linux LDFLAGS: -lefence

#include "mproton.h"


#include <stdio.h>
#include <stdlib.h>
*/
import "C"
import "fmt"

var callbackMap = make(map[string]func(v string) (string, error))

func registerCallback(name string, callback func(v string) (string, error)) int {
	callbackMap[name] = callback
	return 0
}

//export goTrampoline
func goTrampoline(a int, param1 *C.char, param2 *C.char) (*C.char, *C.char) {
	p1 := C.GoString(param1)
	p2 := C.GoString(param2)

	println("[golang]", a, p1, p2)
	callback := callbackMap[p1]

	if callback == nil {
		msg := fmt.Sprintf("No callback registered for: %s", p1)
		return nil, C.CString(msg)
	}

	r1, err := callback(p2)

	if err != nil {
		return C.CString(r1), C.CString(err.Error())
	} else {
		return C.CString(r1), nil
	}

}

type mProtonApp interface {
	Run()
	// Destroy()
	SetTitle(path string)
	SetContentPath(path string)
	Bind(name string, callback func(string) (string, error))
	SetMenuExtraText(name string)
	AddMenuExtra(name string)
	//	SetTitle(title string)
}

type mprotonHandle struct {
}

func New() mProtonApp {
	h := &mprotonHandle{}
	C.initialize()
	return h
}

func (handle *mprotonHandle) Run() {
	C.xmain()
}

func (handle *mprotonHandle) SetTitle(path string) {
	C.set_title(C.CString(path)) // TODO: cleanup mem
}

func (handle *mprotonHandle) SetMenuExtraText(text string) {
	C.set_menu_extra_text(C.CString(text))
}

func (handle *mprotonHandle) AddMenuExtra(text string) {
	C.add_menu_extra_item(C.CString(text))
}

func (handle *mprotonHandle) SetContentPath(path string) {
	C.add_content_path(C.CString(path))
}

func (handle *mprotonHandle) Bind(name string, callback func(string) (string, error)) {
	registerCallback(name, callback)
	C.add_script_message_handler(C.CString(name))
}
