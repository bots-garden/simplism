package server

import "sync"

var currentSimplismProcess = SimplismProcess{}

// TODO:
// - add methods (getter and setter) to SimplismProcess

var protection = sync.Mutex{}
var wasmServices = make(map[string]SimplismProcess) // Map or Slice ? 🤔
