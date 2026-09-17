fn longest<'a>(x: &'a str, y: &'a str) -> &'a str {
    if x.len() > y.len() {
        x
    } else {
        y
    }
}

fn main() {
    let s1 = String::from("Rust 第一课");
    let result;

    {
        let s2 = String::from("事不过三");
        result = longest(s1.as_str(), s2.as_str());
    }

    println!("最长的字符串是：{}", result);
}
