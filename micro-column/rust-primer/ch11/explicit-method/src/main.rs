struct Excerpt<'a> {
    part: &'a str,
}

impl<'a> Excerpt<'a> {
    // 完整显式写出的版本，和省略版本完全等价
    fn announce_explicit<'b, 'c>(&'b self, prefix: &'c str) -> &'b str {
        println!("提示：{}", prefix);
        self.part
    }
}

fn main() {
    let novel = String::from("事不过三。这一次，我们扬帆起航，绝不返航。");
    let first_sentence = novel.split('。').next().expect("找不到句号");
    let excerpt = Excerpt {
        part: first_sentence,
    };

    let prefix = String::from("重要摘录如下");
    let announced = excerpt.announce_explicit(prefix.as_str());
    println!("宣布结果：{}", announced);
}
