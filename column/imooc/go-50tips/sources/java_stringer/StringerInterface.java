package stringer;

interface Stringer {
    String String();
}


public class StringerInterface {
	public static String concat(Stringer a, Stringer b) {
		return a.String()+b.String();
	}

	public static void main(String args[]){
		bar b = new bar();
		b.s = "hello";
		foo f = new foo();
		f.i = 5;
		System.out.println(concat(b, f));
	}
}


