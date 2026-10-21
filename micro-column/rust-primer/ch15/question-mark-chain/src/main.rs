fn parse_and_sum(a: &str, b: &str) -> Result<i32, std::num::ParseIntError> {
    let x = a.parse::<i32>()?;
    let y = b.parse::<i32>()?;
    Ok(x + y)
}

fn main() {
    println!("{:?}", parse_and_sum("10", "20"));
    println!("{:?}", parse_and_sum("10", "abc"));
}
