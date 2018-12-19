//day18

PShader day18Shader;
// alternating buffers
PGraphics renderGraphics1;
PGraphics renderGraphics2;
PImage startImage;
int inputWidth;
int inputHeight;

int stepCount = 0;

String inputfile = "input.txt";


String[] readFile(String filepath) {
  String[] lines = loadStrings(filepath);
  println("there are " + lines.length + " lines");
  for (int i = 0 ; i < lines.length; i++) {
    println(lines[i]);
  }
  return lines;
}

void setup() {
    size(100,100,P2D);
    String[] lines = loadStrings(inputfile);
    inputWidth = lines[0].length();
    inputHeight = lines.length;
    //size(inputWidth, inputHeight, P2D);
    println("Input: " + inputWidth + " x " + inputHeight);
    
    day18Shader = loadShader("sampleShader.glsl", "vertexShader.glsl");
    day18Shader.set("pixelWidth", 1.0/inputWidth);
    day18Shader.set("pixelHeight", 1.0/inputHeight);
    println("Pixel dimensions (ratio):"+ 1.0/inputWidth + "," + 1.0/inputHeight);
    startImage = createImage(inputWidth, inputHeight, RGB);
    renderGraphics1 = createGraphics(inputWidth, inputHeight, P2D);
    renderGraphics2 = createGraphics(inputWidth, inputHeight, P2D);
    
    renderGraphics1.shader(day18Shader);
    renderGraphics2.shader(day18Shader);
    
    for(int i = 0; i < inputWidth; i++) {
       for(int j = 0; j < inputHeight; j++) {
         int val = 0;
         if(lines[j].charAt(i) == '|') {
            val = 1;
         } else if(lines[j].charAt(i) == '#') {
            val = 2;
         }  
         startImage.pixels[i + inputWidth*j] = color(val == 0 ? 255 : 0, val == 1 ? 255 : 0, val == 2 ? 255 : 0);
       } 
    }
    
}

void runSimulation() {
    
    int current_steps = 0;
    while(stepCount < 1000000000 && current_steps < 100000) {
        if(stepCount == 0) {
          // render the input image thru the shader
          renderGraphics1.beginDraw();
          renderGraphics1.image(startImage,0,0);
          renderGraphics1.endDraw(); 
          //shader(day18Shader);
          //image(renderGraphics1,0,0);
          //resetShader();
          //return;
        }
        else {
          if(stepCount % 2 == 0) {
             // render buffer 2 into buffer 1 using the shader
             renderGraphics1.beginDraw();
             renderGraphics1.image(renderGraphics2,0,0);
             //renderGraphics1.resetShader();
             renderGraphics1.endDraw(); 
             image(startImage,0,0);
          } 
          else {
            // render buffer 1 into buffer 2 using the shader
             renderGraphics2.beginDraw();
             renderGraphics2.image(renderGraphics1,0,0);
             //renderGraphics2.resetShader();
             renderGraphics2.endDraw();
          }
          
        }
         stepCount++;
          current_steps++;
        
    }
    println(stepCount);
         background(0);
         image(renderGraphics1,0,0); 
         println(red(renderGraphics1.get(0,0)), green(renderGraphics1.get(0,0)), blue(renderGraphics1.get(0,0)));
}

void draw() {
  //image(startImage,0,0);
  runSimulation();
  //runSimulation();
}