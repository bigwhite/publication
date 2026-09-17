fn first_word(s: &str) -> &str {
    let bytes = s.as_bytes();
    for (i, &item) in bytes.iter().enumerate() {
        if item == b' ' {
            return &s[0..i];
        }
    }
    s
}

fn announce_and_return_first<'a>(x: &'a str, y: &str) -> &'a str {
    println!("即将比较的另一个字符串是：{}", y);
    x
}

fn main() {
    let sentence = String::from("Rust 第一课 事不过三");
    println!("第一个单词：{}", first_word(&sentence));

    let s1 = String::from("Rust 第一课");
    let result;
    {
        let s2 = String::from("这是一个临时字符串");
        result = announce_and_return_first(s1.as_str(), s2.as_str());
        println!("函数内返回值：{}", result);
    }
    println!("函数外依然可用：{}", result);
}
