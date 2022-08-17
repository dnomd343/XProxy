package main

import (
    "os/exec"
)

var subProcess []*Process

func killSub(sub *Process) {
    defer func() {
        recover()
    }()

}

func waitSub(sub *exec.Cmd) {
    defer func() {
        recover()
    }()
}

func exit() {

}
