package demo;

public class Demo {
    public static void main(String[] args) {
        int height = 300;
        try {
            HeightInput.checkHeight(height);
        } catch (HeightOutOfBound e) {
             System.out.printf("%s %s\n", "Are you a real human?", e);
        }
    }
}
