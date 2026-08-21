fn main() {
    let x = 5;

    let y = {
        let inner = x * 2;
        inner + 1
    };

    println!("x = {}", x);
    println!("y = {}", y);

    let z = if y > 10 { "大于10" } else { "不大于10" };
    println!("z = {}", z);
}
