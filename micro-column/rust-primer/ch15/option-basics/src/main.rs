fn find_user_age(name: &str) -> Option<u32> {
    let users = [("Tony", 30), ("Alice", 25)];

    for (n, age) in users.iter() {
        if *n == name {
            return Some(*age);
        }
    }
    None
}

fn main() {
    let result1 = find_user_age("Tony");
    let result2 = find_user_age("Bob");

    println!("{:?}", result1);
    println!("{:?}", result2);

    match result1 {
        Some(age) => println!("Tony 的年龄是 {}", age),
        None => println!("找不到这个人"),
    }

    match result2 {
        Some(age) => println!("Bob 的年龄是 {}", age),
        None => println!("找不到这个人"),
    }
}
