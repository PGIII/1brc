use memchr::memchr;
use std::{
    cmp::{max, min},
    collections::HashMap,
    fmt::Display,
    fs,
    io::{BufRead, BufReader},
    str,
    sync::{Arc, Mutex},
    thread,
};

const MEASUREMENTS_FILE: &'static str = "measurements.txt";
const ZERO_ASCII: u8 = 48;
#[derive(Debug, Default, Clone)]
pub struct Station {
    pub name: Box<[u8]>,
    pub min: i32,
    pub max: i32,
    pub sum: i64,
    pub count: u32,
}

impl Display for Station {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        let min = self.min as f64 / 10.0;
        let max = self.max as f64 / 10.0;
        let mean = (self.sum / self.count as i64) as f64 / 10.0;
        let name = str::from_utf8(&self.name).unwrap();
        write!(f, "{name}={min:.1}/{mean:.1}/{max:.1}, ")
    }
}
impl Station {
    pub fn minf(&self) -> f64 {
        self.min as f64 / 10.0
    }

    pub fn new(name: &[u8], temp: i32) -> Self {
        Self {
            name: name.into(),
            min: temp,
            max: temp,
            sum: temp as i64,
            count: 1,
        }
    }
}
// takes 45sec
fn main() {
    let file = fs::OpenOptions::new()
        .read(true)
        .open(MEASUREMENTS_FILE)
        .unwrap();

    let reader = BufReader::with_capacity(4096 * 1024, file);
    let reader = Arc::new(Mutex::new(reader));
    const THREADS: i32 = 12;
    const LINES: usize = 10_000;
    let mut handles = vec![];

    for _ in 0..THREADS {
        let reader = reader.clone();
        let handle = thread::spawn(move || {
            let mut stations = HashMap::new();
            const NEW_VEC: Vec<u8> = Vec::new();
            let mut read_buf: [Vec<u8>; LINES] = [NEW_VEC; LINES];
            'reader: loop {
                {
                    let mut reader_locked = reader.lock().unwrap();
                    for line in &mut read_buf {
                        if let Ok(n) = reader_locked.read_until(b'\n', line) {
                            if n == 0 {
                                break 'reader;
                            }
                        } else {
                            break 'reader;
                        }
                    }
                }
                for line in &mut read_buf {
                    parse(&line[0..line.len() - 1], &mut stations);
                    line.clear();
                }
            }

            for line in read_buf {
                if !line.is_empty() {
                    parse(&line[0..line.len() - 1], &mut stations);
                }
            }
            stations
        });
        handles.push(handle);
    }

    let mut stations = HashMap::new();
    for handle in handles {
        let thread_stations = handle.join().unwrap();
        stations.extend(thread_stations);
    }

    print!("{{");
    for station in stations.values() {
        print!("{station}");
    }
    println!("}}");
}

fn parse(line: &[u8], stations: &mut HashMap<Box<[u8]>, Station>) {
    if line.len() == 0 {
        return;
    }
    if let Some(semi_pos) = memchr(b';', line) {
        let station_name = &line[0..semi_pos];
        let temp_str = &line[semi_pos + 1..];
        let temp_int = parse_float_str_to_int(temp_str);
        if let Some(station) = stations.get_mut(station_name) {
            station.max = max(station.max, temp_int);
            station.min = min(station.min, temp_int);
            station.sum += temp_int as i64;
            station.count += 1;
        } else {
            let station = Station::new(station_name, temp_int);
            stations.insert(station.name.clone(), station);
        }
    } else {
        panic!("Received invalid line");
    }
}

fn parse_float_str_to_int(temperature: &[u8]) -> i32 {
    let mut result: i32 = 0;
    let mut neg: i32 = 1;
    for b in temperature {
        match b {
            b'-' => {
                neg = -1;
            }
            b'0'..=b'9' => {
                let digit = b - ZERO_ASCII;
                result = result * 10 + (digit as i32);
            }
            _ => {}
        }
    }

    result * neg
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn float_parse_pos() {
        assert_eq!(parse_float_str_to_int(b"10.20"), 1020);
    }

    #[test]
    fn float_parse_neg() {
        assert_eq!(parse_float_str_to_int(b"-10.20"), -1020);
    }

    #[test]
    fn float_parse_large() {
        assert_eq!(parse_float_str_to_int(b"11110.20"), 1111020);
    }

    #[test]
    fn station_convertions() {
        let station = Station::new(b"Test", 600);
        assert_eq!(station.minf(), 60.0);
    }

    #[test]
    fn test_parse() {
        let mut stations = HashMap::new();
        let measurment_str = b"test;60.0";
        let station_name: &[u8] = b"test";

        parse(measurment_str, &mut stations);
        {
            let station = stations.get(station_name).unwrap();
            assert_eq!("test=60.0/60.0/60.0, ", format!("{station}"));
        }

        parse(b"test;30.0", &mut stations);
        {
            let station = stations.get(station_name).unwrap();
            assert_eq!("test=30.0/45.0/60.0, ", format!("{station}"));
        }
        parse(b"test;-45.0", &mut stations);
        {
            let station = stations.get(station_name).unwrap();
            assert_eq!("test=-45.0/15.0/60.0, ", format!("{station}"));
        }
    }
}
