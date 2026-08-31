fn print_key(key: String) {
    println!("使用钥匙：{}", key);
}

fn main() {
    let key = String::from("大门钥匙");

    print_key(key.clone());

    println!("我还想再用一次钥匙：{}", key);
}
