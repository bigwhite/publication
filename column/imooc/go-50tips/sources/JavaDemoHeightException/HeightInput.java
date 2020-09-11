package demo;

public class HeightInput {
    public static void checkHeight(int height) throws HeightOutOfBound {
        if(height>20 && height<300){
            System.out.print("ok");
        }else{
            throw new HeightOutOfBound();
        }
    }
}
