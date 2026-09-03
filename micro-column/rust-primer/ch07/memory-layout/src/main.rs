fn main() {
    let s1 = String::from("hello");

    println!("s1 在栈上的三个字段：");
    println!("  指针地址（ptr）  : {:p}", s1.as_ptr());
    println!("  长度（len）      : {}", s1.len());
    println!("  容量（capacity） : {}", s1.capacity());
    println!(
        "  String 本身在栈上占用的字节数：{}",
        std::mem::size_of::<String>()
    );

    let s2 = s1;

    println!();
    println!("s2 = s1（发生 move）之后：");
    println!("  s2 指针地址（ptr）: {:p}", s2.as_ptr());
    println!("  说明：s2 的 ptr 和刚才 s1 的 ptr 完全相同，");
    println!("  因为只是把栈上那三个字段拷贝了一份给 s2，堆上的数据没有被复制。");
}
