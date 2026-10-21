fn main() {
    let some_num: Option<i32> = Some(5);
    let none_num: Option<i32> = None;

    // map：有值就变换，没值就还是 None
    println!("map: {:?}", some_num.map(|n| n * 2));
    println!("map: {:?}", none_num.map(|n| n * 2));

    // and_then：链式调用可能失败的操作（也叫 flatMap）
    fn half_if_even(n: i32) -> Option<i32> {
        if n % 2 == 0 {
            Some(n / 2)
        } else {
            None
        }
    }
    println!("and_then: {:?}", Some(10).and_then(half_if_even));
    println!("and_then: {:?}", Some(7).and_then(half_if_even));

    // unwrap_or / unwrap_or_else：提供默认值
    println!("unwrap_or: {}", none_num.unwrap_or(0));
    println!("unwrap_or_else: {}", none_num.unwrap_or_else(|| 100));

    // ok_or：把 Option 转成 Result，方便和 ? 配合使用
    let result: Result<i32, &str> = none_num.ok_or("值不存在");
    println!("ok_or: {:?}", result);

    // Result 上也有类似的 map / unwrap_or
    let r: Result<i32, String> = Ok(21);
    println!("Result map: {:?}", r.map(|v| v * 2));
}
