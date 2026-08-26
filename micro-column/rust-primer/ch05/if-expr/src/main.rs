fn main() {
    let temperature = 28;

    let feel = if temperature > 30 {
        "炎热"
    } else if temperature > 20 {
        "舒适"
    } else {
        "凉爽"
    };

    println!("当前 {} 度，体感：{}", temperature, feel);
}
