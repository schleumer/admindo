package main

import (
	"bitbucket.org/kardianos/osext"
	"bufio"
	"fmt"
	"github.com/lxn/win"
	"os"
	"os/exec"
	"strings"
	"syscall"
	"unsafe"
)

func SetRegString(h win.HKEY, key string, value string) {
	win.RegSetValueEx(h, syscall.StringToUTF16Ptr(key), 0, win.REG_SZ, (*byte)(unsafe.Pointer(syscall.StringToUTF16Ptr(value))), 32)
}

func ReadRegString(h win.HKEY, key string) (value string) {
	var typ uint32
	var buf [74]uint16
	n := uint32(len(buf))
	win.RegQueryValueEx(h, syscall.StringToUTF16Ptr(key), nil, &typ, (*byte)(unsafe.Pointer(&buf[0])), &n)
	return syscall.UTF16ToString(buf[:])
}

func main() {
	var execpath, _ = osext.Executable()

	var h win.HKEY
	var h2 win.HKEY
	win.RegOpenKeyEx(win.HKEY_CURRENT_USER, syscall.StringToUTF16Ptr("SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\AppCompatFlags\\Layers"), 0, syscall.KEY_WRITE, &h)
	win.RegOpenKeyEx(win.HKEY_CURRENT_USER, syscall.StringToUTF16Ptr("SOFTWARE\\Microsoft\\Windows NT\\CurrentVersion\\AppCompatFlags\\Layers"), 0, syscall.KEY_READ, &h2)

	var val = ReadRegString(h2, execpath)

	if val != "~ RUNASADMIN" {
		fmt.Println("Putting executable to run as Administrator, now restart your command to be a trully administrator master race")
		SetRegString(h, execpath, "~ RUNASADMIN")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
		return
	}
	defer win.RegCloseKey(h)

	var args = make([]string, 0)
	var file = "powershell.exe"
	var cmdArgs = ""
	var argsList = make([]string, 0)

	if len(os.Args) >= 2 {
		args = os.Args[1:]
		file = args[0]
	}

	if len(args) > 1 {
		argsList = args[1:]
	}

	daShit := " -NoNewWindow "

	if len(argsList) > 0 {
		daShit += " -ArgumentList "
		cmdArgs = "\"" + strings.Join(argsList, " ") + "\""
	}

	cmdLine := "Start-Process " + file + daShit

	p := exec.Command("powershell", cmdLine, cmdArgs)
	p.Dir = "C:\\"
	p.Stdout = os.Stdout
	p.Stdin = os.Stdin
	p.SysProcAttr = &syscall.SysProcAttr{
		HideWindow: true,
	}
	p.Start()
	p.Wait()
}
