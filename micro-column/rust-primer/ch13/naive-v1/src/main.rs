struct StrSplit<'a> {
    remainder: &'a str,
    delimiter: &'a str,
}

impl<'a> StrSplit<'a> {
    fn new(haystack: &'a str, delimiter: &'a str) -> Self {
        Self {
            remainder: haystack,
            delimiter,
        }
    }
}

impl<'a> Iterator for StrSplit<'a> {
    type Item = &'a str;

    fn next(&mut self) -> Option<Self::Item> {
        if let Some(next_delim) = self.remainder.find(self.delimiter) {
            let until_delimiter = &self.remainder[..next_delim];
            self.remainder = &self.remainder[(next_delim + self.delimiter.len())..];
            Some(until_delimiter)
        } else if self.remainder.is_empty() {
            None
        } else {
            let rest = self.remainder;
            self.remainder = "";
            Some(rest)
        }
    }
}

fn main() {
    let haystack = "a,b,";
    let parts: Vec<&str> = StrSplit::new(haystack, ",").collect();
    println!("{:?}", parts);
    println!("期望：[\"a\", \"b\", \"\"]（因为末尾还有一个逗号，理应切出一个空字符串）");
}
