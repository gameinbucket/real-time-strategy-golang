#version 330 core

uniform mat4 u_mvp;

layout (location = 0) in vec2 a_position;
layout (location = 1) in vec2 a_texture;

out vec2 v_texture;

void main()
{
    v_texture = a_texture;

    gl_Position = u_mvp * vec4(a_position, 0.0, 1.0);
}
