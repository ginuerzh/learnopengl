package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

const (
	vertexShaderSource = `
#version 330 core
layout (location = 0) in vec3 aPos;
void main()
{
   gl_Position = vec4(aPos.x, aPos.y, aPos.z, 1.0);
}
` + "\x00"
	fragmentShaderSource = `
#version 330 core
out vec4 FragColor;
void main()
{
   FragColor = vec4(1.0f, 0.5f, 0.2f, 1.0f);
}
` + "\x00"
)

var (
	vertices1 = []float32{
		// first triangle
		-0.9, -0.5, 0.0, // left
		-0.0, -0.5, 0.0, // right
		-0.45, 0.5, 0.0, // top
	}
	vertices2 = []float32{
		// second triangle
		0.0, -0.5, 0.0, // left
		0.9, -0.5, 0.0, // right
		0.45, 0.5, 0.0, // top
	}
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

	program, err := newPragram(vertexShaderSource, fragmentShaderSource)
	if err != nil {
		log.Fatal(err)
	}

	var vao, vbo [2]uint32
	gl.GenVertexArrays(2, &vao[0])
	defer gl.DeleteVertexArrays(2, &vao[0])

	gl.GenBuffers(2, &vbo[0])
	defer gl.DeleteBuffers(2, &vbo[0])

	// first triangle setup
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo[0]) // VAO will not record this operation state, it can be called before the VAO binding
	gl.BindVertexArray(vao[0])
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices1)*4, gl.Ptr(vertices1), gl.STATIC_DRAW)
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 3*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	// no need to unbind at all as we directly bind a different VBO/VAO the next few lines
	// gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	// gl.BindVertexArray(0)

	// second triangle setup
	gl.BindVertexArray(vao[1])             // note that we bind to a different VAO now
	gl.BindBuffer(gl.ARRAY_BUFFER, vbo[1]) // and a different VBO
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices2)*4, gl.Ptr(vertices2), gl.STATIC_DRAW)
	// because the vertex data is tightly packed we can also specify 0 as the vertex attribute's stride to let OpenGL figure it out
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 0, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)
	// gl.BindBuffer(gl.ARRAY_BUFFER, 0)
	// not really necessary as well, but beware of calls that could affect VAOs while this one is bound
	// (like binding element buffer objects, or enabling/disabling vertex attributes)
	// gl.BindVertexArray(0)

	// uncomment this call to draw in wireframe polygons.
	// gl.PolygonMode(gl.FRONT_AND_BACK, gl.LINE)

	// render loop
	for !window.ShouldClose() {
		// input
		processInput(window)

		// render
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// draw our first triangle
		gl.UseProgram(program)
		// draw first triangle using the data from the first VAO
		gl.BindVertexArray(vao[0])
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		// then we draw the second triangle using the data from the second VAO
		gl.BindVertexArray(vao[1])
		gl.DrawArrays(gl.TRIANGLES, 0, 3)
		// no need to unbind it every time
		// gl.BindVertexArray(0);

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

func newPragram(vertexShaderSource string, fragmentShaderSource string) (uint32, error) {
	vertexShader, err := compileShader(vertexShaderSource, gl.VERTEX_SHADER)
	if err != nil {
		return 0, err
	}

	fragmentShader, err := compileShader(fragmentShaderSource, gl.FRAGMENT_SHADER)
	if err != nil {
		return 0, err
	}

	program := gl.CreateProgram()
	gl.AttachShader(program, vertexShader)
	gl.AttachShader(program, fragmentShader)
	gl.LinkProgram(program)

	var status int32
	gl.GetProgramiv(program, gl.LINK_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetProgramiv(program, gl.INFO_LOG_LENGTH, &logLength)
		logs := strings.Repeat("\x00", int(logLength+1))
		gl.GetProgramInfoLog(program, logLength, nil, gl.Str(logs))

		return 0, fmt.Errorf("failed to link program: %v", logs)
	}

	gl.DeleteShader(vertexShader)
	gl.DeleteShader(fragmentShader)

	return program, nil
}

func compileShader(source string, shaderType uint32) (uint32, error) {
	shader := gl.CreateShader(shaderType)
	csources, free := gl.Strs(source)
	gl.ShaderSource(shader, 1, csources, nil)
	free()
	gl.CompileShader(shader)

	var status int32
	gl.GetShaderiv(shader, gl.COMPILE_STATUS, &status)
	if status == gl.FALSE {
		var logLength int32
		gl.GetShaderiv(shader, gl.INFO_LOG_LENGTH, &logLength)
		logs := strings.Repeat("\x00", int(logLength+1))
		gl.GetShaderInfoLog(shader, logLength, nil, gl.Str(logs))

		return 0, fmt.Errorf("failed to compile %v : %v", source, logs)
	}

	return shader, nil
}
