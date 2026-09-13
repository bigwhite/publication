fn add_suffix(book: &mut String) {
    book.push_str("（第一课）");
}

fn main() {
    let mut book = String::from("Rust");
    println!("修改前：{}", book);

    add_suffix(&mut book);

    println!("修改后：{}", book);
}
