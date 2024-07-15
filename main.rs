extern crate rodio;
extern crate aubio;
extern crate minifb;

use std::fs::File;
use std::io::BufReader;
use rodio::{Decoder, Source};
use aubio::{Pitch, PitchMode};
use minifb::{Key, Window, WindowOptions};

const WIDTH: usize = 800;
const HEIGHT: usize = 600;

fn main() {
    let file = File::open("path_to_audio_file.mp3").unwrap();
    let source = Decoder::new(BufReader::new(file)).unwrap();

    let mut pitch = Pitch::new(PitchMode::Yin, 2048, 1024, source.sample_rate()).unwrap();

    let mut buffer: Vec<u32> = vec![0; WIDTH * HEIGHT];
    let mut window = Window::new(
        "Synesthesia Art Generator - Rust",
        WIDTH,
        HEIGHT,
        WindowOptions::default(),
    )
    .unwrap_or_else(|e| {
        panic!("Window creation failed: {}", e);
    });

    for frame in source.into_iter() {
        let frame = frame.unwrap();
        let freq = pitch.do_(&frame);

        let color = map_frequency_to_color(freq);

        for pixel in buffer.iter_mut() {
            *pixel = color;
        }

        window
            .update_with_buffer(&buffer, WIDTH, HEIGHT)
            .unwrap_or_else(|e| {
                println!("Window update failed: {}", e);
            });

        if !window.is_open() {
            break;
        }
        window.get_keys_pressed(Key::Escape);
    }
}

fn map_frequency_to_color(freq: f32) -> u32 {
    let gray_value = (freq / 20000.0 * 255.0) as u32;
    (gray_value << 16) | (gray_value << 8) | gray_value
}
