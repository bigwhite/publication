mod shipping {
    pub fn calculate_fee(weight_kg: f64) -> f64 {
        weight_kg * 8.0
    }

    pub mod tracking {
        pub fn generate_id() -> String {
            String::from("SF1234567890")
        }
    }
}

mod billing {
    pub fn generate_invoice(amount: f64) -> String {
        format!("发票金额：{:.2} 元", amount)
    }
}

fn main() {
    let fee = shipping::calculate_fee(3.5);
    println!("运费：{:.2} 元", fee);

    let id = shipping::tracking::generate_id();
    println!("运单号：{}", id);

    let invoice = billing::generate_invoice(fee);
    println!("{}", invoice);
}
