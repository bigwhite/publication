#[derive(Debug)]
enum OrderStatus {
    PendingPayment,
    Paid,
    Shipped { tracking_number: String },
    Cancelled { reason: String },
}

struct Order {
    id: u32,
    status: OrderStatus,
}

impl Order {
    fn new(id: u32) -> Self {
        Order {
            id,
            status: OrderStatus::PendingPayment,
        }
    }

    fn describe(&self) -> String {
        match &self.status {
            OrderStatus::PendingPayment => "待支付".to_string(),
            OrderStatus::Paid => "已支付，待发货".to_string(),
            OrderStatus::Shipped { tracking_number } => {
                format!("已发货，运单号：{}", tracking_number)
            }
            OrderStatus::Cancelled { reason } => format!("已取消，原因：{}", reason),
        }
    }
}

fn main() {
    let mut order = Order::new(1001);
    println!("订单 {}：{}", order.id, order.describe());

    order.status = OrderStatus::Paid;
    println!("订单 {}：{}", order.id, order.describe());

    order.status = OrderStatus::Shipped {
        tracking_number: "SF123456".to_string(),
    };
    println!("订单 {}：{}", order.id, order.describe());

    // 现在，一个订单同一时刻只能处于这四种状态中的"恰好一种"，
    // 根本不可能出现"既发货又取消"这种自相矛盾的状态。
}
