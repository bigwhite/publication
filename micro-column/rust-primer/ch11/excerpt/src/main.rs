struct Excerpt<'a> {
    part: &'a str,
}

impl<'a> Excerpt<'a> {
    // 这个方法的签名没有显式写生命周期，但它其实等价于：
    // fn announce<'b, 'c>(
    //     &'b self,
    //     prefix: &'c str,
    // ) -> &'b str {
    //     println!("提示：{}", prefix);
    //     self.part
    // }
    // 编译器靠"省略规则"自动帮我们补全了，我们会在下一讲把这条规则讲透。
    fn announce(&self, prefix: &str) -> &str {
        println!("提示：{}", prefix);
        self.part
    }

    fn part_len(&self) -> usize {
        self.part.len()
    }
}

fn main() {
    let novel = String::from("事不过三。这一次，我们扬帆起航，绝不返航。");
    let first_sentence = novel.split('。').next().expect("找不到句号");

    let excerpt = Excerpt {
        part: first_sentence,
    };

    println!("摘录：{}", excerpt.part);
    println!("摘录长度：{}", excerpt.part_len());

    let announced = excerpt.announce("重要摘录如下");
    println!("宣布结果：{}", announced);
}
