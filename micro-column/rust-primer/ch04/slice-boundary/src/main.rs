fn main() {
    let s = String::from("Rust中文");
    let slice = &s[0..4];
    println!("正常切片（前4字节，正好是 Rust）：{}", slice);

    let bad_slice = &s[0..5]; // 5 落在"中"这个字的内部字节边界上
    println!("{}", bad_slice);
}
