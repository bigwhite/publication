fn main() {
    let mut book = String::from("Rust");

    let r1 = &book;
    let r2 = &book;
    println!("只读借用：{}, {}", r1, r2);
    // r1、r2 的最后一次使用就在上面这一行，NLL 之下它们的"生命"到这里就结束了

    let r3 = &mut book;
    r3.push_str("（第一课）");
    println!("可变借用之后：{}", r3);
}
