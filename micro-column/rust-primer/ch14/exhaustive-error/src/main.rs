enum OrderStatus {
    PendingPayment,
    Paid,
    Shipped { tracking_number: String },
    Cancelled { reason: String },
}

fn describe(status: &OrderStatus) -> String {
    match status {
        OrderStatus::PendingPayment => "待支付".to_string(),
        OrderStatus::Paid => "已支付，待发货".to_string(),
        OrderStatus::Shipped { tracking_number } => {
            format!("已发货，运单号：{}", tracking_number)
        } // 故意漏掉了 Cancelled 分支
    }
}

fn main() {}
