fn main() {
    // 裸指针允许是空指针，编译器完全不会阻止你创建它
    // （只要不去解引用它，创建本身是安全的、能正常编译运行的）
    let null_ptr: *const i32 = std::ptr::null();
    println!("空裸指针：{:p}", null_ptr);
    println!("空指针本身合法存在，但绝不能对它解引用，那样会是未定义行为");
}
