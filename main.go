package main

import (
	"log"
	//"sort"
	"fmt"
	"os/exec"
	"time"

	"runtime"

	"github.com/donomii/goof"

	"github.com/donomii/nuklear-templates"
	"github.com/mitchellh/go-ps"

	"github.com/golang-ui/nuklear/nk"

	"github.com/xlab/closer"

	"github.com/go-gl/glfw/v3.2/glfw"
)

var procs map[int]ps.Process
var procsstr []string

func getProcs() map[int]ps.Process {
	procs, _ := ps.Processes()
	procHash := map[int]ps.Process{}
	procsstr = []string{}
	for _, v := range procs {
		procHash[v.Pid()] = v
	}
	return procHash
}
func updateProcs() {

	procs = getProcs()
	for _, vv := range procs {
		//log.Println(extendedPS(vv.Pid()))
		procsstr = append(procsstr, extendedPS(vv.Pid()))
	}

}

func extendedPS(pid int) string {
	if runtime.GOOS == "windows" {
		cmd := exec.Command("tasklist.exe", "/fo", "csv", "/nh")
		//Only compiles on windows?
		//cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
		out, err := cmd.Output()
		if err != nil {
			return ""
		}
		return string(out)
	} else {
		//if runtime.GOOS == "darwin" {
		//if runtime.GOOS == "linux" {
		out := goof.Chomp(goof.Chomp(goof.Command("ps", []string{"-o", "command=", "-p", fmt.Sprintf("%v", pid)})))
		return out
	}
	return ""
}

func main() {

	runtime.LockOSThread()
	runtime.GOMAXPROCS(1)
	if err := glfw.Init(); err != nil {
		closer.Fatalln(err)
	}
	log.Println("Loading processes")
	updateProcs()
	log.Println("Load complete")
	win, ctx := nktemplates.StartNuke()

	pane1 := func() {
		//A horizontal button bar.  Fires the callback when clicked
		nktemplates.ButtonBar(ctx, []string{"Kill", "Continuous Kill"}, func(i int, s string) { log.Println(s) })
	}
	pane2 := func() {
		/*
			for _, vv := range procsps {
				procs = append(procs, vv.Executable())
			}

			sort.Strings(procs)
		*/

		//nk.NkLayoutRowDynamic(ctx, 40, 1)
		//{

		for _, name := range procsstr {
			//node := vv.SubNodes[i]

			//time.Sleep(1 * time.Millisecond)
			//if nk.NkButtonLabel(ctx, name.Executable()) > 0 {

			if nk.NkButtonLabel(ctx, name) > 0 {
				//Kill process or whatever
			}
		}
		//}
	}

	exitC := make(chan struct{}, 1)
	doneC := make(chan struct{}, 1)
	closer.Bind(func() {
		close(exitC)
		<-doneC
	})

	for {
		select {
		case <-exitC:
			nk.NkPlatformShutdown()
			glfw.Terminate()
			close(doneC)
			return
		default:
			if win.ShouldClose() {
				close(exitC)
				continue
			}
			glfw.PollEvents()

			nktemplates.LeftCol(win, ctx, nil, pane1, pane2, pane1)
			//TkRatWin(win, ctx, state, pane1, pane2, pane3)
			time.Sleep(100 * time.Millisecond)
		}
	}
}
