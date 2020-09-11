package stringer;

public class foo implements Stringer {
        int i;
        public String String() {
                return Integer.toString(i);
        }
}
