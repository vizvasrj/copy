use std::fs;
use std::io;
use std::path::Path;
use std::sync::{Arc, Mutex};
use std::thread;
use std::time::Instant;

fn copy_item(source: &Path, destination: &Path) -> io::Result<()> {
    if source.is_file() {
        fs::copy(source, destination)?;
    } else if source.is_dir() {
        if !destination.exists() {
            fs::create_dir(destination)?;
        }
        let entries = fs::read_dir(source)?;
        for entry in entries {
            let entry = entry?;
            let entry_source = entry.path();
            let entry_destination = destination.join(entry.file_name());
            copy_item(&entry_source, &entry_destination)?;
        }
    }
    Ok(())
}

fn main() {
    let start_time = Instant::now();
    // Configuration
    let source_path = Path::new("demo_files4/");
    let destination_path = Path::new("/home/tmp");



    // Create destination directory if it doesn't exist
    if !destination_path.exists() {
        fs::create_dir_all(destination_path).unwrap();
    }

    // Get list of files and directories in source path
    let items_res = fs::read_dir(source_path);
    if items_res.is_err() {
        println!("errors at {} {:?}", items_res.err().unwrap(), source_path);
        return;
    }
    let items = items_res
        .unwrap()
        .map(|entry| entry.unwrap().path())
        .collect::<Vec<_>>();

    // Calculate the number of CPU cores
    let num_cores = 4;

    // Create a counter to track the number of completed copy operations
    let counter = Arc::new(Mutex::new(0));

    // Spawn a thread for each CPU core
    let threads: Vec<_> = (0..num_cores)
        .map(|_| {
            let counter = Arc::clone(&counter);
            let items = items.clone();
            thread::spawn(move || {
                loop {
                    let item = {
                        let mut counter = counter.lock().unwrap();
                        let item = items.get(*counter);
                        *counter += 1;
                        item.cloned()
                    };

                    match item {
                        Some(item) => {
                            let destination_item = destination_path.join(
                                item.strip_prefix(&source_path).unwrap(),
                            );
                            copy_item(&item, &destination_item).unwrap();
                        }
                        None => break,
                    }
                }
            })
        })
        .collect();

    // Wait for all threads to finish
    for thread in threads {
        thread.join().unwrap();
    }
    let elapsed_time = start_time.elapsed();
    
    println!("All items copied. took {:?} seconds", elapsed_time);
}