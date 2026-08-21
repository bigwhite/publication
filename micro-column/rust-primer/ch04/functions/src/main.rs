// 这是一行普通注释，编译器会完全忽略它，只是给人看的。

/// 计算圆的面积。
///
/// # 参数
/// * `radius` - 圆的半径
///
/// # 返回值
/// 圆的面积（f64）
fn circle_area(radius: f64) -> f64 {
    std::f64::consts::PI * radius * radius
}

fn add(a: i32, b: i32) -> i32 {
    a + b // 注意：没有分号，这是函数体的最后一个表达式，也就是返回值
}

fn add_with_return(a: i32, b: i32) -> i32 {
    return a + b; // 显式 return 也可以，但要加分号，风格上通常只在提前返回时使用
}

fn main() {
    let area = circle_area(2.0);
    println!("半径为 2 的圆面积：{:.2}", area);

    println!("add(1, 2) = {}", add(1, 2));
    println!("add_with_return(1, 2) = {}", add_with_return(1, 2));
}
