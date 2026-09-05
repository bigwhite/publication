fn calculate_length(s: &String) -> usize {
    s.len()
}

fn main() {
    let book = String::from("Rust 第一课");

    let len = calculate_length(&book);

    println!("《{}》这本书的标题长度是 {} 字节", book, len);
}
