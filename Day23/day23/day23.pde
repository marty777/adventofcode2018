

public class Bot {
   public PVector  p;
   public float    r;
}


int min_x, min_y, min_z, max_x, max_y, max_z;
Bot[] bots = {};
int displaymode = 0;

void setup() {
    size(640,640);
    
    String[] lines = loadStrings("input.txt");
    for(int i = 0; i < lines.length; i++) {
      int startindex = 0;
      String px = lines[i].substring(lines[i].indexOf("<")+1, lines[i].indexOf(","));
      startindex = lines[i].indexOf(",")+1;
      String py = lines[i].substring(startindex, lines[i].indexOf(",",startindex) );
      startindex = lines[i].indexOf(",", startindex)+1;
      String pz = lines[i].substring(startindex, lines[i].indexOf(">",startindex) );
      String r = lines[i].substring(lines[i].indexOf("r=")+2, lines[i].length());
      println(px);
      println(py);
      println(pz);
      println(r);
      Bot bot = new Bot();
      bot.p = new PVector(Integer.parseInt(px), Integer.parseInt(py), Integer.parseInt(pz));
      bot.r = float(Integer.parseInt(r));
      bots = (Bot[])append(bots, bot);
      int minx = int(floor(bot.p.x - bot.r));
      int miny = int(floor(bot.p.y - bot.r));
      int minz = int(floor(bot.p.z - bot.r));
      int maxx = int(floor(bot.p.x + bot.r));
      int maxy = int(floor(bot.p.y + bot.r));
      int maxz = int(floor(bot.p.z + bot.r));
      if(minx < min_x) {
         min_x = minx;
      }
      if(miny < min_y) {
         min_y = miny;
      }
      if(minz < min_z) {
         min_z = minz;
      }
      if(maxx > max_x) {
         max_x = maxx;
      }
      if(maxy > max_y) {
         max_y = maxy;
      }
      if(maxz > max_z) {
         max_z = maxz;
      }
    }
}

void draw() {
   background(0);
   float scale = 1.0;
   int dx = max_x - min_x;
   int dy = max_y - min_y;
   int dz = max_z - min_z;
   scale = float(dx);
   if(float(dy) > scale) {
      scale = float(dy);
   }
   if(float(dz) > scale) {
      scale = float(dz); 
   }
   
   scale = 640/scale;
   
   println(bots.length);
   for(int i=0; i < bots.length; i++) {
     if(displaymode == 0) {
       stroke(255,0,0);
       point((bots[i].p.x - min_x) * scale, (bots[i].p.y - min_y) * scale);
       //fill(255,255,255,1);
       noFill();
       stroke(255,255,255,8);
       //ellipseMode(RADIUS);
       //ellipse((bots[i].p.x - min_x) * scale, (bots[i].p.y - min_y) * scale, bots[i].r*scale, bots[i].r*scale);
       // west-south
       line((bots[i].p.x - min_x - bots[i].r)*scale, (bots[i].p.y - min_y)*scale, (bots[i].p.x - min_x)*scale, (bots[i].p.y - min_y - bots[i].r)*scale);
       // south-east
       line((bots[i].p.x - min_x)*scale, (bots[i].p.y - min_y - bots[i].r)*scale, (bots[i].p.x - min_x + bots[i].r)*scale, (bots[i].p.y - min_y)*scale);
       // east-north
       line((bots[i].p.x - min_x + bots[i].r)*scale, (bots[i].p.y - min_y)*scale, (bots[i].p.x - min_x)*scale, (bots[i].p.y - min_y + bots[i].r)*scale);
       // north-west
       line((bots[i].p.x - min_x)*scale, (bots[i].p.y - min_y + bots[i].r)*scale, (bots[i].p.x - min_x - bots[i].r)*scale, (bots[i].p.y - min_y)*scale);
       
   } 
     else if(displaymode == 1) {
       stroke(255,0,0);
       point((bots[i].p.x - min_x) * scale, (bots[i].p.z - min_z) * scale);
       //fill(255,255,255,1);
       noFill();
       stroke(255,255,255,8);
       //ellipseMode(RADIUS);
       //ellipse((bots[i].p.x - min_x) * scale, (bots[i].p.z - min_z) * scale, bots[i].r*scale, bots[i].r*scale);
      // west-south
       line((bots[i].p.x - min_x - bots[i].r)*scale, (bots[i].p.z - min_z)*scale, (bots[i].p.x - min_x)*scale, (bots[i].p.z - min_z - bots[i].r)*scale);
       // south-east
       line((bots[i].p.x - min_x)*scale, (bots[i].p.z - min_z - bots[i].r)*scale, (bots[i].p.x - min_x + bots[i].r)*scale, (bots[i].p.z - min_z)*scale);
       // east-north
       line((bots[i].p.x - min_x + bots[i].r)*scale, (bots[i].p.z - min_z)*scale, (bots[i].p.x - min_x)*scale, (bots[i].p.z - min_z + bots[i].r)*scale);
       // north-west
       line((bots[i].p.x - min_x)*scale, (bots[i].p.z - min_z + bots[i].r)*scale, (bots[i].p.x - min_x - bots[i].r)*scale, (bots[i].p.z - min_z)*scale);
        
      }
      else if(displaymode == 2) {
       stroke(255,0,0);
       point((bots[i].p.y - min_y) * scale, (bots[i].p.z - min_z) * scale);
       //fill(255,255,255,1);
       noFill();
       stroke(255,255,255,8);
       //ellipseMode(RADIUS);
       //ellipse((bots[i].p.y - min_y) * scale, (bots[i].p.z - min_z) * scale, bots[i].r*scale, bots[i].r*scale);
       // west-south
       line((bots[i].p.y - min_y - bots[i].r)*scale, (bots[i].p.z - min_z)*scale, (bots[i].p.y - min_y)*scale, (bots[i].p.z - min_z - bots[i].r)*scale);
       // south-east
       line((bots[i].p.y - min_y)*scale, (bots[i].p.z - min_z - bots[i].r)*scale, (bots[i].p.y - min_y + bots[i].r)*scale, (bots[i].p.z - min_z)*scale);
       // east-north
       line((bots[i].p.y - min_y + bots[i].r)*scale, (bots[i].p.z - min_z)*scale, (bots[i].p.y - min_y)*scale, (bots[i].p.z - min_z + bots[i].r)*scale);
       // north-west
       line((bots[i].p.y - min_y)*scale, (bots[i].p.z - min_z + bots[i].r)*scale, (bots[i].p.y - min_y - bots[i].r)*scale, (bots[i].p.z - min_z)*scale);
        
     }
 }
   
}

void keyPressed() {
   if(key == 'm' || key == 'M') {
     displaymode = (displaymode + 1) % 3;
   }
}