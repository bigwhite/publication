fn parse_config_nested(input: &str) -> Option<u32> {
    if let Some(value_str) = input.split('=').nth(1) {
        if let Ok(value) = value_str.trim().parse::<u32>() {
            Some(value)
        } else {
            println!("值不是一个合法的数字：{}", value_str);
            None
        }
    } else {
        println!("格式不对，找不到等号后面的值");
        None
    }
}

fn parse_config(input: &str) -> Option<u32> {
    let Some(value_str) = input.split('=').nth(1) else {
        println!("格式不对，找不到等号后面的值");
        return None;
    };

    let Ok(value) = value_str.trim().parse::<u32>() else {
        println!("值不是一个合法的数字：{}", value_str);
        return None;
    };

    Some(value)
}

fn main() {
    println!("{:?}", parse_config_nested("port=8080"));
    println!("{:?}", parse_config_nested("port=abc"));

    println!("{:?}", parse_config("port"));
}
