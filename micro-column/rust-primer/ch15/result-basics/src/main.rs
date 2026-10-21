fn divide(a: f64, b: f64) -> Result<f64, String> {
    if b == 0.0 {
        Err(String::from("除数不能为 0"))
    } else {
        Ok(a / b)
    }
}

fn main() {
    let r1 = divide(10.0, 2.0);
    let r2 = divide(10.0, 0.0);

    match r1 {
        Ok(v) => println!("10 / 2 = {}", v),
        Err(e) => println!("出错了：{}", e),
    }

    match r2 {
        Ok(v) => println!("10 / 0 = {}", v),
        Err(e) => println!("出错了：{}", e),
    }
}
