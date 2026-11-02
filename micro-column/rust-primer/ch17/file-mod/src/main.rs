mod billing;
mod shipping;

fn main() {
    let fee = shipping::calculate_fee(3.5);
    println!("运费：{:.2} 元", fee);

    let id = shipping::tracking::generate_id();
    println!("运单号：{}", id);

    let invoice = billing::generate_invoice(fee);
    println!("{}", invoice);
}
