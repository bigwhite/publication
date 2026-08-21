fn main() {
    let mut count = 0;
    count += 1;
    count += 1;
    println!("count（mut，同一类型累加）: {}", count);

    let spaces = "   ";
    let spaces = spaces.len();
    println!("spaces（shadowing，类型从 &str 变成 usize）: {}", spaces);

    let x = 5;
    {
        let x = x * 2;
        println!("内层作用域里的 x: {}", x);
    }
    println!("外层作用域里的 x 依然是: {}", x);
}
