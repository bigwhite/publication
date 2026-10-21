fn divide(a: f64, b: f64) -> Result<f64, String> {
    if b == 0.0 {
        Err(String::from("除数不能为 0"))
    } else {
        Ok(a / b)
    }
}

fn main() {
    let result = divide(10.0, 0.0).unwrap();
    println!("{}", result);
}
