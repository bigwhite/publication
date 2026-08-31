fn print_key(key: String) {
    println!("使用钥匙：{}", key);
}

fn main() {
    let key = String::from("大门钥匙");

    print_key(key);

    // 下面这行会报错：key 的所有权已经被移动到 print_key 函数里了
    println!("我还想再用一次钥匙：{}", key);
}
