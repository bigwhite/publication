fn main() {
    // 字符串字面量：写死在代码里，类型是 &str
    let literal: &str = "Rust 第一课";
    println!("字面量：{}（类型是 &str）", literal);

    // String：拥有所有权、可增长
    let mut owned: String = String::from(literal);
    owned.push_str("，事不过三");
    owned.push('！');
    println!("String：{}", owned);

    // &String 会被自动解引用（deref coercion）成 &str，两者可以无缝互通
    fn print_it(s: &str) {
        println!("统一接收 &str：{}", s);
    }
    print_it(&owned);
    print_it(literal);

    // 转义字符与原始字符串
    let escaped = "第一行\n第二行\t制表符\\反斜杠\"引号";
    println!("转义写法：{}", escaped);

    let raw = r"C:\Users\tony\rust-first-course";
    println!("原始字符串（不转义）：{}", raw);

    // 多行字符串
    let multi = "第一行
第二行";
    println!("多行字符串：{}", multi);
}
