struct Order {
    is_paid: bool,
    is_shipped: bool,
    is_cancelled: bool,
    tracking_number: Option<String>,
    cancel_reason: Option<String>,
}

fn describe(order: &Order) -> String {
    if order.is_cancelled {
        format!("已取消，原因：{:?}", order.cancel_reason)
    } else if order.is_shipped {
        format!("已发货，运单号：{:?}", order.tracking_number)
    } else if order.is_paid {
        "已支付，待发货".to_string()
    } else {
        "待支付".to_string()
    }
}

fn main() {
    // 这是一个"非法状态"：既标记为已发货，又标记为已取消，
    // 编译器完全不会阻止你构造出这样一个自相矛盾的订单！
    let broken_order = Order {
        is_paid: true,
        is_shipped: true,
        is_cancelled: true,
        tracking_number: Some("SF123456".to_string()),
        cancel_reason: Some("用户申请退款".to_string()),
    };

    println!("{}", describe(&broken_order));
    println!("这个订单的状态自相矛盾，但类型系统完全没有拦住它！");
}
