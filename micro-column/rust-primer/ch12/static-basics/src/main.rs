// static 变量：整个程序运行期间只有一份，存放在固定的内存地址
static GREETING: &str = "事不过三，我们一起学 Rust";

fn get_greeting() -> &'static str {
    GREETING
}

fn get_literal() -> &'static str {
    // 字符串字面量本身的类型就是 &'static str，直接返回完全没问题
    "这也是一个 'static 的字符串"
}

fn main() {
    println!("{}", get_greeting());
    println!("{}", get_literal());

    let s: &'static str = "显式标注 'static 的字符串字面量";
    println!("{}", s);
}
