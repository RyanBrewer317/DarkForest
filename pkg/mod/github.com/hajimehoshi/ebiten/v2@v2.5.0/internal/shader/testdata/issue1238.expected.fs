vec4 F0(in vec4 l0);

vec4 F0(in vec4 l0) {
	if (true) {
		return l0;
	}
	return l0;
}

void main(void) {
	gl_FragColor = F0(gl_FragCoord);
}
