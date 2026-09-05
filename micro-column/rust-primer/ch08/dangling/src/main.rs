/*
fn dangle() -> &String {
    let s = String::from("hello");
    &s
}
*/

fn dangle() -> &'static String {
    let s = String::from("hello");
    &s
}

fn main() {
    let reference = dangle();
    println!("{}", reference);
}
