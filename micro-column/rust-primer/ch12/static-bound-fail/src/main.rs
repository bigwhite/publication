use std::fmt::Display;

fn print_it<T: Display + 'static>(t: T) {
    println!("打印：{}", t);
}

fn main() {
    let owned = String::from("局部字符串");
    let borrowed: &str = owned.as_str();
    print_it(borrowed); // borrowed 的生命周期受限于 owned，不满足 'static
}
