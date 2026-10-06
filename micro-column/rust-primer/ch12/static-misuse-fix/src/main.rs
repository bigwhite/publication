fn make_greeting(name: &str) -> &'static str {
    let s = format!("你好，{}！", name);
    Box::leak(s.into_boxed_str())
}

fn main() {
    let greeting = make_greeting("Tony");
    println!("{}", greeting);
    println!("这份内存永远不会被回收了，这就是 Box::leak 的代价");
}
