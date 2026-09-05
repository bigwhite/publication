fn print_book(who: &str, book: &String) {
    println!("{} 正在阅读：{}", who, book);
}

fn main() {
    let book = String::from("Rust 第一课");

    let reader1 = &book;
    let reader2 = &book;
    let reader3 = &book;

    print_book("张三", reader1);
    print_book("李四", reader2);
    print_book("王五", reader3);

    println!("book 本身也依然可用：{}", book);
}
