package main

import (
	"log"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

func init() {
	log.SetFlags(log.Lshortfile | log.LstdFlags)
}

func main() {
	// glfw: initialize
	if err := glfw.Init(); err != nil {
		log.Fatal(err)
	}
	// glfw: terminate, clearing all previously allocated GLFW resources.
	defer glfw.Terminate()

	// glfw: configure
	glfw.WindowHint(glfw.ContextVersionMajor, 3)
	glfw.WindowHint(glfw.ContextVersionMinor, 3)
	glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)

	// glfw window creation
	window, err := glfw.CreateWindow(800, 600, "LearnOpenGL", nil, nil)
	if err != nil {
		log.Fatal(err)
	}
	defer window.Destroy()

	window.MakeContextCurrent()
	window.SetFramebufferSizeCallback(frameBufferSizeCallback)

	// gl: initializes the OpenGL bindings by loading the function pointers
	// (for each OpenGL function) from the active OpenGL context.
	if err := gl.Init(); err != nil {
		log.Fatal(err)
	}

	// render loop
	for !window.ShouldClose() {
		// input
		processInput(window)

		// render
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// glfw: swap buffers and poll IO events (keys pressed/released, mouse moved etc.)
		window.SwapBuffers()
		glfw.PollEvents()
	}
}

// process all input: query GLFW whether relevant keys are pressed/released this frame and react accordingly
func processInput(w *glfw.Window) {
	if w.GetKey(glfw.KeyEscape) == glfw.Press {
		w.SetShouldClose(true)
	}
}

// glfw: whenever the window size changed (by OS or user resize) this callback function executes
func frameBufferSizeCallback(w *glfw.Window, width int, height int) {
	// make sure the viewport matches the new window dimensions; note that width and
	// height will be significantly larger than specified on retina displays.
	gl.Viewport(0, 0, int32(width), int32(height))
	// log.Printf("frameBufferSizeCallback (%d, %d)", width, height)
}
