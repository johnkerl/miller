// Claude

use std::fs::File;
use std::io::{self, BufRead, BufReader};

// Method 3: Using custom buffer size for larger files
fn read_lines_custom_buffer(path: &str) -> io::Result<()> {
    const BUFFER_SIZE: usize = 128 * 1024; // 128KB buffer
    let file = File::open(path)?;
    let reader = BufReader::with_capacity(BUFFER_SIZE, file);
    
    for line in reader.lines() {
        let _line = line?;
        // Process line here
    }
    Ok(())
}

fn main() {
  let _ = read_lines_custom_buffer("/Users/kerl/tmp/big");
}
