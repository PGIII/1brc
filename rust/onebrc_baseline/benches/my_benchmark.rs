use criterion::{criterion_group, criterion_main, Criterion};

fn criterion_benchmark(c: &mut Criterion) {
    c.bench_function("Baseline Calculate Averages", |b| {
        b.iter(|| onebrc_baseline::calculate_averages())
    });
}

criterion_group!(
    name = benches;
    config = Criterion::default().sample_size(10);
    targets = criterion_benchmark,
);
criterion_main!(benches);
