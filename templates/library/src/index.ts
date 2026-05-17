export interface GreetOptions {
  /** Idioma da saudação. */
  locale?: "pt" | "en";
}

const greetings: Record<NonNullable<GreetOptions["locale"]>, string> = {
  pt: "Olá",
  en: "Hello",
};

/**
 * Retorna uma saudação personalizada.
 *
 * @example
 *   greet("mundo")                   // "Olá, mundo!"
 *   greet("world", { locale: "en" }) // "Hello, world!"
 */
export function greet(name: string, opts: GreetOptions = {}): string {
  const prefix = greetings[opts.locale ?? "pt"];
  return `${prefix}, ${name}!`;
}
