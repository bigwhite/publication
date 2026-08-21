fn main() {
    let s = String::from("Rust中文");

    println!("len()（字节数）        : {}", s.len());
    println!("chars().count()（字符数）: {}", s.chars().count());

    for (i, c) in s.chars().enumerate() {
        println!("  第 {} 个字符：{}", i, c);
    }

    let third_char = s.chars().nth(4);
    println!("第 5 个字符（.chars().nth(4)）：{:?}", third_char);
}
