extern crate taubyte_sdk;
use taubyte_sdk::event::Event;

#[no_mangle]
pub fn ping(event: Event) -> u32 {
    let http = event.http().unwrap();
    _ = http.write("pong".as_bytes()).unwrap();

    0
}
