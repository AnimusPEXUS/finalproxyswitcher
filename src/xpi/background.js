
const go = new Go();
WebAssembly.instantiateStreaming(fetch("background.wasm"), go.importObject).then((result) => {
  go.run(result.instance);
});

