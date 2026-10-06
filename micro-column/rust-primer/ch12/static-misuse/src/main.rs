fn make_greeting(name: &str) -> &'static str {
    let s = format!("你好，{}！", name);
    s.as_str()
}

fn main() {
    let greeting = make_greeting("Tony");
    println!("{}", greeting);
}
