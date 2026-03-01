// CLI: step (1-6) arg1, comma-separated field names arg2, filenames arg3+.
// Step controls how far the pipeline runs (for profiling): 1=read, 2=+parse, 3=+select, 4=+build, 5=+newline, 6=+write.

use std::collections::HashMap;
use std::env;
use std::fs::File;
use std::io::{self, BufRead, BufReader, Read, Write};
use std::sync::mpsc;
use std::sync::Arc;
use std::thread;

const PIPELINE_CAP: usize = 64;

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
    let input: Box<dyn Read + Send> = if file_name == "-" {
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

    let (read_tx, read_rx) = mpsc::sync_channel::<(usize, String)>(PIPELINE_CAP);
    let (write_tx, write_rx) = mpsc::sync_channel::<(usize, String)>(PIPELINE_CAP);
    let (err_tx, err_rx) = mpsc::sync_channel::<std::io::Error>(1);
    let err_tx = Arc::new(err_tx);

    const READ_BUF_SIZE: usize = 256 * 1024;
    let reader = BufReader::with_capacity(READ_BUF_SIZE, input);

    let err_tx_reader = Arc::clone(&err_tx);
    let h_reader = thread::spawn(move || {
        let mut reader = reader;
        let mut line = String::new();
        let mut index = 0;
        loop {
            line.clear();
            match reader.read_line(&mut line) {
                Ok(0) => break,
                Ok(_) => {}
                Err(e) => {
                    let _ = err_tx_reader.send(e);
                    break;
                }
            }
            if read_tx.send((index, line.clone())).is_err() {
                break;
            }
            index += 1;
        }
        drop(read_tx);
    });

    let inc: Vec<String> = include_fields.iter().map(|s| s.to_string()).collect();
    let h_processor = thread::spawn(move || {
        let mut mymap: HashMap<String, String> = HashMap::new();
        let mut newmap: HashMap<String, String> = HashMap::with_capacity(inc.len());
        let mut out = String::new();
        while let Ok((idx, line)) = read_rx.recv() {
            if step <= 1 {
                continue;
            }
            mymap.clear();
            for field in line.trim_end_matches('\n').split(',') {
                let mut kvps = field.splitn(2, '=');
                if let (Some(k), Some(v)) = (kvps.next(), kvps.next()) {
                    mymap.insert(k.to_string(), v.to_string());
                }
            }
            if step <= 2 {
                continue;
            }
            newmap.clear();
            for k in &inc {
                if let Some(v) = mymap.get(k) {
                    newmap.insert(k.clone(), v.clone());
                }
            }
            if step <= 3 {
                continue;
            }
            out.clear();
            let mut first = true;
            for k in &inc {
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
            if write_tx.send((idx, out.clone())).is_err() {
                break;
            }
        }
        drop(write_tx);
    });

    let err_tx_writer = Arc::clone(&err_tx);
    let h_writer = thread::spawn(move || {
        for (_, out) in write_rx {
            if let Err(e) = io::stdout().write_all(out.as_bytes()) {
                let _ = err_tx_writer.send(e);
                return;
            }
        }
    });

    let _ = h_reader.join();
    let _ = h_processor.join();
    let _ = h_writer.join();

    if let Ok(e) = err_rx.try_recv() {
        eprintln!("{}", e);
        return false;
    }
    true
}
