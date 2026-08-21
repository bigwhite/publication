fn main() {
    let name = "Tony";
    let score = 95.5;

    println!("{}, 你的分数是 {}", name, score);
    println!("{name}, 你的分数是 {score}");
    println!("{0} 排名比 {1} 靠前，{0} 加油", "张三", "李四");
    println!("保留两位小数：{:.2}", score);
    println!("宽度对齐：[{:>8}]", "右对齐");
    println!("十六进制：{:x}，二进制：{:b}", 255, 255);
    println!("Debug 格式打印元组：{:?}", (1, "a", 3.0));

    let msg = format!("{} 的成绩是 {:.1} 分", name, score);
    println!("format! 生成的字符串：{}", msg);
}
