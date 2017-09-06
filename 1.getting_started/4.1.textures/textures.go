package main

import (
	"log"

	"github.com/ginuerzh/learnopengl/utils/shader"
	"github.com/ginuerzh/learnopengl/utils/texture"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
)

var (
	vertices = []float32{
		// positions   // colors      // texture coords
		0.5, 0.5, 0.0, 1.0, 0.0, 0.0, 1.0, 1.0, // top right
		0.5, -0.5, 0.0, 0.0, 1.0, 0.0, 1.0, 0.0, // bottom right
		-0.5, -0.5, 0.0, 0.0, 0.0, 1.0, 0.0, 0.0, // bottom left
		-0.5, 0.5, 0.0, 1.0, 1.0, 0.0, 0.0, 1.0, // top left
	}
	indices = []uint32{
		0, 1, 3, // first triangle
		1, 2, 3, // second triangle
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
	// there is no need to explicitly create and bind a VAO when using the compatibility profile
	// see https://www.opengl.org/discussion_boards/showthread.php/199916-vertex-array-and-buffer-objects?p=1288280&viewfull=1#post1288280
	// glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCompatProfile)

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

	shader, err := shader.NewShader("4.1.texture.vs", "4.1.texture.fs")
	if err != nil {
		log.Fatal(err)
	}

	var vao, vbo, ebo uint32
	gl.GenVertexArrays(1, &vao)
	defer gl.DeleteVertexArrays(1, &vao)

	gl.GenBuffers(1, &vbo)
	defer gl.DeleteBuffers(1, &vbo)

	gl.GenBuffers(1, &ebo)
	defer gl.DeleteBuffers(1, &ebo)

	// bind the Vertex Array Object first, then bind and set vertex buffer(s), and then configure vertex attributes(s).
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	gl.BindBuffer(gl.ELEMENT_ARRAY_BUFFER, ebo)
	gl.BufferData(gl.ELEMENT_ARRAY_BUFFER, len(indices)*4, gl.Ptr(indices), gl.STATIC_DRAW)

	// position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// color attribute
	gl.VertexAttribPointer(1, 3, gl.FLOAT, false, 8*4, gl.PtrOffset(3*4)) // note that the offset is 12 in BYTEs.
	gl.EnableVertexAttribArray(1)

	// texture coord attribute
	gl.VertexAttribPointer(2, 2, gl.FLOAT, false, 8*4, gl.PtrOffset(6*4))
	gl.EnableVertexAttribArray(2)

	texture := texture.NewTexture2D()
	// all upcoming GL_TEXTURE_2D operations now have effect on this texture object
	texture.Use()
	// set the texture wrapping parameters
	texture.SetParameter(gl.TEXTURE_WRAP_S, gl.REPEAT)
	texture.SetParameter(gl.TEXTURE_WRAP_T, gl.REPEAT)
	// set texture filtering parameters
	texture.SetParameter(gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	texture.SetParameter(gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	// load image, create texture and generate mipmaps
	image, err := texture.Load("../../resources/textures/container.jpg")
	if err != nil {
		log.Fatal(err)
	}
	log.Println(image.Rect, image.Stride)

	// note that this is allowed,
	// the call to glVertexAttribPointer registered VBO as the vertex attribute's bound vertex buffer object
	// so afterwards we can safely unbind
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// You can unbind the VAO afterwards so other VAO calls won't accidentally modify this VAO, but this rarely happens. Modifying other
	// VAOs requires a call to glBindVertexArray anyways so we generally don't unbind VAOs (nor VBOs) when it's not directly necessary.
	// gl.BindVertexArray(0)

	shader.Use()

	// render loop
	for !window.ShouldClose() {
		// input
		processInput(window)

		// render
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.COLOR_BUFFER_BIT)

		// seeing as we only have a single VAO there's no need to bind it every time.
		// gl.BindVertexArray(vao)
		gl.DrawElements(gl.TRIANGLES, 6, gl.UNSIGNED_INT, gl.PtrOffset(0))
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
