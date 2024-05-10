# Processing 1 Billion lines in Rust

## Performance Summary

### Tested on AMD Threadripper 1950X, 64GB of ram, 1TB WD Blue NVME. Running Fedora 40. 5/10/24
| Project | Run Time (Seconds)|
| ------- | ----------------- |
| baseline | 77.096s |
| memmapped | 76.636s |
| threaded | 5.4173s |


## Baseline

Most Straight forward approach using a BufReader to read the file in a single thread, reading until we find a new line and then processing that.
All solutions use Integers to represent fixed point decimals. This is 1. Faster, and 2. Allows us to get the correct rounding to match the output of the samples

## Memmap

Baseline but memmap the file. Last tested seems to have nearly identical performance to the baseline implementation.

## Threaded

Taking the memmap version and adding multi-threading greatly speeds up processing.
We split the file into a number of sections, based on how many threads available. 
Ensuring that each section ends with a new line.
Lastly we spawn a thread with each chunk that creates its own hashmap of the stations provided.
All threads are joined and the partial hashmaps are combined to make one final master list.

### Things to try
- Use a vector instead of hashmaps in threads
- use a different hashmap made for 10k entries (max allowed in challenge)
