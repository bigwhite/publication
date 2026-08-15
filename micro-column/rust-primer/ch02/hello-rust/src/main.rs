fn main() {
    let name = "Rustacean";
    greet(name);
}

fn greet(name: &str) {
    println!("Hello, {}! 欢迎开启 Rust 第一课。", name);
}
