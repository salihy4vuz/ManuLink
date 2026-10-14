# --- Stage 1: Chef Recipe ---
FROM docker.io/library/rust:latest AS chef
RUN cargo install cargo-chef
WORKDIR /app

# --- Stage 2: Prepare Dependencies ---
FROM chef AS planner
COPY . .
RUN cargo chef prepare --recipe-path recipe.json

# --- Stage 3: Build Application ---
FROM chef AS builder
COPY --from=planner /app/recipe.json recipe.json
# Build dependencies - this is the caching Docker layer!
RUN cargo chef cook --release --recipe-path recipe.json

# Build application
COPY . .
RUN cargo build --release --bin geometry-engine

# --- Stage 4: Runtime Environment ---
FROM gcr.io/distroless/cc-debian12

# Copy compiled binary from builder
COPY --from=builder /app/target/release/geometry-engine /geometry-engine

CMD ["/geometry-engine"]
