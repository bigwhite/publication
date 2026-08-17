pub fn build_greeting(name: &str) -> String {
    let base = format!("Hello, {}! 欢迎来到 Rust 第一课。", name);

    #[cfg(feature = "loud")]
    {
        base.to_uppercase()
    }

    #[cfg(not(feature = "loud"))]
    {
        base
    }
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn it_builds_greeting() {
        let msg = build_greeting("Tony");
        assert!(msg.contains("Tony"));
    }
}
