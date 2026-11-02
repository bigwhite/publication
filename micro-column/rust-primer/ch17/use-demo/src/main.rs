mod shipping {
    pub mod tracking {
        pub fn generate_id() -> String {
            String::from("SF1234567890")
        }
    }

    // 重导出：把内部深层路径的功能，"提升"到 shipping 模块自己这一层直接暴露出去
    pub use self::tracking::generate_id;
}

fn main() {
    // 不用 use，每次都写完整路径，比较啰嗦
    let id1 = shipping::tracking::generate_id();

    // 用 use 引入路径之后，可以直接用短名字
    use shipping::tracking::generate_id as gen_id_full_path;
    let id2 = gen_id_full_path();

    // 因为 shipping 模块里做了 pub use 重导出，调用者甚至不需要知道 tracking 这一层的存在
    let id3 = shipping::generate_id();

    println!("{} / {} / {}", id1, id2, id3);
}
