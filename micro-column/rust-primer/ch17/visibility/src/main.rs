mod billing {
    pub struct Invoice {
        pub amount: f64,
        tax_id: String, // 默认私有，外部无法访问
    }

    impl Invoice {
        pub fn new(amount: f64) -> Self {
            Invoice {
                amount,
                tax_id: String::from("TAX-001"),
            }
        }

        // pub(crate)：只在当前 crate 内部可见，crate 外部（比如未来发布成库给别人用）无法调用
        pub(crate) fn internal_tax_id(&self) -> &str {
            &self.tax_id
        }
    }
}

fn main() {
    let invoice = billing::Invoice::new(100.0);
    //println!("{}", invoice.tax_id);
    println!("{}", invoice.amount);

    // 因为 main 和 billing 模块在同一个 crate 内，pub(crate) 的方法可以被调用
    println!("税号（crate 内部可见）：{}", invoice.internal_tax_id());
}
