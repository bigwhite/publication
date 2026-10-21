fn parse_and_double_verbose(s: &str) -> Result<i32, std::num::ParseIntError> {
    let n = match s.parse::<i32>() {
        Ok(v) => v,
        Err(e) => return Err(e),
    };
    Ok(n * 2)
}

fn parse_and_double(s: &str) -> Result<i32, std::num::ParseIntError> {
    let n = s.parse::<i32>()?; // 解析失败会直接把 Err 提前返回，不需要手写 match
    Ok(n * 2)
}

fn main() {
    println!("{:?}", parse_and_double_verbose("21"));
    println!("{:?}", parse_and_double_verbose("21a"));
    println!("{:?}", parse_and_double("21"));
    println!("{:?}", parse_and_double("21a"));
}
