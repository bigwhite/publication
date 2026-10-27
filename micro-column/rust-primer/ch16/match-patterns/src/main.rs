fn describe_number(n: i32) -> String {
    match n {
        0 => "零".to_string(),
        1 | 2 | 3 => "一到三之间".to_string(), // | 多重匹配
        4..=9 => "四到九之间".to_string(),     // 范围匹配
        n if n < 0 => format!("负数：{}", n),  // 守卫（guard）
        n @ 10..=99 => format!("是一个两位数：{}", n), // @ 绑定，既匹配范围又拿到具体值
        _ => "太大了，懒得分类".to_string(),   // 通配符兜底
    }
}

fn main() {
    for n in [-5, 0, 2, 7, 42, 999] {
        println!("{} -> {}", n, describe_number(n));
    }

    // 解构：元组和结构体
    let point = (3, 0);
    match point {
        (0, 0) => println!("原点"),
        (x, 0) => println!("在 x 轴上，x = {}", x),
        (0, y) => println!("在 y 轴上，y = {}", y),
        (x, y) => println!("普通点：({}, {})", x, y),
    }

    struct Point {
        x: i32,
        y: i32,
    }
    let p = Point { x: 0, y: 7 };
    match p {
        Point { x: 0, y } => println!("在 y 轴上，y = {}", y),
        Point { x, y: 0 } => println!("在 x 轴上，x = {}", x),
        Point { x, y } => println!("普通点：({}, {})", x, y),
    }
}
