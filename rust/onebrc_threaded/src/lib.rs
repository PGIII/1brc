use memchr::{memchr, memchr_iter, memrchr};
use memmap2::MmapOptions;
use std::{
    cmp::{max, min},
    collections::HashMap,
    fmt::Display,
    fs, io, str,
    thread::{self, available_parallelism},
};

const MEASUREMENTS_FILE: &'static str = "measurements.txt";
const ZERO_ASCII: u8 = 48;

#[derive(Debug, Default, Clone)]
pub struct Station<'a> {
    pub name: &'a [u8],
    pub min: i32,
    pub max: i32,
    pub sum: i64,
    pub count: u32,
}

impl Display for Station<'_> {
    fn fmt(&self, f: &mut std::fmt::Formatter<'_>) -> std::fmt::Result {
        let min = self.min as f64 / 100.0;
        let max = self.max as f64 / 100.0;
        let mean = round_to_1decimal_f64(self.sum / self.count as i64);
        let name = str::from_utf8(&self.name).unwrap();
        write!(f, "{name}={min:.1}/{mean:.1}/{max:.1}")
    }
}

impl<'a> Station<'a> {
    pub fn new(name: &'a [u8], temp: i32) -> Self {
        Self {
            name,
            min: temp,
            max: temp,
            sum: temp as i64,
            count: 1,
        }
    }
}

type StationMap<'a> = HashMap<&'a [u8], Station<'a>>;

#[inline]
pub fn calculate_averages() {
    let file = fs::OpenOptions::new()
        .read(true)
        .open(MEASUREMENTS_FILE)
        .unwrap();
    let file_mm = unsafe { MmapOptions::new().map(&file).unwrap() };
    let file_size = file_mm.len();
    let mut stations = StationMap::new();
    let threads = available_parallelism().unwrap();
    let chunk_size = file_size / threads;
    let mut last_end = 0;

    thread::scope(|s| {
        let mut handles = vec![];
        while last_end < file_size {
            //split up file, but make sure last byte of chunk is \n
            let potential_end = min(last_end + chunk_size, file_size);
            let end = if potential_end != file_size {
                min(
                    memchr(b'\n', &file_mm[last_end + chunk_size..]).unwrap()
                        + 1
                        + chunk_size
                        + last_end,
                    file_size,
                )
            } else {
                file_size
            };
            let chunk = &file_mm[last_end..end];
            last_end = end;
            let handle = s.spawn(move || {
                let mut stations = HashMap::new();
                let iter = memchr_iter(b'\n', &chunk);
                let mut last = 0;
                for i in iter {
                    parse(&chunk[last..i], &mut stations);
                    last = i + 1;
                }
                stations
            });
            handles.push(handle);
        }

        let mut partial_maps = vec![];
        for handle in handles {
            let map = handle.join().unwrap();
            partial_maps.push(map);
        }
        for map in partial_maps {
            for (name, new_station) in map {
                //check if exists in map, if so merge them, other wise insert it
                if let Some(station) = stations.get_mut(name) {
                    station.min = min(new_station.min, station.min);
                    station.max = max(new_station.max, station.max);
                    station.count += new_station.count;
                    station.sum += new_station.sum;
                } else {
                    stations.insert(name, new_station);
                }
            }
        }
    });
    let mut output = io::stdout().lock();
    print_string(&stations, &mut output);
}

/// Convert out fixed 2 decimal type to a f64 with 1 decimal place thats properly rounded
fn round_to_1decimal_f64(num: i64) -> f64 {
    let hund = num % 10;
    let mut num = num / 10;
    if hund >= 5 {
        num += 1;
    }
    num as f64 / 10.0
}

fn print_string(stations: &StationMap, output: &mut impl io::Write) {
    let mut sorted: Vec<_> = stations.iter().collect();
    sorted.sort_by_key(|&(name, _)| name);
    let mut sorted = sorted.into_iter();
    let (_, first) = sorted.next().unwrap();
    write!(output, "{{").unwrap();
    write!(output, "{first}").unwrap();
    for (_, station) in sorted {
        write!(output, ", {station}").unwrap();
    }
    writeln!(output, "}}").unwrap();
}

fn parse<'a>(line: &'a [u8], stations: &mut HashMap<&'a [u8], Station<'a>>) {
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
            stations.insert(station.name.into(), station);
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

    result * neg * 10
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn float_parse_pos() {
        assert_eq!(parse_float_str_to_int(b"10.20"), 10200);
    }

    #[test]
    fn float_parse_neg() {
        assert_eq!(parse_float_str_to_int(b"-10.20"), -10200);
    }

    #[test]
    fn float_parse_large() {
        assert_eq!(parse_float_str_to_int(b"11110.20"), 11110200);
    }

    #[test]
    fn station_convertions() {
        let station = Station::new(b"Test", 6000);
        assert_eq!(station.min as f64 / 100.0, 60.0);
    }

    #[test]
    fn test_parse() {
        let mut stations = HashMap::new();
        let measurment_str = b"test;60.0";
        let station_name: &[u8] = b"test";

        parse(measurment_str, &mut stations);
        {
            let station = stations.get(station_name).unwrap();
            assert_eq!("test=60.0/60.0/60.0", format!("{station}"));
        }

        parse(b"test;30.0", &mut stations);
        {
            let station = stations.get(station_name).unwrap();
            assert_eq!("test=30.0/45.0/60.0", format!("{station}"));
        }
        parse(b"test;-45.0", &mut stations);
        {
            let station = stations.get(station_name).unwrap();
            assert_eq!("test=-45.0/15.0/60.0", format!("{station}"));
        }
    }
}
