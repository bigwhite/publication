use std::fs;

pub struct Config {
    pub pattern: String,
    pub filename: String,
    pub case_insensitive: bool,
    pub show_line_numbers: bool,
}

impl Config {
    pub fn build(args: &[String]) -> Result<Config, String> {
        // args[0] 是程序自身路径，真正的参数从 args[1] 开始
        let mut positional = Vec::new();
        let mut case_insensitive = false;
        let mut show_line_numbers = false;

        for arg in &args[1..] {
            match arg.as_str() {
                "-i" => case_insensitive = true,
                "-n" => show_line_numbers = true,
                other => positional.push(other.to_string()),
            }
        }

        if positional.len() < 2 {
            return Err(format!(
                "用法: cligrep [-i] [-n] <pattern> <filename>\n收到的参数不足：{:?}",
                positional
            ));
        }

        Ok(Config {
            pattern: positional[0].clone(),
            filename: positional[1].clone(),
            case_insensitive,
            show_line_numbers,
        })
    }
}

pub fn search_case_sensitive<'a>(pattern: &str, contents: &'a str) -> Vec<(usize, &'a str)> {
    let mut results = Vec::new();
    for (i, line) in contents.lines().enumerate() {
        if line.contains(pattern) {
            results.push((i + 1, line));
        }
    }
    results
}

pub fn search_case_insensitive<'a>(pattern: &str, contents: &'a str) -> Vec<(usize, &'a str)> {
    let pattern = pattern.to_lowercase();
    let mut results = Vec::new();
    for (i, line) in contents.lines().enumerate() {
        if line.to_lowercase().contains(&pattern) {
            results.push((i + 1, line));
        }
    }
    results
}

pub fn run(config: &Config) -> Result<(), String> {
    let contents = fs::read_to_string(&config.filename)
        .map_err(|e| format!("无法读取文件 '{}'：{}", config.filename, e))?;

    let results = if config.case_insensitive {
        search_case_insensitive(&config.pattern, &contents)
    } else {
        search_case_sensitive(&config.pattern, &contents)
    };

    if results.is_empty() {
        println!("没有找到匹配 \"{}\" 的行", config.pattern);
    } else {
        for (line_no, line) in results {
            if config.show_line_numbers {
                println!("{}: {}", line_no, line);
            } else {
                println!("{}", line);
            }
        }
    }

    Ok(())
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn case_sensitive_finds_exact_match() {
        let contents = "\
Rust第一课
事不过三
RUST is great";
        assert_eq!(
            search_case_sensitive("Rust", contents),
            vec![(1, "Rust第一课")]
        );
    }

    #[test]
    fn case_insensitive_finds_all_case_variants() {
        let contents = "\
Rust第一课
事不过三
RUST is great";
        assert_eq!(
            search_case_insensitive("rust", contents),
            vec![(1, "Rust第一课"), (3, "RUST is great")]
        );
    }

    #[test]
    fn build_rejects_insufficient_args() {
        let args = vec!["cligrep".to_string(), "onlyone".to_string()];
        assert!(Config::build(&args).is_err());
    }

    #[test]
    fn build_parses_flags() {
        let args = vec![
            "cligrep".to_string(),
            "-i".to_string(),
            "-n".to_string(),
            "rust".to_string(),
            "poem.txt".to_string(),
        ];
        let config = Config::build(&args).unwrap();
        assert!(config.case_insensitive);
        assert!(config.show_line_numbers);
        assert_eq!(config.pattern, "rust");
        assert_eq!(config.filename, "poem.txt");
    }
}
