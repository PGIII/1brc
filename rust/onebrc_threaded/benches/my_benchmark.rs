use criterion::{criterion_group, criterion_main, Criterion};

fn criterion_benchmark(c: &mut Criterion) {
    c.bench_function("Threaded Calculate Averages", |b| {
        b.iter(|| onebrc_threaded::calculate_averages())
    });
}

criterion_group!(
    name = benches;
    config = Criterion::default().sample_size(10);
    targets = criterion_benchmark,
);
criterion_main!(benches);
