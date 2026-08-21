fn add_one(x: u8) -> u8 {
    x + 1
}

fn main() {
    let a: u8 = 255;
    let b = add_one(a);
    println!("{}", b);
}
