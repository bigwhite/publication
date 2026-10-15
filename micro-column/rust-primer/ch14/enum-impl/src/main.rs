enum OrderStatus {
    PendingPayment,
    Paid,
    Shipped { tracking_number: String },
    Cancelled { reason: String },
}

impl OrderStatus {
    // 直接给 Enum 本身定义方法，判断是否已经是"终态"
    fn is_final(&self) -> bool {
        matches!(self, OrderStatus::Cancelled { .. })
    }

    // 关联函数：不依赖具体实例，常用作"构造器"
    fn new_cancelled(reason: &str) -> Self {
        OrderStatus::Cancelled {
            reason: reason.to_string(),
        }
    }
}

fn main() {
    let s1 = OrderStatus::Paid;
    let s2 = OrderStatus::new_cancelled("库存不足");

    println!("s1 是终态吗？{}", s1.is_final());
    println!("s2 是终态吗？{}", s2.is_final());
}
