fn main() {
    let mut book = String::from("Rust");

    let r1 = &mut book;
    let r2 = &mut book;

    println!("{}, {}", r1, r2);
}
