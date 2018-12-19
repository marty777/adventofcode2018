#ifdef GL_ES
precision mediump float;
precision mediump int;
#endif

#define PROCESSING_TEXTURE_SHADER

uniform sampler2D texture;

//uniform float texWidth;
//uniform float texHeight;
uniform float pixelWidth;
uniform float pixelHeight;

varying vec4 vertTexCoord;

bool neighborTest(vec2 pos, int field, int min) {
	int count = 0;
	for(float i = pos.x - pixelWidth; i <= pos.x + pixelWidth; i+=pixelWidth) {
		for( float j = pos.y - pixelHeight; j <= pos.y + pixelHeight; j+=pixelHeight) {
			if(i == pos.x && j == pos.y || i < 0.0 || j < 0.0 || i > 1.0 || j > 1.0) {
				continue;
			}
			if(field == 0) {
					if(texture2D(texture, vec2(i,j)).x > 0.0) {
						count++;
					}
			}
			else if(field == 1) {
					if(texture2D(texture, vec2(i,j)).y > 0.0) {
						count++;
					}
			}
			else {
					if(texture2D(texture, vec2(i,j)).z > 0.0) {
						count++;
					}
			}
		}
	}
	if(count >= min) {
		return true;
	}
	return false;
}


void main() {
	
  
  vec2 coords = vec2(vertTexCoord.x, vertTexCoord.y);
  gl_FragColor = texture2D(texture, coords);
  
  // open field to trees
  if(texture2D(texture, coords).x > 0.0 && neighborTest(coords, 1, 3)) {
	gl_FragColor = vec4(0.0,1.0,0.0,1.0);
  }
  // trees to lumberyard
  else if(texture2D(texture,coords).y > 0.0 && neighborTest(coords, 2, 3)) {
	gl_FragColor = vec4(0.0,0.0,1.0,1.0);
  }
  // lumberyard to open
  else if(texture2D(texture, coords).z > 0.0 && (!neighborTest(coords, 2, 1) || !neighborTest(coords, 1, 1))) {
		gl_FragColor = vec4(1.0,0.0,0.0,1.0);
  }
  
 
  //gl_FragColor = vec4(1.0,0.0,0.0,1.0);
}