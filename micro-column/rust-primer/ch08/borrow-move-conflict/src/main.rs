fn print_book(book: &String) {
    println!("正在阅读：{}", book);
}

fn take_book(book: String) {
    println!("这本书被拿走了：{}", book);
}

fn main() {
    let book = String::from("Rust 第一课");

    let borrowed = &book;

    take_book(book);

    print_book(borrowed); // 修正方法：挪到take_book函数调用之前
}
