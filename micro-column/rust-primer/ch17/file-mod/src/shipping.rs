pub fn calculate_fee(weight_kg: f64) -> f64 {
    weight_kg * 8.0
}

pub mod tracking;
