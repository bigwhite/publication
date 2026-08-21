fn main() {
    let raw = "  Rust 第一课  ";

    println!("trim 去除首尾空白：[{}]", raw.trim());
    println!("是否包含子串：{}", raw.contains("第一课"));
    println!("替换：{}", raw.trim().replace("第一课", "事不过三"));
    println!("转大写：{}", "rust".to_uppercase());
    println!("转小写：{}", "RUST".to_lowercase());

    let parts: Vec<&str> = "a,b,c".split(',').collect();
    println!("split 拆分结果：{:?}", parts);
}
