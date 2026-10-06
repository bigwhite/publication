use std::fmt::Display;

fn print_it<T: Display + 'static>(t: T) {
    println!("打印：{}", t);
}

fn main() {
    let owned = String::from("这是一个拥有所有权的 String");
    // owned 本身完全是在函数调用这一刻才创建的，运行没多久，
    // 但它依然满足 T: 'static —— 因为它不包含任何"借来的、可能过期的"数据。
    print_it(owned);

    print_it(42);
    print_it("字符串字面量");
}
