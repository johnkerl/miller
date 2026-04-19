// CLI: step (1-6) arg1, comma-separated field names arg2, filenames arg3+.
// Step controls how far the pipeline runs (for profiling): 1=read, 2=+parse, 3=+select, 4=+build, 5=+newline, 6=+write.

use std::collections::HashMap;
use std::env;
use std::fs::File;
use std::io::{self, BufRead, BufReader, Read, Write};

fn main() {
    let args: Vec<String> = env::args().collect();
    if args.len() < 3 {
        eprintln!("usage: {} <step 1-6> <field1,field2,...> [file ...]", args[0]);
        std::process::exit(1);
    }
    let step: u32 = match args[1].parse() {
        Ok(n) if (1..=6).contains(&n) => n,
        _ => {
            eprintln!("step must be 1-6, got \"{}\"", args[1]);
            std::process::exit(1);
        }
    };
    let include_fields: Vec<&str> = args[2].split(',').collect();
    let filenames: Vec<&str> = if args.len() > 3 {
        args[3..].iter().map(String::as_str).collect()
    } else {
        vec!["-"]
    };

    let mut ok = true;
    for name in filenames {
        ok = handle(name, step, &include_fields) && ok;
    }
    std::process::exit(if ok { 0 } else { 1 });
}

fn handle(file_name: &str, step: u32, include_fields: &[&str]) -> bool {
    let input: Box<dyn Read> = if file_name == "-" {
        Box::new(io::stdin())
    } else {
        match File::open(file_name) {
            Ok(f) => Box::new(f),
            Err(e) => {
                eprintln!("{}", e);
                return false;
            }
        }
    };

    let mut reader = BufReader::new(input);
    let mut line = String::new();

    loop {
        line.clear();
        let n = match reader.read_line(&mut line) {
            Ok(0) => break,
            Ok(n) => n,
            Err(e) => {
                eprintln!("{}", e);
                return false;
            }
        };
        if n > 0 && !line.ends_with('\n') {
            // EOF without newline; line still has content
        }

        if step <= 1 {
            continue;
        }

        // Step 2: line to map
        let mymap: HashMap<String, String> = line
            .trim_end_matches('\n')
            .split(',')
            .filter_map(|field| {
                let mut kvps = field.splitn(2, '=');
                let k = kvps.next()?;
                let v = kvps.next()?;
                Some((k.to_string(), v.to_string()))
            })
            .collect();
        if step <= 2 {
            continue;
        }

        // Step 3: map-to-map transform (newmap for step 4 lookup; order from include_fields)
        let newmap: HashMap<String, String> = include_fields
            .iter()
            .filter_map(|&k| mymap.get(k).map(|v| (k.to_string(), v.clone())))
            .collect();
        if step <= 3 {
            continue;
        }

        // Step 4–5: map to string + newline (iterate include_fields to preserve order)
        let mut first = true;
        let mut out = String::new();
        for &k in include_fields {
            if let Some(v) = newmap.get(k) {
                if !first {
                    out.push(',');
                }
                out.push_str(k);
                out.push('=');
                out.push_str(v);
                first = false;
            }
        }
        out.push('\n');
        if step <= 5 {
            continue;
        }

        // Step 6: write to stdout
        if let Err(e) = io::stdout().write_all(out.as_bytes()) {
            eprintln!("{}", e);
            return false;
        }
    }

    true
}
