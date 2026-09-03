fn main() {
    let value = 42;

    // 引用：编译器全程保驾护航，&value 保证一定指向一个有效的 i32
    let reference: &i32 = &value;
    println!("通过引用读取：{}", *reference);

    // 裸指针：只是记录了一个地址，编译器不再对它做任何有效性检查
    let raw_ptr: *const i32 = &value as *const i32;
    println!("裸指针存的地址：{:p}", raw_ptr);

    // 解引用裸指针必须显式包在 unsafe 块里——
    // 这是 Rust 在提醒你："从这里开始，安全保证由你自己负责"
    unsafe {
        println!("通过裸指针读取：{}", *raw_ptr);
    }
}
