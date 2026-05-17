import { Component, signal } from "@angular/core";

@Component({
  selector: "app-root",
  standalone: true,
  template: `
    <main class="min-h-screen flex items-center justify-center p-8">
      <section class="max-w-md w-full bg-white rounded-2xl shadow-sm border border-slate-200 p-8">
        <h1 class="text-2xl font-bold mb-2">🌱 {{ title }}</h1>
        <p class="text-slate-600 mb-6">
          App Angular + Tauri scaffolded com
          <code class="text-sm bg-slate-100 px-1.5 py-0.5 rounded">eco</code>.
        </p>

        <button
          type="button"
          (click)="increment()"
          class="px-4 py-2 bg-slate-900 text-white rounded-lg hover:bg-slate-700 transition-colors"
          data-testid="counter-btn"
        >
          Cliques: {{ count() }}
        </button>
      </section>
    </main>
  `,
})
export class AppComponent {
  readonly title = "{{name}}";
  readonly count = signal(0);

  increment(): void {
    this.count.update((n) => n + 1);
  }
}
