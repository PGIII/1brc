# Processing 1 Billion lines in Rust

## Baseline

Most Straight forward approach using a BufReader to read the file in a single thread, reading until we find a new line and then processing that.
All solutions use Integers to represent fixed point decimals. This is 1. Faster, and 2. Allows us to get the correct rounding to match the output of the samples

## Memmap

Baseline but with memmaped file instead of buf reader. This is actually a slower implementation due to no buffering i assume.

## Threaded

Taking the memmap version and adding multi-threading greatly speeds up processing, running in 5.6sec on a 1950x with 64GB or ram.
We split the file into sections based on threads available.
Then make sure that each section ends with a new line.
Lastly we spawn a thread with this chunk that creates its own hashmap and returns that for combining at the end.

### Things to try
- Use a vector instead of hashmaps in threads
- use a different hashmap made for 10k entries (max allowed in challenge)
