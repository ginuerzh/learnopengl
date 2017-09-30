package main

import (
	"log"
	"math"
	"runtime"

	"github.com/ginuerzh/learnopengl/utils/shader"
	"github.com/ginuerzh/learnopengl/utils/texture"

	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.2/glfw"
	"github.com/go-gl/mathgl/mgl32"
)

const (
	screenWidth  = 800
	screenHeight = 600
)

var (
	vertices = []float32{
		-0.5, -0.5, -0.5, 0.0, 0.0,
		0.5, -0.5, -0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 0.0,

		-0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 1.0,
		-0.5, 0.5, 0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,

		-0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, -0.5, 1.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, 0.5, 1.0, 0.0,

		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, 0.5, 0.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,

		-0.5, -0.5, -0.5, 0.0, 1.0,
		0.5, -0.5, -0.5, 1.0, 1.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		0.5, -0.5, 0.5, 1.0, 0.0,
		-0.5, -0.5, 0.5, 0.0, 0.0,
		-0.5, -0.5, -0.5, 0.0, 1.0,

		-0.5, 0.5, -0.5, 0.0, 1.0,
		0.5, 0.5, -0.5, 1.0, 1.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		0.5, 0.5, 0.5, 1.0, 0.0,
		-0.5, 0.5, 0.5, 0.0, 0.0,
		-0.5, 0.5, -0.5, 0.0, 1.0,
	}
	cubePositions = []mgl32.Vec3{
		{0.0, 0.0, 0.0},
		{2.0, 5.0, -15.0},
		{-1.5, -2.2, -2.5},
		{-3.8, -2.0, -12.3},
		{2.4, -0.4, -3.5},
		{-1.7, 3.0, -7.5},
		{1.3, -2.0, -2.5},
		{1.5, 2.0, -2.5},
		{1.5, 0.2, -1.5},
		{-1.3, 1.0, -1.5},
	}
)

func init() {
	runtime.LockOSThread()
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
	window, err := glfw.CreateWindow(screenWidth, screenHeight, "LearnOpenGL", nil, nil)
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

	shader, err := shader.NewShader("7.1.camera.vs", "7.1.camera.fs")
	if err != nil {
		log.Fatal(err)
	}

	var vao, vbo uint32
	gl.GenVertexArrays(1, &vao)
	defer gl.DeleteVertexArrays(1, &vao)

	gl.GenBuffers(1, &vbo)
	defer gl.DeleteBuffers(1, &vbo)

	// bind the Vertex Array Object first, then bind and set vertex buffer(s), and then configure vertex attributes(s).
	gl.BindVertexArray(vao)

	gl.BindBuffer(gl.ARRAY_BUFFER, vbo)
	gl.BufferData(gl.ARRAY_BUFFER, len(vertices)*4, gl.Ptr(vertices), gl.STATIC_DRAW)

	// position attribute
	gl.VertexAttribPointer(0, 3, gl.FLOAT, false, 5*4, gl.PtrOffset(0))
	gl.EnableVertexAttribArray(0)

	// texture coord attribute
	gl.VertexAttribPointer(1, 2, gl.FLOAT, false, 5*4, gl.PtrOffset(3*4))
	gl.EnableVertexAttribArray(1)

	// texture 1
	texture1 := texture.NewTexture2D()
	texture1.Use()
	// set the texture wrapping parameters
	texture1.SetParameter(gl.TEXTURE_WRAP_S, gl.REPEAT)
	texture1.SetParameter(gl.TEXTURE_WRAP_T, gl.REPEAT)
	// set texture filtering parameters
	texture1.SetParameter(gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	texture1.SetParameter(gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	// load image, create texture and generate mipmaps
	image, err := texture1.Load("../../resources/textures/container.jpg", false, false)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("container.jpg", image.Rect, image.Stride, len(image.Pix))

	texture2 := texture.NewTexture2D()
	texture2.Use()
	texture2.SetParameter(gl.TEXTURE_WRAP_S, gl.REPEAT)
	texture2.SetParameter(gl.TEXTURE_WRAP_T, gl.REPEAT)
	texture2.SetParameter(gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	texture2.SetParameter(gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	image, err = texture2.Load("../../resources/textures/awesomeface.jpg", false, true)
	if err != nil {
		log.Fatal(err)
	}
	log.Println("awesomeface.jpg", image.Rect, image.Stride, len(image.Pix))

	// note that this is allowed,
	// the call to glVertexAttribPointer registered VBO as the vertex attribute's bound vertex buffer object
	// so afterwards we can safely unbind
	gl.BindBuffer(gl.ARRAY_BUFFER, 0)

	// You can unbind the VAO afterwards so other VAO calls won't accidentally modify this VAO, but this rarely happens. Modifying other
	// VAOs requires a call to glBindVertexArray anyways so we generally don't unbind VAOs (nor VBOs) when it's not directly necessary.
	// gl.BindVertexArray(0)

	shader.Use()
	if err := shader.SetUniformName("texture1", 0); err != nil {
		log.Fatal(err)
	}
	if err := shader.SetUniformName("texture2", 1); err != nil {
		log.Fatal(err)
	}

	projection := mgl32.Perspective(mgl32.DegToRad(45.0), float32(screenWidth)/float32(screenHeight), 0.1, 100.0)
	// projection = mgl32.Ortho(0.0, 800.0, 0.0, 600.0, 0.1, 100.0)
	if err := shader.SetUniformMatrixName("projection", false, projection); err != nil {
		log.Fatal(err)
	}

	gl.Enable(gl.DEPTH_TEST)

	// render loop
	for !window.ShouldClose() {
		// input
		processInput(window)

		// render
		gl.ClearColor(0.2, 0.3, 0.3, 1.0)
		gl.Clear(gl.DEPTH_BUFFER_BIT | gl.COLOR_BUFFER_BIT)

		// bind textures on corresponding texture units
		gl.ActiveTexture(gl.TEXTURE0)
		texture1.Use()
		gl.ActiveTexture(gl.TEXTURE1)
		texture2.Use()

		radius := 10.0
		camX := math.Sin(glfw.GetTime()) * radius
		camZ := math.Cos(glfw.GetTime()) * radius
		view := mgl32.LookAtV(
			mgl32.Vec3{float32(camX), 0, float32(camZ)},
			mgl32.Vec3{},
			mgl32.Vec3{0, 1, 0},
		)
		if err := shader.SetUniformMatrixName("view", false, view); err != nil {
			log.Fatal(err)
		}

		for i, pos := range cubePositions {
			angle := 20.0 * float32(i)

			model := mgl32.HomogRotate3D(float32(mgl32.DegToRad(angle)), mgl32.Vec3{0.5, 1.0, 0.0})
			model = mgl32.Translate3D(pos.Elem()).Mul4(model)
			if err := shader.SetUniformMatrixName("model", false, model); err != nil {
				log.Fatal(err)
			}
			gl.DrawArrays(gl.TRIANGLES, 0, 36)
		}

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
