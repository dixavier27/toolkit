import pc from "picocolors";
import { getLogLevel } from "./logger.ts";

const FRAMES = ["⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"];

export interface Spinner {
  start: () => void;
  stop: (finalMessage?: string) => void;
  update: (message: string) => void;
}

export function spinner(initialMessage: string): Spinner {
  const isTTY = process.stdout.isTTY === true;
  const enabled = isTTY && getLogLevel() !== "silent";

  let message = initialMessage;
  let frameIndex = 0;
  let interval: ReturnType<typeof setInterval> | undefined;

  function render() {
    const frame = FRAMES[frameIndex % FRAMES.length];
    process.stdout.write(`\r${pc.cyan(frame)} ${message}`);
    frameIndex++;
  }

  function clear() {
    if (!isTTY) return;
    process.stdout.write(`\r${" ".repeat(message.length + 4)}\r`);
  }

  return {
    start() {
      if (!enabled) {
        console.log(message);
        return;
      }
      render();
      interval = setInterval(render, 80);
    },
    update(newMessage: string) {
      message = newMessage;
      if (!enabled) {
        console.log(message);
        return;
      }
      clear();
      render();
    },
    stop(finalMessage?: string) {
      if (interval) {
        clearInterval(interval);
        interval = undefined;
      }
      clear();
      if (finalMessage && getLogLevel() !== "silent") {
        console.log(finalMessage);
      }
    },
  };
}
