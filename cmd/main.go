package main

import (
	"github.com/NeerajRijhwani/code-editor/internal/editor"
	"log"
	"os"
	"runtime/debug"
)

func main() {
	file, err := os.OpenFile("debug.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	defer file.Close()
	log.SetOutput(file)

	defer func() {
		if r := recover(); r != nil {
			log.Printf("[PANIC RECOVERED] Value: %v\nStack Trace:\n%s", r, debug.Stack())
		}
	}()

	log.Println("Editor Started")
	editor, err := editor.InitEditor()

	editor.Renderer.Show()
	if err != nil {
		log.Fatal(err)
	}
	editor.Run()

}
