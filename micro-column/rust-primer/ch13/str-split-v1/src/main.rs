struct StrSplit<'haystack, 'delimiter> {
    remainder: Option<&'haystack str>,
    delimiter: &'delimiter str,
}

impl<'haystack, 'delimiter> StrSplit<'haystack, 'delimiter> {
    fn new(haystack: &'haystack str, delimiter: &'delimiter str) -> Self {
        Self {
            remainder: Some(haystack),
            delimiter,
        }
    }
}

impl<'haystack, 'delimiter> Iterator for StrSplit<'haystack, 'delimiter> {
    type Item = &'haystack str;

    fn next(&mut self) -> Option<Self::Item> {
        let remainder = self.remainder.as_mut()?;

        if let Some(next_delim) = remainder.find(self.delimiter) {
            let until_delimiter = &remainder[..next_delim];
            *remainder = &remainder[(next_delim + self.delimiter.len())..];
            Some(until_delimiter)
        } else {
            self.remainder.take()
        }
    }
}

fn main() {
    let haystack = "a b c d e";
    let letters: Vec<&str> = StrSplit::new(haystack, " ").collect();
    println!("{:?}", letters);

    let csv = "Rust,第一课,事不过三";
    for part in StrSplit::new(csv, ",") {
        println!("part: {}", part);
    }

    let with_trailing = "a,b,";
    let parts: Vec<&str> = StrSplit::new(with_trailing, ",").collect();
    println!("{:?}", parts);
}
