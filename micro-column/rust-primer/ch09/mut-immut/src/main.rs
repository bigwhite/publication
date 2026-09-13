fn main() {
    let mut book = String::from("Rust");

    let r1 = &book;
    let r2 = &mut book;

    println!("{}, {}", r1, r2);
}
