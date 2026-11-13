/** Normalize all SSE line endings before the generated parser sees a chunk. */
export function normalizeEventBody(body: ReadableStream<Uint8Array>): ReadableStream<Uint8Array> {
  let pendingCR = false;
  return body.pipeThrough(
    new TransformStream<Uint8Array, Uint8Array>({
      transform(chunk, controller) {
        const out = new Uint8Array(chunk.length + 1);
        let length = 0;
        for (const byte of chunk) {
          if (pendingCR) {
            out[length++] = 10;
            pendingCR = false;
            if (byte === 10) continue;
          }
          if (byte === 13) pendingCR = true;
          else out[length++] = byte;
        }
        controller.enqueue(out.slice(0, length));
      },
      flush(controller) {
        if (pendingCR) controller.enqueue(new Uint8Array([10]));
      },
    }),
  );
}
