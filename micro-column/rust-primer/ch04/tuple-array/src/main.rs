fn main() {
    let point: (i32, i32, &str) = (10, 20, "起点");
    let (x, y, label) = point;
    println!("解构：x={}, y={}, label={}", x, y, label);
    println!("按索引访问：point.0={}, point.1={}", point.0, point.1);

    let scores: [i32; 5] = [90, 85, 77, 88, 95];
    println!("数组长度：{}", scores.len());
    println!("第一个元素：{}", scores[0]);

    let zeros = [0; 3];
    println!("用简写语法初始化：{:?}", zeros);

    let total: i32 = scores.iter().sum();
    println!("总分：{}", total);
}
