#version 330 core

in vec3 ourPos;
out vec4 FragColor;

void main()
{
    FragColor = vec4(ourPos, 1.0f);
}