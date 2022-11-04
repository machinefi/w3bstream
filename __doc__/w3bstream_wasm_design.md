## Wasm programming model for w3bStream

### Rational

Several programming models is proposed for w3bStream

1. Data In Data Out

```rust
fn main(data_ptr *i32, data_size i32) -> i32 {
    copy(scr, data_ptr, output_size);
    return output_size
}
```

In this model, wasm acts as a streaming data handler. It is straightforward to pass a input data in and return the filtered data from the func.

2. Struct Wrapper

```rust
fn main(in Request) -> Response {
}
```

Examples:

[wasmedge_wasi_socket](https://github.com/second-state/wasmedge_wasi_socket/blob/main/_examples/http_server/src/main.rs#L15)

[Rust on Compute@Edge](https://developer.fastly.com/learning/compute/rust/#main-interface)

This model enforces the adaptation of pre-defined struct on the bost host side and wasm dev side, which doesn't provide any advantages compared with the data-in data-out model in data handling scenario.

3. Main-entry style (preferred one)

```rust
fn main(resource_id i32) -> i32 {
    return status_code;
}
```

In this model, the data from IoT devices is not directly passed to the wasm once it is initiated. But it has to be requested data with the resource_id from the host. This model is preferred because it transfer to the [memory ownership](https://github.com/proxy-wasm/spec/tree/master/abi-versions#memory-ownership) to wasm code, who can manage the memory by itself. Besides, a bundle of ABIs is used to build the bridge between host and wasm.

ABI examples:

[proxy-wasm-go-host/imports.go](https://github.com/mosn/proxy-wasm-go-host/blob/main/proxywasm/v2/imports.go)

[proxy-wasm-go-host/exports.go](https://github.com/mosn/proxy-wasm-go-host/blob/main/proxywasm/v2/exports.go)

4. Injection

```rust
fn map() -> i32 {
    return status_code; 
}

fn reduce() -> i32 {
    return status_code;
}
```

In this model, the developer is supposed to inject the wasm into the interface the host provides. When certain func of the interface is invoked on the host side, the corresponding wasm is run. The development of this model is easier than other models, but, unlike HTTP handling, the data flow for streaming isn't fixed. So this model isn't suitable for our application.

Examples:

[proxy-wasm](https://github.com/proxy-wasm/spec/tree/master/abi-versions/vNEXT)

### Implements

#### ABIs

- Func exported by wasm

```rust
fn alloc(size: usize) -> *mut c_void {
    return ptr;
}

fn start(resource_id i32) -> i32 {
    return status_code;
}
```

- Func imported to wasm

```rust
fn ws_get_data(resource_id i32, return_ptr i32, return_size i32) -> i32 {
    copy(data_ptr, return_ptr, return_size)
    return Result_OK
}

fn ws_set_data(resource_id i32, ptr i32, size i32) -> i32 {
    return Result_OK
}

fn ws_get_dB(key_ptr i32, key_size i32, return_value_ptr i32, return_value_size i32) -> i32 {
    return Result_OK
}

fn ws_set_dB(key_ptr i32, key_size i32, value_ptr i32, value_size i32) -> i32 {
    return Result_OK
}

fn ws_log(logLevel i32, ptr i32, size i32) -> i32 {
    return Result_OK
}

// an example of encoded data
// `{
//        "to":    "0xb576c141e5659137ddda4223d209d4744b2106be",
//        "value": "0",
//        "data":  "..."  // hex encoding 
// }`
fn ws_send_tx(encoded_data i32, encoded_size i32, return_value_ptr i32, return_value_size i32) -> i32 {}

// an example of encoded data
// `{
//        "to":    "0xb576c141e5659137ddda4223d209d4744b2106be",
//        "data":  "..."  // hex encoding 
// }`
fn ws_call_contract(encoded_ptr i32, encoded_size i32) -> i32 {}
```

